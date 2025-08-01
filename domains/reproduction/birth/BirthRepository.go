package birth

import (
	"fmt"

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

func (r *BirthRepository) GetBirthsStats(userId string) (*BirthStats, error) {
	intervalHistChan := make(chan *[]BirthIntervalHist)
	currentIntervalChan := make(chan float64)
	intervalTrendChan := make(chan float64)

	indexHistChan := make(chan *[]DeathIndexHist)
	currentDeathIndexChan := make(chan float64)
	indexTrendChan := make(chan float64)

	errChan := make(chan error, 2)

	go r.getBirthIntervalHistory(userId, intervalHistChan, currentIntervalChan, intervalTrendChan, errChan)
	go r.getDeathIndex(userId, indexHistChan, currentDeathIndexChan, indexTrendChan, errChan)

	errs := []error{}
	for range cap(errChan) {
		if err := <-errChan; err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) != 0 {
		return nil, fmt.Errorf("Erros ocorreram: %v", errs)
	}

	birthStats := &BirthStats{
		BirthIntervalHist: *<-intervalHistChan,
		CurrentInterval:   <-currentIntervalChan,
		IntervalTrend:     <-intervalTrendChan,
		DeathIndexHist:    *<-indexHistChan,
		DeathIndex:        <-currentDeathIndexChan,
		DeathIndexTrend:   <-indexTrendChan,
	}
	return birthStats, nil
}

func (r *BirthRepository) getBirthIntervalHistory(
	userId string,
	resultsChan chan *[]BirthIntervalHist,
	intervalChan chan float64,
	trendChan chan float64,
	errChan chan error,
) {
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
            order by 1 desc
        )
        select * from cte where age(month) < interval '1 year'
    `
	results, err := repositoriesUtil.GetList[BirthIntervalHist](r.DB, query, userId)
	errChan <- err
	resultsChan <- results

	if err != nil {
		intervalChan <- 0
		trendChan <- 0
		return
	}

	intervalHist := *results
	var currentInterval, previousInterval, intervalTrend float64

	switch lenght := len(intervalHist); lenght {
	case 0, 1:
		currentInterval = 0
		previousInterval = 0
		intervalTrend = 0
	case 2:
		currentInterval = intervalHist[1].IntervalAverage
		previousInterval = 0
		intervalTrend = 0
	default:
		currentInterval = intervalHist[1].IntervalAverage
		previousInterval = intervalHist[2].IntervalAverage
		intervalTrend = ((currentInterval / previousInterval) - 1) * 100
	}

	intervalChan <- currentInterval
	trendChan <- intervalTrend
}

func (r *BirthRepository) getDeathIndex(
	userId string,
	resultsChan chan *[]DeathIndexHist,
	indexChan chan float64,
	trendChan chan float64,
	errChan chan error,
) {
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
        )
        select
            birth_date as date_month,
            (coalesce(death_total, 0)::float / coalesce(birth_total, 0)::float)*100 as death_index
        from birth_tbl
            full join death_tbl on birth_tbl.birth_date = death_tbl.death_date
        where age(birth_date) < interval '2 years'
        order by 1 desc
    `
	results, err := repositoriesUtil.GetList[DeathIndexHist](r.DB, query, userId)
	errChan <- err
	resultsChan <- results

	if err != nil {
		indexChan <- 0
		trendChan <- 0
		return
	}

	indexHist := *results
	var currentIndex, previousIndex, indexTrend float64

	switch lenght := len(indexHist); lenght {
	case 0, 1:
		currentIndex = 0
		previousIndex = 0
		indexTrend = 0
	case 2:
		currentIndex = indexHist[1].DeathIndex
		previousIndex = 0
		indexTrend = 0
	default:
		currentIndex = indexHist[1].DeathIndex
		previousIndex = indexHist[2].DeathIndex
		indexTrend = ((currentIndex / previousIndex) - 1) * 100
	}

	indexChan <- currentIndex
	trendChan <- indexTrend
}

func (r *BirthRepository) GetBirthHistory(userId string) (*[]BirthsByDate, error) {
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
        )
        select
            birth_data.birth_month as date,
            coalesce(birth_data.birth_total,0) as birth_total,
            coalesce(death_data.death_total, 0) as death_total
        from birth_data
            full join death_data on death_data.death_month = birth_data.birth_month
        order by birth_data.birth_month
    `
	return repositoriesUtil.GetList[BirthsByDate](r.DB, query, userId)
}

func (r *BirthRepository) TotalBySex(userId string) (*[]TotalBirthsBySex, error) {
	query := `
        select 
            date_trunc('month', calf.birth_date) as birth_month,
            count(births.id) filter (where calf.sex = 'M') as males,
            count(births.id) filter (where calf.sex = 'F') as females
        from births
            left join animals as calf on births.calf_id = calf.id
        where births.user_id = $1 and births.deleted_at is null
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
