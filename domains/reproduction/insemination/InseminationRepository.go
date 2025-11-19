package insemination

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type InseminationRepository struct {
	DB *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *InseminationRepository {
	return &InseminationRepository{db}
}

func (r *InseminationRepository) GetBirthRateStats(userId string) (*BirthRateStats, error) {
	query := `
		with totals as (
			select 
				i.insemination_date,
				count(i.*) as total,
				count(distinct i.id) filter (where a.birth_date is not null) as birth_success
			from insemination_entries i
				left join animals a on a.mother_id = i.animal_id
					and a.birth_date > i.insemination_date
					and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
					and not exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.pregnancy_status = 'FAILED'
							and t.test_date between i.insemination_date and a.birth_date
					)
			where i.user_id = $1 and i.deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select 
			insemination_date, 
			(birth_success::float / nullif(total, 0)::float) * 100 as birth_rate
		from totals
		order by 1
    `
	result, err := repositoriesUtil.GetList[BirthRateHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	birthRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(birthRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = birthRates[lenght-1].BirthRate
		previousRate = 0
		trend = 0
	default:
		currentRate = birthRates[lenght-1].BirthRate
		previousRate = birthRates[lenght-2].BirthRate
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &BirthRateStats{
		Hist:    birthRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *InseminationRepository) GetPregnancyRateStats(userId string) (*CardStats, error) {
	query := `
		with insemination_status as (
			select
				i.insemination_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
									and t.test_date between i.insemination_date and a.birth_date
									and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.insemination_date
						  and age(t.test_date, i.insemination_date) <= interval '340 days'
						  and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.insemination_date
						  and age(t.test_date, i.insemination_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status
			from insemination_entries i
			where i.user_id = $1 and i.deleted_at is null 
		),
		cte as (
			select
				t.insemination_date,
				count(t.*) as total,
				count(t.*) filter (where t.pregnancy_status = 'SUCCESS') as pregnancy_success
			from insemination_status t
			group by 1
			order by 1 desc
			limit 10
		)
		select 
			insemination_date,
			(pregnancy_success::float / nullif(total, 0)) * 100 as pregnancy_rate
		from cte
		order by 1
    `

	result, err := repositoriesUtil.GetList[PregnancyRateHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	pregnancyRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(pregnancyRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = pregnancyRates[lenght-1].PregnancyRate
		previousRate = 0
		trend = 0
	default:
		currentRate = pregnancyRates[lenght-1].PregnancyRate
		previousRate = pregnancyRates[lenght-2].PregnancyRate
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardStats{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *InseminationRepository) GetAnimalsNumber(userId string) (*CardStats, error) {
	query := `
		with cte as (
			select
				insemination_date,
				count(*) as animals_number
			from insemination_entries
			where user_id = $1 and deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select *
		from cte
		order by 1
    `

	result, err := repositoriesUtil.GetList[AnimalsHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	pregnancyRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(pregnancyRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = pregnancyRates[lenght-1].AnimalsNumber
		previousRate = 0
		trend = 0
	default:
		currentRate = pregnancyRates[lenght-1].AnimalsNumber
		previousRate = pregnancyRates[lenght-2].AnimalsNumber
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardStats{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *InseminationRepository) GetInseminationStats(userId string) (*[]InseminationHist, error) {
	query := `
        with cte as (
			select 
				i.insemination_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
									and t.test_date between i.insemination_date and a.birth_date
									and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.insemination_date
						  and age(t.test_date, i.insemination_date) <= interval '340 days'
						  and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.insemination_date
						  and age(t.test_date, i.insemination_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
									and t.test_date between i.insemination_date and a.birth_date
									and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from insemination_entries i
			where i.user_id = $1 and i.deleted_at is null
		),
        totals as (
            select
                insemination_date,
                count(*) total,
                count(*) filter (where birth_status = 'SUCCESS') birth_numbers,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_numbers
            from cte
            group by 1
            order by 1 desc
            limit 30
        )
        select * from totals order by insemination_date
    `
	return repositoriesUtil.GetList[InseminationHist](r.DB, query, userId)
}

func (r *InseminationRepository) GetFutureBirths(userId string) (*[]FutureBirths, error) {
	query := `
		with upcoming_births as (
			select 
				i.id,
				t.birth_forecast
			from insemination_entries i
				join pregnancy_tests t on t.animal_id = i.animal_id
					and t.test_date > i.insemination_date
					and age(t.test_date, i.insemination_date) <= interval '340 days'
					and t.pregnancy_status = 'SUCCESS'
			where i.user_id = $1
				and i.deleted_at is null
				and t.birth_forecast >= now()  
				and not exists (
					select 1
					from animals a 
					where a.mother_id = i.animal_id
						and a.birth_date > i.insemination_date
						and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
						and not exists (
							select 1
							from pregnancy_tests f
							where f.animal_id = i.animal_id
								and f.test_date between i.insemination_date and a.birth_date
								and f.pregnancy_status = 'FAILED'
						)

				)
		)
		select
			date_trunc('month', birth_forecast) as birth_forecast,
			count(distinct id) as births_number
		from upcoming_births
		group by 1
		order by 1;
	`
	return repositoriesUtil.GetList[FutureBirths](r.DB, query, userId)
}

func (r *InseminationRepository) GetBestBull(userId string) (*[]InseminationBulls, error) {
	query := `
		with status as (
			select
				concat_ws(' - ', b.ring_number, b.name) bull_name,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from insemination_entries i
			left join animals b on i.bull_id = b.id 
			where i.user_id = $1 and i.deleted_at is null
		),
		totals as (
			select
				s.bull_name,
				count(s.*) total,
				count(s.*) filter (where s.birth_status = 'SUCCESS') birth_success,
				count(s.*) filter (where s.pregnancy_status = 'SUCCESS') pregnancy_success
			from status s
			group by 1
		),
		rates as (
			select 
				bull_name,
				total,
				(birth_success::float / nullif(total, 0)::float) * 100 birth_rate,
				(pregnancy_success::float / nullif(total, 0)::float) * 100 pregnancy_rate
			from totals
		)
		select
			bull_name,
			total,
			birth_rate,
			pregnancy_rate,
			(birth_rate / nullif(avg(birth_rate) over (), 0) - 1) * 100 as birth_comparison_rate,
			(pregnancy_rate / nullif(avg(pregnancy_rate) over (), 0) - 1) * 100 as pregnancy_comparison_rate
		from rates
		order by birth_rate desc;
    `
	return repositoriesUtil.GetList[InseminationBulls](r.DB, query, userId)
}

func (r *InseminationRepository) GetLastGroups(userId string) (*[]InseminationGroup, error) {
	query := `
		with insemination_data as (
			select
				i.insemination_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from insemination_entries i
			where i.user_id = $1 and i.deleted_at is null
		),
		daily_stats as (
			select
				insemination_date,
				count(*) as cow_number,
				count(*) filter (where birth_status = 'SUCCESS') as birth_success,
				count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancy_success
			from insemination_data
			group by insemination_date
		),
		rates as (
			select
				insemination_date,
				cow_number,
				(birth_success::float * 100 / nullif(cow_number, 0)) as birth_rate,
				(pregnancy_success::float * 100 / nullif(cow_number, 0)) as pregnancy_rate
			from daily_stats
		)
		select
			insemination_date,
			cow_number,
			birth_rate,
			pregnancy_rate,
			coalesce(
				(birth_rate / nullif(lag(birth_rate) over win, 0) - 1) * 100, 0
			) as birth_comparison_rate,
			coalesce(
				(pregnancy_rate / nullif(lag(pregnancy_rate) over win, 0) - 1) * 100, 0
			) as pregnancy_comparison_rate
		from rates
		window win as (order by insemination_date)
		order by insemination_date desc
		limit 5;
    `
	return repositoriesUtil.GetList[InseminationGroup](r.DB, query, userId)
}

func (r *InseminationRepository) GetLastEntries(userId string) (*LastEntry, error) {

	lastDateQuery := `
		select max(insemination_date) max_date
		from insemination_entries 
		where deleted_at is null and user_id = $1
	`

	var lastDate time.Time
	err := repositoriesUtil.GetPrimitive(r.DB, lastDateQuery, &lastDate, userId)
	if err != nil {
		return nil, err
	}

	query := `
		select 
			i.id,
			i.insemination_date,
			i.bull_id,
			concat_ws(' - ', a.ring_number, a.name) as animal_info,
			b.name as bull_name,
			case
				when exists (
					select 1 
					from animals a
					where a.mother_id = i.animal_id
						and a.birth_date > i.insemination_date
						and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
						and not exists (
							select 1
							from pregnancy_tests t
							where t.animal_id = a.mother_id
							and t.test_date between i.insemination_date and a.birth_date
							and t.pregnancy_status = 'FAILED'
						)
				) then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
				) and age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when exists (
					select 1 
					from animals a
					where a.mother_id = i.animal_id
						and a.birth_date > i.insemination_date
						and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
						and not exists (
							select 1
							from pregnancy_tests t
							where t.animal_id = a.mother_id
							and t.test_date between i.insemination_date and a.birth_date
							and t.pregnancy_status = 'FAILED'
						)
				) then 'SUCCESS'
				when age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status
		from insemination_entries i
			left join animals a on a.id = i.animal_id
			left join animals b on b.id = i.bull_id
		where i.user_id = $1 
			and i.insemination_date = $2
			and i.deleted_at is null
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0);
    `
	result, err := repositoriesUtil.GetList[InseminationEntry](r.DB, query, userId, lastDate)
	if err != nil {
		return nil, err
	}

	lastEntry := &LastEntry{
		InseminationDate: lastDate,
		Entries:          *result,
	}

	return lastEntry, nil
}

func (r *InseminationRepository) FindEntriesPage(
	userId string,
	filter InseminationEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[InseminationEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":      {Field: "cte.animal_order", Order: "asc"},
		"animal_name":       {Field: "cte.animal_name", Order: "asc"},
		"insemination_date": {Field: "coalesce(cte.insemination_date, '-infinity')", Order: "asc"},
		"id":                {Field: "cte.id", Order: "asc"},
		"created_at":        {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        with cte as (
			select 
				i.id,
				i.animal_id,
				a.name as animal_name,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
				concat_ws(' - ', a.ring_number, a.name) animal_info,
				i.insemination_date,
				i.bull_id,
				b.name as bull_name,
				case
					when c.child_name is not null then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					when not exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.insemination_date
						  and age(t.test_date, i.insemination_date) <= interval '340 days'
					) and age(i.insemination_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as pregnancy_status,
				case
					when c.child_name is not null then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.insemination_date
						  and age(t.test_date, i.insemination_date) <= interval '340 days'
						  and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when age(i.insemination_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as birth_status,
				case 
					when c.child_name is null then 'Sem Cria'
					else c.child_name
				end as child_information,
				i.observation,
				i.created_at
			from insemination_entries i
				left join animals a on a.id = i.animal_id
				left join animals b on b.id = i.bull_id
				left join lateral (
					select
					concat_ws(
						' - ', 
						a.ring_number, 
						coalesce(a.name, a.sex), 
						to_char(a.birth_date, 'DD/MM/YYYY')
					) as child_name
					from animals a
					where a.mother_id = i.animal_id
						and a.birth_date > i.insemination_date
						and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
						and not exists (
							select 1
							from pregnancy_tests t
							where t.animal_id = i.animal_id
								and t.pregnancy_status = 'FAILED'
								and t.test_date between a.birth_date and i.insemination_date
						)
					order by a.birth_date
					limit 1
				) c on true
			where i.user_id = $1 and i.deleted_at is null
		)
		select * from cte
	`
	orderExpression := " order by "

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
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

	orderExpression += sortExpression
	query += whereExpression + orderExpression
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[InseminationEntry](r.DB, query, sort, 100, args...)
}

func (r *InseminationRepository) GetEntriesFoot(
	userId string,
	filter InseminationEntryFilter,
) (*InseminationFooter, error) {

	statusQuery := `
		with cte as  (
			select
				i.*,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from insemination_entries i
			where i.user_id = $1 and i.deleted_at is null
		)
		select pregnancy_status, birth_status
		from cte i
	`

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "i", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := ""
	if filterExpression != "" {
		whereExpression = " where " + filterExpression
	}

	statusQuery += whereExpression

	query := fmt.Sprintf(`
		with status as (%s),
		totals as (
			select 
				count(*) totals,
				count(*) filter (where birth_status = 'SUCCESS') birth_success,
				count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
			from status
		)
        select 
            totals,
            coalesce(birth_success::float / nullif(totals, 0), 0) * 100 average_birth_rate,
            coalesce(pregnancy_success::float / nullif(totals, 0), 0) * 100 average_pregnancy_rate
		from totals
    `, statusQuery)

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[InseminationFooter](r.DB, query, args...)
}

func (r *InseminationRepository) FindEntriesByGroup(userId string, date time.Time) (*[]InseminationEntry, error) {

	query := `
        select 
            i.id,
			i.animal_id,
			i.bull_id,
			concat_ws(' - ', b.ring_number, b.name) as bull_name,
            concat_ws(' - ', a.ring_number, a.name) as animal_info,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.insemination_date
					  and age(t.test_date, i.insemination_date) <= interval '340 days'
				) and age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case
				when c.child_name is null then 'Sem Cria'
				else child_name
			end as child_information,
            i.observation
        from insemination_entries i
            left join animals a on a.id = i.animal_id
            left join animals b on b.id = i.bull_id
			left join lateral (
				select
				concat_ws(
					' - ', 
					a.ring_number, 
					coalesce(a.name, a.sex), 
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = i.animal_id
					and a.birth_date > i.insemination_date
					and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
					and not exists (
						select 1
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.pregnancy_status = 'FAILED'
							and t.test_date between a.birth_date and i.insemination_date
					)
				order by a.birth_date
				limit 1
			) c on true
		where i.user_id = $1 and i.deleted_at is null and i.insemination_date = $2
        order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return repositoriesUtil.GetList[InseminationEntry](r.DB, query, userId, date)
}

func (r *InseminationRepository) GetEntriesByGroupFoot(userId string, date time.Time) (*InseminationFooter, error) {
	query := `
		with status as (
			select
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
									and t.test_date between i.insemination_date and a.birth_date
									and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from insemination_entries i
			where i.user_id = $1 
				and i.insemination_date = $2
				and i.deleted_at is null
		),
        counting as (
            select 
                count(*) totals,
                count(*) filter (where birth_status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from status
        )
        select 
            totals,
            (birth_success::float / nullif(totals, 0)) * 100 average_birth_rate,
            (pregnancy_success::float / nullif(totals, 0)) * 100 average_pregnancy_rate
        from counting
    `
	return repositoriesUtil.GetOne[InseminationFooter](r.DB, query, userId, date)
}

func (r *InseminationRepository) FindGroups(userId string) (*[]InseminationGroup, error) {
	query := `
		with status as (
			select
				i.insemination_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
									and t.test_date between i.insemination_date and a.birth_date
									and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.test_date > i.insemination_date
							and age(t.test_date, i.insemination_date) <= interval '340 days'
							and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.insemination_date
							and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
							and not exists (
								select 1
								from pregnancy_tests t
								where t.animal_id = a.mother_id
								and t.test_date between i.insemination_date and a.birth_date
								and t.pregnancy_status = 'FAILED'
							)
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from insemination_entries i
			where i.user_id = $1 and i.deleted_at is null
		),
        totals as (
            select 
                insemination_date,
				count(*) cow_number,
                count(*) filter (where birth_status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from status i
            group by insemination_date
        ),
        rates as (
            select
                insemination_date,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            from totals
        )
        select 
            s.insemination_date,
            s.cow_number,
            s.birth_rate,
            s.pregnancy_rate,
            coalesce(
				(s.birth_rate / nullif(lag(s.birth_rate) over win, 0)) - 1, 0
			) * 100 as birth_comparison_rate,
            coalesce(
				(s.pregnancy_rate / nullif(lag(s.pregnancy_rate) over win, 0)) - 1, 0
			) * 100 as pregnancy_comparison_rate
        from rates s
		window win as (order by s.insemination_date)
        order by s.insemination_date desc
    `
	return repositoriesUtil.GetList[InseminationGroup](r.DB, query, userId)
}

func (r *InseminationRepository) SearchInseminationBulls(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select a.id, a.name as label
        from animals a 
        where a.is_insemination_bull = true
			and a.user_id = $1 
			and a.deleted_at is null 
        order by a.name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *InseminationRepository) AddInsemination(entry *InseminationEntrySave) *apiError.APIError {

	validateErr := inseminationExists(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		insert into insemination_entries (animal_id, bull_id, insemination_date, observation, user_id)
		values (:animal_id, :bull_id, :insemination_date, :observation, :user_id)
    `

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *InseminationRepository) ReplaceInsemination(entry *InseminationEntrySave) *apiError.APIError {

	query := `
		update insemination_entries
		set bull_id = :bull_id,
			observation = :observation
		where animal_id = :animal_id
			and insemination_date = :insemination_date
			and user_id = :user_id
    `

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *InseminationRepository) Delete(id string, userId string) *apiError.APIError {

	oldQuery := `
		select id, animal_id, bull_id, insemination_date, user_id
		from insemination_entries
		where id = $1 and user_id = $2
	`

	oldEntry, err := repositoriesUtil.GetOne[InseminationEntrySave](r.DB, oldQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	validateErr := validateDelete(r.DB, oldEntry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		update insemination_entries
		set deleted_at = now()
		where id = :id and and user_id = :user_id
    `

	err = repositoriesUtil.NamedExec(r.DB, query, oldEntry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *InseminationRepository) DeleteNoValidation(id string, userId string) *apiError.APIError {

	query := `
		update insemination_entries
		set deleted_at = now()
		where id = :id and and user_id = :user_id
    `

	err := repositoriesUtil.Exec(r.DB, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *InseminationRepository) DeleteAndChangeFather(id string, userId string) *apiError.APIError {

	tx, err := r.DB.Beginx()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	oldQuery := `
		select id, animal_id, bull_id, insemination_date, user_id
		from insemination_entries
		where id = $1 and user_id = $2
	`

	oldEntry, err := repositoriesUtil.GetOneTx(tx, oldQuery, &InseminationEntrySave{}, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	query := `
		update insemination_entries
		set deleted_at = now()
		where id = :id and user_id = :user_id
    `

	err = repositoriesUtil.NamedExecTx(tx, query, oldEntry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	entriesQuery := `
		with children_query as (
			select
				a.id,
				a.mother_id,
				a.birth_date
			from animals a
			where a.mother_id = $1
				and a.birth_date > $2
				and age(a.birth_date, $2) between interval '240 days' and interval '340 days'
				and a.deleted_at is null 
				and a.user_id = $3
			order by a.birth_date 
			limit 1
		),
		entries_query as (
			select p.pasture_id
			from pasture_entries p
				cross join children_query c
			where p.deleted_at is null
				and p.user_id = $3
				and p.animal_id = $1
				and (c.birth_date - interval '308 days') between p.entry_date and coalesce(p.exit_date, now())
		)
		select pe.animal_id
		from pasture_entries pe
			cross join entries_query e
			cross join children_query c
			join animals a on a.id = pe.animal_id
				and a.animal_type = 'REPRODUCTION_ANIMAL'
				and a.sex = 'M'
		where pe.deleted_at is null
			and pe.user_id = $3
			and (c.birth_date - interval '308 days') between pe.entry_date and coalesce(pe.exit_date, now())
			and pe.pasture_id = e.pasture_id
		order by (coalesce(pe.exit_date, now()) - pe.entry_date) desc
		limit 1
	`

	var fatherId sql.NullString
	err = repositoriesUtil.GetPrimitiveTx(
		tx,
		entriesQuery,
		&fatherId,
		oldEntry.AnimalId,
		oldEntry.InseminationDate,
		oldEntry.UserId,
	)

	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if !fatherId.Valid {
		tx.Commit()
		return nil
	}

	fatherQuery := `
		update animals
		set father_id = $1
		where mother_id = $2
			and father_id = $3
			and birth_date > $4
			and age(birth_date, $4) between interval '240 days' and interval '340 days'
			and user_id = $5
			and deleted_at is null
	`
	err = repositoriesUtil.ExecTx(tx, fatherQuery, fatherId.String, oldEntry.AnimalId, oldEntry.BullId, oldEntry.InseminationDate, oldEntry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	tx.Commit()
	return nil
}

func (r *InseminationRepository) Update(newEntry *InseminationEntrySave) (*InseminationEntry, *apiError.APIError) {

	oldQuery := `
		select id, animal_id, bull_id, insemination_date, observation, user_id
		from insemination_entries
		where id = $1 and user_id = $2
	`

	oldEntry, err := repositoriesUtil.GetOne[InseminationEntrySave](r.DB, oldQuery, newEntry.Id, newEntry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	validateErr := validateUpdate(r.DB, oldEntry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		update insemination_entries
		set bull_id = :bull_id, 
	 		insemination_date = :insemination_date, 
			observation = :observation
		where id = :id and user_id = :user_id
	`

	err = repositoriesUtil.NamedExec(r.DB, query, newEntry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		select 
			i.id,
			i.animal_id,
			concat_ws(' - ', a.ring_number, a.name) animal_info,
			i.insemination_date,
			i.bull_id,
			b.name as bull_name,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.insemination_date
					  and age(t.test_date, i.insemination_date) <= interval '340 days'
				) and age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case 
				when c.child_name is null then 'Sem Cria'
				else c.child_name
			end as child_information,
			i.observation
		from insemination_entries i
			left join animals a on a.id = i.animal_id
			left join animals b on b.id = i.bull_id
			left join lateral (
				select
				concat_ws(
					' - ', 
					a.ring_number, 
					coalesce(a.name, a.sex), 
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = i.animal_id
					and a.birth_date > i.insemination_date
					and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
					and not exists (
						select 1
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.pregnancy_status = 'FAILED'
							and t.test_date between a.birth_date and i.insemination_date
					)
				order by a.birth_date
				limit 1
			) c on true
		where i.id = $1
			and i.user_id = $2
			and i.deleted_at is null
	`

	res, err := repositoriesUtil.GetOne[InseminationEntry](r.DB, selectQuery, newEntry.Id, newEntry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return res, nil
}

func (r *InseminationRepository) UpdateNoValidation(entry *InseminationEntrySave) (*InseminationEntry, *apiError.APIError) {

	query := `
		update insemination_entries 
		set bull_id = :bull_id, 
	 		insemination_date = :insemination_date, 
			observation = :observation
		where id = :id and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		select 
			i.id,
			i.animal_id,
			concat_ws(' - ', a.ring_number, a.name) animal_info,
			i.insemination_date,
			i.bull_id,
			b.name as bull_name,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.insemination_date
					  and age(t.test_date, i.insemination_date) <= interval '340 days'
				) and age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
						and t.test_date > i.insemination_date
						and age(t.test_date, i.insemination_date) <= interval '340 days'
						and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(i.insemination_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case 
				when c.child_name is null then 'Sem Cria'
				else c.child_name
			end as child_information,
			i.observation
		from insemination_entries i
			left join animals a on a.id = i.animal_id
			left join animals b on b.id = i.bull_id
			left join lateral (
				select
				concat_ws(
					' - ', 
					a.ring_number, 
					coalesce(a.name, a.sex), 
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = i.animal_id
					and a.birth_date > i.insemination_date
					and age(a.birth_date, i.insemination_date) between interval '240 days' and interval '340 days'
					and not exists (
						select 1
						from pregnancy_tests t
						where t.animal_id = i.animal_id
							and t.pregnancy_status = 'FAILED'
							and t.test_date between a.birth_date and i.insemination_date
					)
				order by a.birth_date
				limit 1
			) c on true
		where i.id = $1 
			and i.user_id = $2
			and i.deleted_at is null
	`
	res, err := repositoriesUtil.GetOne[InseminationEntry](r.DB, selectQuery, entry.Id, entry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return res, nil
}
