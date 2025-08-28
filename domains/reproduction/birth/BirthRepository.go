package birth

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type BirthRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *BirthRepository {
	return &BirthRepository{db}
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
                100 * (375 / al.interval_average) as interval_score
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
        with min_date as (
            select min(date_trunc('month', created_at)) as month 
            from births
            where 
                user_id = $1
                and birth_interval is not null
                and deleted_at is null
        ),
        date_series as (
            select generate_series(
                min.month, 
                date_trunc('month', now()), 
                interval '1 month'
            ) as month
            from min_date min
        ),
		birth_month as (
			select 
				date_trunc('month', created_at) as month,
				birth_interval
			from births
			where 
                user_id = $1
				and birth_interval is not null
				and deleted_at is null
		),
		cte as (
			select distinct
				d.month,
				avg(b.birth_interval) over (order by d.month) interval_average
			from date_series d left join birth_month b using (month)
			order by 1
		)
		select * from cte where age(month) < interval '1 year'
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
        with date_series as (
            select date
			from date_trunc('month', now() - interval '1 month') max_date,
                date_trunc('month', now() - interval '2 years') min_date, 
                generate_series(min_date, max_date, interval '1 month') date	
        ),
        death_tbl as (
            select 
                date_trunc('quarter', d.date) date,
                count(a.*) deaths
            from date_series d
                left join animals a on date_trunc('month', a.death_date) = d.date 
            where
                a.user_id = $1
                and a.death_date is not null
                and age(a.death_date, a.birth_date) < interval '1 year'
                and a.deleted_at is null
            group by 1
        ),
        birth_tbl as (
            select
                date_trunc('quarter', d.date) date,
                count(b.*) births
            from births b
                left join animals a on a.id = b.calf_id
                right join date_series d on date_trunc('month', a.birth_date) = d.date
            where 
                b.user_id = $1
                and b.deleted_at is null
            group by 1
        )
        select
            date as date_month,
            coalesce((deaths::float / nullif(births, 0)::float)*100, 0) as death_index
        from birth_tbl full join death_tbl using(date)
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
                date_trunc('month', death_date) date,
                count(*) death_total 
                from animals  
            where 
                deleted_at is null 
                and death_date is not null
                and age(death_date, birth_date) < interval '1 year'
                and user_id = $1
            group by 1
        ), 
        birth_data as (
            select 
                date_trunc('month', calf.birth_date) date,
                count(births.*) as birth_total
            from births 
                left join animals as calf on births.calf_id = calf.id
            where 
                births.user_id = $1
                and births.deleted_at is null 
            group by 1
        )
        select
            date,
            coalesce(birth_data.birth_total,0) birth_total,
            coalesce(death_data.death_total, 0) death_total
        from birth_data full join death_data using(date)
        where date >= now() - interval '5 years'
        order by date
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
        select month, count(l.*) losses
        from generate_series(
			date_trunc('month', now() - interval '13 months'),
			date_trunc('month', now() - interval '1 month'),
			interval '1 month'
		) month left join losses l on 
            date_trunc('month', loss_date) = month
            and l.deleted_at is null 
            and l.user_id = $1 
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
            date_trunc('month', calf.birth_date) birth_month,
            count(births.id) filter (where calf.sex = 'M') males,
            count(births.id) filter (where calf.sex = 'F') females
        from births
            left join animals calf on births.calf_id = calf.id
        where 
            births.user_id = $1 
            and calf.birth_date >= now() - interval '5 years'
            and births.deleted_at is null
        group by birth_month
        order by birth_month
    `
	return repositoriesUtil.GetList[TotalBirthsBySex](r.DB, query, userId)
}

func (r *BirthRepository) FindPage(
	userId string,
	sort string,
	order string,
	filter BirthEntryFilter,
	cursor string,
) (*entity.Page[BirthEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"calf_birth_date": {Field: "c.birth_date", Order: "asc"},
		"mother_order":    {Field: "coalesce(regexp_replace(m.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"mother_name":     {Field: "m.name", Order: "asc"},
		"id":              {Field: "m.id", Order: "asc"},
		"created_at":      {Field: "m.created_at", Order: "asc"},
	}

	query := `
        select 
            b.id,
            b.mother_id,
            b.calf_id,
            concat_ws(' - ', m.ring_number, m.name) mother_name,
            coalesce(regexp_replace(m.ring_number, '[^0-9]', '', 'g')::int, 0) mother_order,
            c.birth_date calf_birth_date,
            c.sex calf_sex,
            case 
                when c.name is null then ''
                else concat_ws(' - ', c.ring_number, c.name)
                end as calf_name,
            c.father_id calf_father_id,
            concat_ws(' - ', f.ring_number, f.name) calf_father,
            b.birth_interval,
            b.observation
        from births b
            left join animals m on m.id = b.mother_id
            left join animals c on c.id = b.calf_id
            left join animals f on f.id = c.father_id
    `
	whereExpression := "where b.user_id = $1 and b.deleted_at is null "
	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "b", 2)
	cursorExpression, nextParam, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression += " and " + cursorExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	orderExpression := " order by " + sortExpression
	query += whereExpression + orderExpression
	return repositoriesUtil.GetPage[BirthEntry](r.DB, query, sort, 100, args...)
}

func (r *BirthRepository) FindPageFooter(userId string, filter BirthEntryFilter) (*BirthFooter, error) {
	query := `
        select 
            count(b.*) total,
            avg(b.birth_interval) filter (where b.birth_interval is not null) intervalAverage
        from births b
            left join animals m on m.id = b.mother_id
            left join animals c on c.id = b.calf_id
            left join animals f on f.id = c.father_id
    `
	whereExpression := "where b.user_id = $1 and b.deleted_at is null "
	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "b", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query += whereExpression
	return repositoriesUtil.GetOne[BirthFooter](r.DB, query, args...)
}
