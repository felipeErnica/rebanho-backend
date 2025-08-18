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

func (r *TestEntryRepository) GetLastEntries(userId string) (*[]TestEntry, error) {
	query := `
        with max_date as (
            select max(test_date) max_date
            from birth_tests 
            where user_id = $1 and deleted_at is null
        )
        select
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            t.test_date,
            t.pregnancy_status,
            t.birth_status
        from max_date m, birth_tests t
            left join animals a on a.id = t.animal_id
        where t.test_date = m.max_date and t.user_id = $1 and t.deleted_at is null
        order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
    return repositoriesUtil.GetList[TestEntry](r.DB, query, userId)
}

func (r *TestEntryRepository) GetLastGroups(userId string) (*[]TestGroups, error) {
	query := `
        with grouped_count as (
            select 
                test_date,
                count(*) animals_number,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success,
                count(*) filter (where birth_status = 'SUCCESS') birth_success
            from birth_tests
            where deleted_at is null and user_id = $1 
            group by 1
            limit 10
        ),
        grouped_stats as (
            select
                g.test_date,
                g.animals_number,
                (g.pregnancy_success::float / g.animals_number::float)*100 pregnancy_rate,
                (g.birth_success::float / g.animals_number::float)*100 birth_rate
            from grouped_count g
        )
        select
            g.test_date,
            g.animals_number,
            g.pregnancy_rate,
            g.birth_rate,
            coalesce(((g.pregnancy_rate / lag(g.pregnancy_rate) over (order by test_date)) - 1)*100, 0) pregnancy_comparison,
            coalesce(((g.birth_rate / lag(g.birth_rate) over (order by test_date)) - 1)*100, 0) birth_comparison
        from grouped_stats g
        order by g.test_date desc
    `
    return repositoriesUtil.GetList[TestGroups](r.DB, query, userId)
}

func (r *TestEntryRepository) GetNextBirths(userId string) (*[]NextBirths, error) {
	query := `
        select 
            date_trunc('month', birth_forecast) birth_forecast,
            count(*) birth_numbers
        from birth_tests
        where 
            deleted_at is null 
            and birth_forecast > now()
            and pregnancy_status = 'SUCCESS'
            and birth_status = 'STAND_BY'
            and user_id = $1
        group by 1
        order by 1
    `
    return repositoriesUtil.GetList[NextBirths](r.DB, query, userId)
}

func (r *TestEntryRepository) GetBestResults(userId string) (*[]TestAnimal, error) {
    query := `
        with grouped_animals as (
            select 
                animal_id,
                count(*) totals,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success,
                count(*) filter (where birth_status = 'SUCCESS') birth_success
            from birth_tests
            where deleted_at is null and user_id = $1
            group by 1
        ),
        grouped_rates as (
            select
                animal_id,
                totals,
                pregnancy_success,
                birth_success,
                (pregnancy_success::float / totals::float)*100 pregnancy_rate,
                (birth_success::float / totals::float)*100 birth_rate
            from grouped_animals
        ),
        total_counts as (
            select 
                count(*) totals,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success,
                count(*) filter (where birth_status = 'SUCCESS') birth_success
            from birth_tests
            where deleted_at is null and user_id = $1
        ),
        total_rates as (
            select
                (pregnancy_success::float / totals::float)*100 total_pregnancy_rate,
                (birth_success::float / totals::float)*100 total_birth_rate
            from total_counts
        )
        select
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            g.totals,
            g.pregnancy_rate,
            g.birth_rate,
            ((g.pregnancy_rate / total_pregnancy_rate) - 1)*100 pregnancy_comparison,
            ((g.birth_rate / total_birth_rate) - 1)*100 birth_comparison
        from total_rates, grouped_rates g
            join animals a on a.id = g.animal_id
        order by pregnancy_success*0.6 + birth_success*0.4 desc
    `
    return repositoriesUtil.GetList[TestAnimal](r.DB, query, userId)
}
