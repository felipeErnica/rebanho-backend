package birth

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type BirthRepository struct {
	SelectQuery string
	TableName   string
	DB          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *BirthRepository {
	selectQuery := `
        SELECT birth.*, 
            animals.name as animal_name, animals.number as animal_number, animals.order as animal_order,
            calf.birth_date as calf_birth_date, calf.sex as calf_sex, father.name as calf_father
        FROM birth_entries as birth 
            LEFT JOIN animals ON animals.id = birth.animal_id
            LEFT JOIN animals as calf ON calf.id = birth.calf_id
            LEFT JOIN animals as father ON father.id = calf.father_id
    `
	return &BirthRepository{selectQuery, "birth", db}
}

func (r *BirthRepository) GetBestIntervals(userId string) (*[]IntervalAnimal, error) {
	query := `
        with average_list as (
            select
                concat_ws(' - ', animals.ring_number, animals.name) as animal_name,
                count(births.*) as birth_numbers,
                avg(births.birth_interval) as interval_average
            from births
                left join animals on animals.id = births.mother_id
            where
                births.user_id = $1
                and births.deleted_at is null
            group by animal_name
        ),
        birth_stats as (
            select 
                min(birth_numbers) as min_births,
                max(birth_numbers) as max_births
            from average_list
        ),
        general_average as (
            select avg(birth_interval) as general_average
            from births 
            where 
                user_id = $1 
                and deleted_at is null 
                and birth_interval is not null
        ),
        scores as (
            select 
                al.*,
                100 * (al.birth_numbers - bs.min_births/bs.max_births - bs.min_births) as birth_score,
                100 * (308 / al.interval_average) as interval_score
            from average_list as al
                cross join birth_stats as bs
            where al.interval_average is not null
        ),
        best_scored as (
            select 
                animal_name,
                interval_average,
                birth_numbers,
                (interval_score*0.6 + birth_score*0.4) as reproductive_score,
                ((interval_average / ga.general_average) - 1) * 100 as average_rate
            from scores
                cross join general_average as ga
            order by reproductive_score desc
            limit 10
        )
        select 
            animal_name,
            interval_average,
            birth_numbers,
            average_rate
        from best_scored
        order by interval_average
    `
	return repositoriesUtil.GetList[IntervalAnimal](r.DB, query, userId)
}

func (r *BirthRepository) GetBirthsStats(userId string) (*BirthStats, error) {
	intervalHistChan := make(chan entity.Result)
	indexHistChan := make(chan entity.Result)
	birthHistChan := make(chan entity.Result)
	pregnantsNumberChan := make(chan entity.Result)
	lossHistoryChan := make(chan entity.Result)

	go r.getBirthIntervalHistory(userId, intervalHistChan)
	go r.getDeathIndex(userId, indexHistChan)
	go r.getBirthHistory(userId, birthHistChan)
	go r.getPregnantsNumber(userId, pregnantsNumberChan)
	go r.getLossHist(userId, lossHistoryChan)

	errs := []error{}

	var intervalStats IntervalStats
	err := util.GetResults(<-intervalHistChan, &intervalStats)
	if err != nil {
		errs = append(errs, err)
	}

	var deathStats DeathStats
	err = util.GetResults(<-indexHistChan, &deathStats)
	if err != nil {
		errs = append(errs, err)
	}

	var currentStats CurrentStats
	err = util.GetResults(<-birthHistChan, &currentStats)
	if err != nil {
		errs = append(errs, err)
	}

	var pregnantNumbers int
	err = util.GetResults(<-pregnantsNumberChan, &pregnantNumbers)
	if err != nil {
		errs = append(errs, err)
	}

	var lossStats LossStats
	err = util.GetResults(<-lossHistoryChan, &lossStats)
	if err != nil {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		return nil, fmt.Errorf("Erros ocorreram: %v", errs)
	}

	birthStats := &BirthStats{
		IntervalHist:        intervalStats.BirthIntervalHist,
		CurrentInterval:     intervalStats.CurrentInterval,
		IntervalTrend:       intervalStats.IntervalTrend,
		DeathIndexHist:      deathStats.DeathIndexHist,
		DeathIndex:          deathStats.DeathIndex,
		DeathTrend:          deathStats.DeathIndexTrend,
		BirthNumbersTrend:   currentStats.BirthNumbersTrend,
		CurrentBirthNumbers: currentStats.CurrentBirthNumbers,
		BirthHistory:        currentStats.BirthHistory,
		PregnantsNumber:     pregnantNumbers,
		LossHist:            lossStats.LossHist,
		Losses:              lossStats.LossNumbers,
		LossTrend:           lossStats.LossTrend,
	}

	return birthStats, nil
}

func (r *BirthRepository) getBirthIntervalHistory(userId string, resultChan chan entity.Result) {
	query := `
        with cte as (
            select distinct
                date_trunc('month', created_at) as month,
                avg(birth_interval) OVER (ORDER BY date_trunc('month', created_at)) AS interval_average
            from births
            where 
                user_id = $1 
                and birth_interval is not null
                and deleted_at is null
                and created_at is not null
            order by 1
        ),
        max_tbl as (select max(month) as max_date from cte)
        select cte.* 
        from cte cross join max_tbl
        where 
            age(month) > interval '1 month'
            and month >= max_date - interval '1 year'
    `
	results, err := repositoriesUtil.GetList[BirthIntervalHist](r.DB, query, userId)

	if err != nil {
		resultChan <- entity.Result{Result: nil, Err: err}
		return
	}

	intervalHist := *results
	var currentInterval, previousInterval, intervalTrend float64

	switch lenght := len(intervalHist); lenght {
	case 0:
		currentInterval = 0
		previousInterval = 0
		intervalTrend = 0
	case 1:
		currentInterval = intervalHist[lenght-1].IntervalAverage
		previousInterval = 0
		intervalTrend = 0
	default:
		currentInterval = intervalHist[lenght-1].IntervalAverage
		previousInterval = intervalHist[lenght-2].IntervalAverage
		intervalTrend = ((currentInterval / previousInterval) - 1) * 100
	}

	intervalResult := IntervalStats{
		IntervalTrend:     intervalTrend,
		BirthIntervalHist: intervalHist,
		CurrentInterval:   currentInterval,
	}

	resultChan <- entity.Result{Result: intervalResult, Err: nil}
}

func (r *BirthRepository) getDeathIndex(userId string, resultsChan chan entity.Result) {
	query := `
        with death_tbl as (
            select 
                date_trunc('quarter', death_date) as death_date,
                count(*) as death_total
        from animals
        where
            user_id = $1
            and animals.death_date is not null
            and age(death_date, birth_date) < interval '1 year'
            and animals.deleted_at is null
        group by 1
        ),
        birth_tbl as (
            select
                date_trunc('quarter', birth_date) as birth_date,
                count(*) as birth_total
            from births
                left join animals on animals.id = births.calf_id
            where 
                births.user_id = $1
                and births.deleted_at is null
            group by 1
        ),
        max_birth_tbl as (select max(birth_date) as max_birth_date from birth_tbl)
        select
            birth_date as date_month,
            (coalesce(death_total, 0)::float / coalesce(birth_total, 0)::float)*100 as death_index
        from birth_tbl
            full join death_tbl on birth_tbl.birth_date = death_tbl.death_date
            cross join max_birth_tbl m
        where 
            age(birth_date) > interval '1 month' 
            and birth_date >= m.max_birth_date - interval '1 year'
        order by 1
    `
	results, err := repositoriesUtil.GetList[DeathIndexHist](r.DB, query, userId)

	if err != nil {
		resultsChan <- entity.Result{Result: nil, Err: err}
		return
	}

	indexHist := *results
	var currentIndex, previousIndex, indexTrend float64

	switch lenght := len(indexHist); lenght {
	case 0:
		currentIndex = 0
		previousIndex = 0
		indexTrend = 0
	case 1:
		currentIndex = indexHist[lenght-1].DeathIndex
		previousIndex = 0
		indexTrend = 0
	default:
		currentIndex = indexHist[lenght-1].DeathIndex
		previousIndex = indexHist[lenght-2].DeathIndex
		indexTrend = ((currentIndex / previousIndex) - 1) * 100
	}

	deathStats := DeathStats{
		DeathIndexHist:  indexHist,
		DeathIndexTrend: indexTrend,
		DeathIndex:      currentIndex,
	}

	resultsChan <- entity.Result{Result: deathStats, Err: nil}
}

func (r *BirthRepository) getBirthHistory(userId string, resultChan chan entity.Result) {
	query := `
        with death_data as ( 
        select 
            date_trunc('month', death_date) as death_month, 
            count(*) as death_total 
            from animals 
            where 
                deleted_at is null 
                and death_date is not null
                and age(death_date, birth_date) < interval '1 year'
                and user_id = $1
            group by death_month 
        ), 
        birth_data as (
            select 
                date_trunc('month', calf.birth_date) as birth_month, 
                count(births.*) as birth_total
            from births 
                left join animals as calf on births.calf_id = calf.id
            where 
                births.user_id = $1
                and births.deleted_at is null 
            group by birth_month 
        ),
        max_date as (select max(birth_month) as max_date from birth_data)
        select
            birth_data.birth_month as date,
            coalesce(birth_data.birth_total,0) as birth_total,
            coalesce(death_data.death_total, 0) as death_total
        from birth_data
            full join death_data on death_data.death_month = birth_data.birth_month
            cross join max_date m
        where birth_data.birth_month  >= m.max_date - interval '5 years'
        order by birth_data.birth_month
    `
	birthHist, err := repositoriesUtil.GetList[BirthsByDate](r.DB, query, userId)

	if err != nil {
		resultChan <- entity.Result{Result: nil, Err: err}
		return
	}

	birthHistArray := *birthHist
	var previousBirth, currentBirth, birthTrend int

	switch lenght := len(birthHistArray); lenght {
	case 0, 1:
		currentBirth = 0
		previousBirth = 0
		birthTrend = 0
	case 2:
		currentBirth = birthHistArray[lenght-1].BirthTotal
		previousBirth = 0
		birthTrend = 0
	default:
		currentBirth = birthHistArray[lenght-1].BirthTotal
		previousBirth = birthHistArray[lenght-2].BirthTotal
		birthTrend = currentBirth - previousBirth
	}

	currentStats := CurrentStats{
		BirthNumbersTrend:   birthTrend,
		CurrentBirthNumbers: currentBirth,
		BirthHistory:        birthHistArray,
	}

	resultChan <- entity.Result{Result: currentStats, Err: nil}
}

func (r *BirthRepository) getPregnantsNumber(userId string, result chan entity.Result) {
	query := `
        select count(*) as pregnant_numbers
        from birth_tests
        where 
            user_id = $1
            and deleted_at is null
            and status = 'PREGNANT'
    `
	var pregnantNumbers int
	err := r.DB.Get(&pregnantNumbers, query, userId)
	if err != nil {
		result <- entity.Result{Result: nil, Err: err}
	}

	result <- entity.Result{Result: pregnantNumbers, Err: nil}
}

func (r *BirthRepository) getLossHist(userId string, resultsChan chan entity.Result) {
	query := `
        with max_date as (select date_trunc('month', max(loss_date)) as max_date from losses)
        select 
            date_trunc('month', loss_date) as month,
            count(*) as losses
        from losses
            cross join max_date m
        where 
            user_id = $1
            and deleted_at is null
            and age(loss_date) > interval '1 month'
            and loss_date >= m.max_date - interval '1 year'
        group by month
        order by month
    `
	results, err := repositoriesUtil.GetList[LossHist](r.DB, query, userId)

	if err != nil {
		resultsChan <- entity.Result{Result: nil, Err: err}
		return
	}

	lossHist := *results
	var currentLosses, previousLosses, lossesTrend int

	switch lenght := len(lossHist); lenght {
	case 0:
		currentLosses = 0
		previousLosses = 0
		lossesTrend = 0
	case 1:
		currentLosses = lossHist[lenght-1].Losses
		previousLosses = 0
		lossesTrend = 0
	default:
		currentLosses = lossHist[lenght-1].Losses
		previousLosses = lossHist[lenght-2].Losses
		lossesTrend = currentLosses - previousLosses
	}

	deathStats := LossStats{
		LossNumbers: currentLosses,
		LossHist:    lossHist,
		LossTrend:   lossesTrend,
	}

	resultsChan <- entity.Result{Result: deathStats, Err: nil}
}

func (r *BirthRepository) TotalBySex(userId string) (*[]TotalBirthsBySex, error) {
	query := `
        select 
            date_trunc('month', calf.birth_date) as birth_month,
            count(births.id) filter (where calf.sex = 'M') as males,
            count(births.id) filter (where calf.sex = 'F') as females
        from births
            left join animals as calf on births.calf_id = calf.id
        where 
            births.user_id = $1 
            and age(calf.birth_date) < interval '5 years'
            and births.deleted_at is null
        group by birth_month
        order by birth_month
    `
	return repositoriesUtil.GetList[TotalBirthsBySex](r.DB, query, userId)
}

func (r *BirthRepository) FindByMotherId(motherId string) (*[]BirthEntry, error) {
	query := r.SelectQuery + " WHERE birth.animal_id = $1 AND birth.deleted_at is null"
	return repositoriesUtil.GetList[BirthEntry](r.DB, query, motherId)
}

func (r *BirthRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.DB, r.TableName, id)
}
