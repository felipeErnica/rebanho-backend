package pregnancyTests

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type TestEntryRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TestEntryRepository {
	return &TestEntryRepository{db}
}

func (r *TestEntryRepository) GetPregnancyRate(userId string) (*CardStats, error) {
	query := `
        with cte as (
            select 
                test_date,
                count(*) as totals,
                count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancies
            from pregnancy_tests
            where deleted_at is null and user_id = $1 
			group by 1
			order by test_date desc
            limit 10
        )
        select 
            test_date,
            (pregnancies::float / nullif(totals, 0)) * 100 pregnancy_rate
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

	stats := &CardStats{
		Trend:   trend,
		Current: current,
		Hist:    pregnancyHist,
	}

	return stats, nil
}

func (r *TestEntryRepository) GetAnimalsNumber(userId string) (*CardStats, error) {
	query := `
        with cte as (
            select 
                test_date,
                count(*) as totals
            from pregnancy_tests
            where deleted_at is null and user_id = $1 
			group by 1
			order by test_date desc
            limit 10
        )
        select test_date, totals
        from cte
        order by test_date
    `
	result, err := repositoriesUtil.GetList[AnimalsNumberHist](r.DB, query, userId)
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
		current = pregnancyHist[lenght-1].AnimalsNumber
		previous = 0
		trend = 0
	default:
		current = pregnancyHist[lenght-1].AnimalsNumber
		previous = pregnancyHist[lenght-2].AnimalsNumber
		trend = ((current / previous) - 1) * 100
	}

	stats := &CardStats{
		Trend:   trend,
		Current: current,
		Hist:    pregnancyHist,
	}

	return stats, nil
}

func (r *TestEntryRepository) GetBirthRate(userId string) (*BirthStats, error) {
	query := `
        with cte as (
            select
                test_date,
                count(*) as totals,
                count(*) filter (where pregnancy_status = 'SUCCESS' 
					and exists (
						select 1
						from animals a
						where a.mother_id = t.animal_id
							and a.birth_date > t.test_date
							and age(a.birth_date, t.test_date) <= interval '340 days'
					)
				) as births
            from pregnancy_tests t
            where user_id = $1 and deleted_at is null
			group by test_date
			order by test_date desc
            limit 10
        )
        select 
            test_date,
            (births::float / nullif(totals, 0)) * 100 as birth_rate
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
            select 
                test_date,
                count(*) totals,
                count(*) filter (where pregnancy_status = 'SUCCESS' 
					and exists (
						select 1
						from animals a
						where a.mother_id = t.animal_id
							and a.birth_date > t.test_date
							and age(a.birth_date, t.test_date) <= interval '340 days'
					)
				) as births,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancies
            from pregnancy_tests t
            where user_id = $1 and deleted_at is null 
			group by 1
			order by test_date desc
            limit 36
        )
        select 
            test_date,
            totals,
			pregnancies,
			births
        from cte
        order by test_date
    `
	return repositoriesUtil.GetList[PregnancyTestHist](r.DB, query, userId)
}

func (r *TestEntryRepository) GetLastEntries(userId string) (*LastEntries, error) {
	dateQuery := `
		select max(test_date) max_date
		from pregnancy_tests 
		where user_id = $1 and deleted_at is null
	`

	var lastDate time.Time
	err := repositoriesUtil.GetPrimitive(r.DB, dateQuery, &lastDate, userId)
	if err != nil {
		return nil, err
	}

	query := `
        select
			t.id,
			t.animal_id,
            concat_ws(' - ', a.ring_number, a.name) animal_info,
            t.test_date,
            t.birth_forecast,
            t.pregnancy_status,
			case
				when pregnancy_status = 'FAILED' then 'FAILED'
				when exists (
					select 1 
					from animals a
					where a.mother_id = t.animal_id
						and a.birth_date > t.test_date
						and age(a.birth_date, t.test_date) <= interval '340 days'
				) then 'SUCCESS'
				when age(t.test_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
            t.observation
        from pregnancy_tests t
            left join animals a on a.id = t.animal_id
        where t.user_id = $1 
			and t.test_date = $2
			and t.deleted_at is null
        order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
	result, err := repositoriesUtil.GetList[TestEntry](r.DB, query, userId, lastDate)
	if err != nil {
		return nil, err
	}

	lastEntry := &LastEntries{
		TestDate: lastDate,
		Entries:  *result,
	}

	return lastEntry, nil
}

func (r *TestEntryRepository) GetLastGroups(userId string) (*[]TestGroups, error) {
	query := `
        with totals as (
            select 
                test_date,
                count(*) animals_number,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success,
                count(*) filter (where pregnancy_status = 'SUCCESS' 
					and exists (
						select 1 
						from animals a
						where a.mother_id = t.animal_id
							and a.birth_date > t.test_date
							and age(a.birth_date, t.test_date) <= interval '340 days'
					)
				) as birth_success
            from pregnancy_tests t
            where deleted_at is null and user_id = $1 
            group by 1
            limit 6
        ),
        rates as (
            select
                test_date,
                animals_number,
                (pregnancy_success::float / nullif(animals_number, 0)) * 100 pregnancy_rate,
                (birth_success::float / nullif(animals_number, 0)) * 100 birth_rate
            from totals
        )
        select
            test_date,
            animals_number,
            pregnancy_rate,
            birth_rate,
            coalesce(
				(pregnancy_rate / lag(pregnancy_rate) over win) - 1, 0
			) * 100 as pregnancy_comparison,
            coalesce(
				(birth_rate / lag(birth_rate) over win) - 1, 0
			) * 100 as birth_comparison
        from rates
		window win as (order by test_date)
        order by test_date desc
    `
	return repositoriesUtil.GetList[TestGroups](r.DB, query, userId)
}

func (r *TestEntryRepository) GetNextBirths(userId string) (*[]NextBirths, error) {
	query := `
        select 
            date_trunc('month', birth_forecast) birth_forecast,
            count(*) birth_numbers
        from pregnancy_tests t
        where 
            deleted_at is null 
            and user_id = $1
            and birth_forecast > now()
            and pregnancy_status = 'SUCCESS'
            and age(t.test_date) < interval '340 days'
			and not exists (
				select 1
				from animals a
				where a.mother_id = t.animal_id
					and a.birth_date > t.test_date
					and age(a.birth_date, t.test_date) <= interval '340 days'
			)
        group by 1
        order by 1
    `
	return repositoriesUtil.GetList[NextBirths](r.DB, query, userId)
}

func (r *TestEntryRepository) GetBestResults(userId string) (*[]TestAnimal, error) {
	query := `
		with birth_test_enriched as (
			select 
				bt.animal_id,
				bt.test_date,
				bt.pregnancy_status,
				exists (
					select 1
					from animals a
					where a.mother_id = bt.animal_id
					  and a.birth_date > bt.test_date
					  and age(a.birth_date, bt.test_date) <= interval '340 days'
				) as has_valid_birth
			from pregnancy_tests bt
			where bt.deleted_at is null
			  and bt.user_id = $1
		),
		totals as (
			select 
				animal_id,
				count(*) as totals,
				count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancy_success,
				count(*) filter (where pregnancy_status = 'SUCCESS' and has_valid_birth) as birth_success
			from birth_test_enriched
			group by animal_id
			having count(*) >= 5
		),
		rates as (
			select
				animal_id,
				totals,
				(pregnancy_success::float / totals) * 100 as pregnancy_rate,
				(birth_success::float / totals) * 100 as birth_rate
			from totals
		),
		general_totals as (
			select 
				count(*) as totals,
				count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancy_success,
				count(*) filter (where pregnancy_status = 'SUCCESS' and has_valid_birth) as birth_success
			from birth_test_enriched
		),
		general_rates as (
			select
				(pregnancy_success::float / nullif(totals, 0)) * 100 as total_pregnancy_rate,
				(birth_success::float / nullif(totals, 0)) * 100 as total_birth_rate
			from general_totals
		),
		scores as (
			select
				animal_id,
				totals,
				pregnancy_rate,
				birth_rate,
				(birth_rate - avg(birth_rate) over ()) / nullif(stddev(birth_rate) over (), 0) as birth_score,
				(pregnancy_rate - avg(pregnancy_rate) over ()) / nullif(stddev(pregnancy_rate) over (), 0) as pregnancy_score
			from rates 
		)
		select
			concat_ws(' - ', a.ring_number, a.name) as animal_name,
			s.totals,
			s.pregnancy_rate,
			s.birth_rate,
			coalesce((s.pregnancy_rate / nullif(gr.total_pregnancy_rate, 0)) - 1, 0) * 100 as pregnancy_comparison,
			coalesce((s.birth_rate / nullif(gr.total_birth_rate, 0)) - 1, 0) * 100 as birth_comparison
		from scores s
		cross join general_rates gr
		join animals a on a.id = s.animal_id
		where (s.birth_score + s.pregnancy_score) > 0
		order by (s.birth_score * 0.7 + s.pregnancy_score * 0.3) desc
		limit 10;
    `
	return repositoriesUtil.GetList[TestAnimal](r.DB, query, userId)
}

func (r *TestEntryRepository) GetWorstResults(userId string) (*[]TestAnimal, error) {
	query := `
		with birth_test_enriched as (
			select 
				bt.animal_id,
				bt.test_date,
				bt.pregnancy_status,
				exists (
					select 1
					from animals a
					where a.mother_id = bt.animal_id
					  and a.birth_date > bt.test_date
					  and age(a.birth_date, bt.test_date) <= interval '340 days'
				) as has_valid_birth
			from pregnancy_tests bt
			where bt.deleted_at is null
			  and bt.user_id = $1
		),
		totals as (
			select 
				animal_id,
				count(*) as totals,
				count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancy_success,
				count(*) filter (where pregnancy_status = 'SUCCESS' and has_valid_birth) as birth_success
			from birth_test_enriched
			group by animal_id
			having count(*) >= 5
		),
		rates as (
			select
				animal_id,
				totals,
				(pregnancy_success::float / totals) * 100 as pregnancy_rate,
				(birth_success::float / totals) * 100 as birth_rate
			from totals
		),
		general_totals as (
			select 
				count(*) as totals,
				count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancy_success,
				count(*) filter (where pregnancy_status = 'SUCCESS' and has_valid_birth) as birth_success
			from birth_test_enriched
		),
		general_rates as (
			select
				(pregnancy_success::float / nullif(totals, 0)) * 100 as total_pregnancy_rate,
				(birth_success::float / nullif(totals, 0)) * 100 as total_birth_rate
			from general_totals
		),
		scores as (
			select
				animal_id,
				totals,
				pregnancy_rate,
				birth_rate,
				(birth_rate - avg(birth_rate) over ()) / nullif(stddev(birth_rate) over (), 0) as birth_score,
				(pregnancy_rate - avg(pregnancy_rate) over ()) / nullif(stddev(pregnancy_rate) over (), 0) as pregnancy_score
			from rates 
		)
		select
			concat_ws(' - ', a.ring_number, a.name) as animal_name,
			totals,
			pregnancy_rate,
			birth_rate,
			coalesce((pregnancy_rate / nullif(total_pregnancy_rate, 0)) - 1, 0) * 100 as pregnancy_comparison,
			coalesce((birth_rate / nullif(total_birth_rate, 0)) - 1, 0) * 100 as birth_comparison
		from scores s
		cross join general_rates gr
		join animals a on a.id = s.animal_id
		where (birth_score + pregnancy_score) < 0
		order by (-birth_score * 0.7 - pregnancy_score * 0.3) desc
		limit 10;
    `
	return repositoriesUtil.GetList[TestAnimal](r.DB, query, userId)
}

func (r *TestEntryRepository) FindEntriesPage(
	filter TestEntryFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[TestEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":   {Field: "cte.animal_order", Order: "asc"},
		"test_date":      {Field: "cte.test_date", Order: "desc"},
		"birth_forecast": {Field: "coalesce(cte.birth_forecast, '-infinity')", Order: "desc"},
		"animal_name":    {Field: "cte.animal_name", Order: "asc"},
		"id":             {Field: "cte.id", Order: "asc"},
		"created_at":     {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        with cte as (
			select
				t.id,
				t.test_date,
				t.animal_id,
				a.name as animal_name,
				concat_ws(' - ', a.ring_number, a.name) animal_info,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
				t.birth_forecast,
				t.pregnancy_status,
				case
					when pregnancy_status = 'FAILED' then 'FAILED'
					when child_name is not null then 'SUCCESS'
					when age(t.test_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as birth_status,
				case 
					when pregnancy_status = 'FAILED' then 'Sem Cria'
					when child_name is not null then child_name
					else 'Sem Cria'
				end as child_information,
				t.observation,
				t.created_at
			from pregnancy_tests t 
				left join animals a on a.id = t.animal_id
				left join lateral (
					select concat_ws(
						' - ',
						a.ring_number,
						coalesce(a.name, a.sex),
						to_char(a.birth_date, 'DD/MM/YYYY')
					) as child_name
					from animals a
					where a.mother_id = t.animal_id
						and a.birth_date > t.test_date
						and age(a.birth_date, t.test_date) <= interval '340 days'
					limit 1
				) c on true
			where t.user_id = $1 and t.deleted_at is null
		)
		select *
		from cte
    `
	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression := repositoriesUtil.GetWhereExpression(filterExpression, cursorExpression)

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}
	sortExpression = " order by " + sortExpression

	query += whereExpression + sortExpression
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[TestEntry](r.DB, query, sort, 100, args...)
}

func (r *TestEntryRepository) GetEntriesFoot(filter TestEntryFilter, userId string) (*TestEntryFoot, error) {

	countQuery := `
		with cte as (
			select 
				t.*,
				case
					when pregnancy_status = 'FAILED' then 'FAILED'
					when exists (
						select 1
						from animals a
						where a.mother_id = t.animal_id
							and a.birth_date > t.test_date
							and age(a.birth_date, t.test_date) <= interval '340 days'
					) then 'SUCCESS'
					when age(t.test_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as birth_status
			from pregnancy_tests t
			where t.user_id = $1 and t.deleted_at is null
		)
		select 
			count(*) totals,
			count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success,
			count(*) filter (where birth_status = 'SUCCESS') as birth_success
		from cte t
    `

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "t", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := repositoriesUtil.GetWhereExpression(filterExpression)
	countQuery += whereExpression

	query := fmt.Sprintf(`
        with count_query as (%s)
        select 
            totals,
            coalesce(birth_success::float / nullif(totals, 0), 0) * 100 birth_rate,
            coalesce(pregnancy_success::float / nullif(totals, 0), 0) * 100 pregnancy_rate
        from count_query
    `, countQuery)

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	return repositoriesUtil.GetOne[TestEntryFoot](r.DB, query, args...)
}

func (r *TestEntryRepository) FindGroups(userId string) (*[]TestGroups, error) {
	query := `
        with totals as (
            select 
                test_date,
                count(*) animals_number,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success,
                count(*) filter (where pregnancy_status = 'SUCCESS' 
					and exists (
						select 1
						from animals a
						where a.mother_id = t.animal_id
							and a.birth_date > t.test_date
							and age(a.birth_date, t.test_date) <= interval '340 days'
					)
				) as birth_success
            from pregnancy_tests t
            where deleted_at is null and user_id = $1 
            group by 1
        ),
        rates as (
            select
                g.test_date,
                g.animals_number,
                (g.pregnancy_success::float / g.animals_number::float)*100 pregnancy_rate,
                (g.birth_success::float / g.animals_number::float)*100 birth_rate
            from totals g
        )
        select
            test_date,
            animals_number,
            pregnancy_rate,
            birth_rate,
            coalesce((pregnancy_rate / lag(pregnancy_rate) over win) - 1, 0) *100 pregnancy_comparison,
            coalesce((birth_rate / lag(birth_rate) over win) - 1, 0) * 100 birth_comparison
        from rates
		window win as (order by test_date)
        order by test_date desc
    `
	return repositoriesUtil.GetList[TestGroups](r.DB, query, userId)
}

func (r *TestEntryRepository) FindEntriesByGroup(
	sort string,
	order string,
	testDate time.Time,
	userId string,
) (*[]TestEntry, error) {

	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":   {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "desc"},
		"birth_forecast": {Field: "coalesce(t.birth_forecast, '-infinity')", Order: "desc"},
		"animal_name":    {Field: "a.name", Order: "asc"},
	}

	query := `
        select
            t.id,
            t.test_date,
            t.animal_id,
            concat_ws(' - ', a.ring_number, a.name) as animal_info,
            coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
            t.birth_forecast,
            t.pregnancy_status,
			case
				when pregnancy_status = 'FAILED' then 'FAILED'
				when child_name is not null then 'SUCCESS'
				when age(t.test_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case 
				when pregnancy_status = 'FAILED' then 'Sem Cria'
				when child_name is not null then child_name
				else 'Sem Cria'
			end as child_information,
            t.observation,
            t.loss_id,
            t.calf_id,
            t.created_at
        from pregnancy_tests t 
			left join animals a on a.id = t.animal_id
			left join lateral (
				select concat_ws(
					' - ',
					a.ring_number,
					coalesce(a.name, a.sex),
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = t.animal_id
					and a.birth_date > t.test_date
					and age(a.birth_date, t.test_date) <= interval '340 days'
				limit 1
			) c on true
		where t.user_id = $1 and t.test_date = $2 and t.deleted_at is null
    `
	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query = query + " order by " + sortExpression
	return repositoriesUtil.GetList[TestEntry](r.DB, query, userId, testDate)
}

func (r *TestEntryRepository) GetEntriesByGroupFoot(testDate time.Time, userId string) (*TestEntryFoot, error) {
	query := `
        with count_query as (
            select 
                count(*) totals,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success,
                count(*) filter (where pregnancy_status = 'SUCCESS' 
					and exists (
						select 1
						from animals a
						where a.mother_id = t.animal_id
							and a.birth_date > t.test_date
							and age(a.birth_date, t.test_date) <= interval '340 days'
					)
				) birth_success
            from pregnancy_tests t
            where t.test_date = $1 and t.user_id = $2 and t.deleted_at is null
        )
        select 
            totals,
            coalesce(birth_success::float / nullif(totals, 0), 0) * 100 birth_rate,
            coalesce(pregnancy_success::float / nullif(totals, 0), 0) * 100 pregnancy_rate
        from count_query
    `
	return repositoriesUtil.GetOne[TestEntryFoot](r.DB, query, testDate, userId)
}

func (r *TestEntryRepository) Add(entry *TestEntrySave) *apiError.APIError {

	validateErr := validateAdd(r.DB, entry) 
	if validateErr != nil {
		return validateErr
	}
	
	query := `
		insert into pregnancy_tests (test_date, animal_id, pregnancy_status, birth_forecast, observation, user_id)
		values (:test_date, :animal_id, :pregnancy_status, :birth_forecast, :observation, :user_id)
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *TestEntryRepository) Replace(entry *TestEntrySave) *apiError.APIError {

	query := `
		update pregnancy_tests 
		set pregnancy_status = :pregnancy_status, 
			birth_forecast = :birth_forecast, 
			observation = :observation
		where test_date = :test_date
			and animal_id = :animal_id
			and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *TestEntryRepository) Update(entry *TestEntrySave) (*TestEntry, *apiError.APIError) {

	validateErr := validateUpdate(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		update pregnancy_tests 
		set test_date = :test_date,
			pregnancy_status = :pregnancy_status, 
			birth_forecast = :birth_forecast, 
			observation = :observation
		where id = :id and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		select
			t.id,
			t.test_date,
			t.animal_id,
			concat_ws(' - ', a.ring_number, a.name) animal_info,
			t.birth_forecast,
			t.pregnancy_status,
			case
				when pregnancy_status = 'FAILED' then 'FAILED'
				when child_name is not null then 'SUCCESS'
				when age(t.test_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case 
				when pregnancy_status = 'FAILED' then 'Sem Cria'
				when child_name is not null then child_name
				else 'Sem Cria'
			end as child_information,
			t.observation
		from pregnancy_tests t 
			left join animals a on a.id = t.animal_id
			left join lateral (
				select concat_ws(
					' - ',
					a.ring_number,
					coalesce(a.name, a.sex),
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = t.animal_id
					and a.birth_date > t.test_date
					and age(a.birth_date, t.test_date) <= interval '340 days'
				limit 1
			) c on true
		where t.id = $1 and t.user_id = $2
	`
	result, err := repositoriesUtil.GetOne[TestEntry](r.DB, selectQuery, entry.Id, entry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return result, nil
}

func (r *TestEntryRepository) Delete(id string, userId string) *apiError.APIError {
	query := `
		update pregnancy_tests
		set deleted_at = now()
		where id = $1 and user_id = $2
	`

	err := repositoriesUtil.Exec(r.DB, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}
