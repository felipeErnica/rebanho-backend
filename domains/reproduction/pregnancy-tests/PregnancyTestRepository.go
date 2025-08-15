package pregnancyTests

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type TestEntryRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TestEntryRepository {
	return &TestEntryRepository{db}
}

func (r *TestEntryRepository) GetPregnancyRate(userId string) (*PregnancyStats, error) {
	query := `
        with cte as (
            select distinct
                date_trunc('month', test_date) test_date,
                count(*) over (order by date_trunc('month', test_date)) totals,
                count(*) filter (where pregnancy_status = 'SUCCESS') over (order by date_trunc('month', test_date)) pregnancies
            from birth_tests
            where user_id = $1 and deleted_at is null
            limit 10
        )
        select 
            test_date,
            (pregnancies::float / totals::float)*100 pregnancy_rate
        from cte
        order by test_date
    `
	result, err := repositoriesUtil.GetList[PregnancyHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	pregnancyHist := *result
	var current, previous, trend float64

	switch lenght := len(pregnancyHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = pregnancyHist[lenght-1].PregnancyRate
		previous = 0
		trend = 0
	default:
		current = pregnancyHist[lenght-1].PregnancyRate
		previous = pregnancyHist[lenght-2].PregnancyRate
		trend = ((current / previous) - 1) * 100
	}

	stats := &PregnancyStats{
		Trend:   trend,
		Current: current,
		Hist:    pregnancyHist,
	}

	return stats, nil
}

func (r *TestEntryRepository) GetBirthRate(userId string) (*BirthStats, error) {
	query := `
        with cte as (
            select distinct
                date_trunc('month', test_date) test_date,
                count(*) over (order by date_trunc('month', test_date)) totals,
                count(*) filter (where birth_status = 'SUCCESS') over (order by date_trunc('month', test_date)) births
            from birth_tests
            where 
                user_id = $1    
                and deleted_at is null
                and age(test_date) > interval '11 months'
            limit 10
        )
        select 
            test_date,
            (births::float / totals::float)*100 birth_rate
        from cte
        order by test_date
    `
	result, err := repositoriesUtil.GetList[BirthHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	birthHist := *result
	var current, previous, trend float64

	switch lenght := len(birthHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = birthHist[lenght-1].BirthRate
		previous = 0
		trend = 0
	default:
		current = birthHist[lenght-1].BirthRate
		previous = birthHist[lenght-2].BirthRate
		trend = ((current / previous) - 1) * 100
	}

	stats := &BirthStats{
		Trend:   trend,
		Current: current,
		Hist:    birthHist,
	}

	return stats, nil
}

func (r *TestEntryRepository) GetPregnancyTestHist(userId string) (*[]PregnancyTestHist, error) {
	query := `
        with cte as (
            select distinct
                date_trunc('month', test_date) test_date,
                count(*) totals,
                count(*) filter (where birth_status = 'SUCCESS') births,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancies
            from birth_tests
            where 
                user_id = $1 
                and deleted_at is null
                and age(test_date) > interval '11 months'
            group by 1
            limit 60
        )
        select 
            test_date,
            totals,
            (births::float / totals::float)*100 birth_rate,
            (pregnancies::float / totals::float)*100 pregnancy_rate
        from cte
        order by test_date
    `
    return repositoriesUtil.GetList[PregnancyTestHist](r.DB, query, userId)
}
