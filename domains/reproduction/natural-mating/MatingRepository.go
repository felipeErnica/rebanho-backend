package naturalMating

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type MatingRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *MatingRepository {
	return &MatingRepository{db}
}

func (r *MatingRepository) GetBirthRateStats(userId string) (*CardEntry, error) {
	query := `
        with totals as (
            select 
                mating_date,
                count(*) as total,
                count(*) filter (where exists (
					select 1 from animals a
					where a.mother_id = i.animal_id
					  and a.birth_date > i.mating_date
					  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
				)) birth_success
            from natural_matings i
			where i.user_id = $1 and i.deleted_at is null
			group by 1
            order by 1 desc
            limit 10
        )
        select 
            mating_date,
			(birth_success::float / nullif(total, 0)::float) * 100 as birth_rate
        from totals
		order by 1
    `
	result, err := repositoriesUtil.GetList[BirthRateEntry](r.DB, query, userId)
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

	stats := &CardEntry{
		Hist:    birthRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *MatingRepository) GetPregnancyRateStats(userId string) (*CardEntry, error) {
	query := `
		with insemination_status as (
			select
				i.mating_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
						  and a.birth_date > i.mating_date
						  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status
			from natural_matings i
			where i.user_id = $1 and i.deleted_at is null 
		),
		cte as (
			select
				t.mating_date,
				count(t.*) as total,
				count(t.*) filter (where t.pregnancy_status = 'SUCCESS') as pregnancy_success
			from insemination_status t
			group by 1
			order by 1 desc
			limit 10
		)
		select 
			mating_date,
			(pregnancy_success::float / nullif(total, 0)) * 100 as pregnancy_rate
		from cte
		order by 1
    `

	result, err := repositoriesUtil.GetList[PregnancyRateEntry](r.DB, query, userId)
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

	stats := &CardEntry{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *MatingRepository) GetAnimalsNumber(userId string) (*CardEntry, error) {
	query := `
		with cte as (
			select
				mating_date,
				count(*) as animals_number
			from natural_matings
			where user_id = $1 and deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select *
		from cte
		order by 1
    `

	result, err := repositoriesUtil.GetList[AnimalsNumberEntry](r.DB, query, userId)
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

	stats := &CardEntry{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (r *MatingRepository) GetInseminationStats(userId string) (*[]MatingHist, error) {
	query := `
        with cte as (
			select 
				i.mating_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
						  and a.birth_date > i.mating_date
						  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
						  and a.birth_date > i.mating_date
						  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from natural_matings i
			where i.user_id = $1 and i.deleted_at is null
		),
        totals as (
            select
                mating_date,
                count(*) total,
                count(*) filter (where birth_status = 'SUCCESS') birth_numbers,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_numbers
            from cte
            group by 1
            order by 1 desc
            limit 30
        )
        select * from totals order by mating_date
    `
	return repositoriesUtil.GetList[MatingHist](r.DB, query, userId)
}

func (r *MatingRepository) GetFutureBirths(userId string) (*[]FutureBirths, error) {
	query := `
		with upcoming_births as (
			select t.birth_forecast
			from natural_matings i
			join pregnancy_tests t
				on t.animal_id = i.animal_id
				and t.test_date > i.mating_date
				and age(t.test_date, i.mating_date) <= interval '340 days'
				and t.pregnancy_status = 'SUCCESS'
			where i.user_id = $1
			  and i.deleted_at is null
			  and not exists (
				  select 1
				  from animals a
				  where a.mother_id = i.animal_id
					and a.birth_date > i.mating_date
					and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
			  )
			  and t.birth_forecast >= now()  
		)
		select
			date_trunc('month', birth_forecast) as birth_forecast,
			count(*) as births_number
		from upcoming_births
		group by 1
		order by 1;
	`
	return repositoriesUtil.GetList[FutureBirths](r.DB, query, userId)
}

func (r *MatingRepository) GetBestBull(userId string) (*[]BestBulls, error) {
	query := `
		with status as (
			select
				concat_ws(' - ', b.ring_number, b.name) bull_name,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
						  and a.birth_date > i.mating_date
						  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 from animals a
						where a.mother_id = i.animal_id
						  and a.birth_date > i.mating_date
						  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from natural_matings i
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
	return repositoriesUtil.GetList[BestBulls](r.DB, query, userId)
}

func (r *MatingRepository) GetLastGroups(userId string) (*[]MatingGroup, error) {
	query := `
		with insemination_data as (
			select
				i.mating_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = i.animal_id
						  and a.birth_date > i.mating_date
						  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 	
						from animals a
						where a.mother_id = i.animal_id
						  and a.birth_date > i.mating_date
						  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from natural_matings i
			where i.user_id = $1 and i.deleted_at is null
		),
		daily_stats as (
			select
				mating_date,
				count(*) as cow_number,
				count(*) filter (where birth_status = 'SUCCESS') as birth_success,
				count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancy_success
			from insemination_data
			group by mating_date
		),
		rates as (
			select
				mating_date,
				cow_number,
				(birth_success::float * 100 / nullif(cow_number, 0)) as birth_rate,
				(pregnancy_success::float * 100 / nullif(cow_number, 0)) as pregnancy_rate
			from daily_stats
		)
		select
			mating_date,
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
		window win as (order by mating_date)
		order by mating_date desc
		limit 5;
    `
	return repositoriesUtil.GetList[MatingGroup](r.DB, query, userId)
}

func (r *MatingRepository) GetLastEntries(userId string) (*LastEntry, error) {

	lastDateQuery := `
		select max(mating_date) max_date
		from natural_matings 
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
			i.mating_date,
			i.bull_id,
			concat_ws(' - ', a.ring_number, a.name) as animal_name,
			b.name as bull_name,
			case
				when exists (
					select 1 
					from animals a
					where a.mother_id = i.animal_id
					  and a.birth_date > i.mating_date
					  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
				) then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.mating_date
					  and age(t.test_date, i.mating_date) <= interval '340 days'
					  and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.mating_date
					  and age(t.test_date, i.mating_date) <= interval '340 days'
				) and age(i.mating_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when exists (
					select 1 
					from animals a
					where a.mother_id = i.animal_id
					  and a.birth_date > i.mating_date
					  and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
				) then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.mating_date
					  and age(t.test_date, i.mating_date) <= interval '340 days'
					  and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(i.mating_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status
		from natural_matings i
			left join animals a on a.id = i.animal_id
			left join animals b on b.id = i.bull_id
		where i.user_id = $1 
			and i.mating_date = $2
			and i.deleted_at is null
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0);
    `
	result, err := repositoriesUtil.GetList[MatingEntry](r.DB, query, userId, lastDate)
	if err != nil {
		return nil, err
	}

	lastEntry := &LastEntry{
		MatingDate: lastDate,
		Entries:    *result,
	}

	return lastEntry, nil
}

func (r *MatingRepository) FindEntriesPage(
	userId string,
	filter MatingEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[MatingEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order": {Field: "i.animal_order", Order: "asc"},
		"name":         {Field: "i.animal_name", Order: "asc"},
		"mating_date":  {Field: "coalesce(i.mating_date, '-infinity')", Order: "asc"},
		"id":           {Field: "i.id", Order: "asc"},
		"created_at":   {Field: "i.created_at", Order: "asc"},
	}

	query := `
        with cte as (
			select 
				i.id,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
				concat_ws(' - ', a.ring_number, a.name) animal_name,
				i.mating_date,
				i.bull_id,
				b.name as bull_name,
				case
					when c.child_name is not null then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					when not exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
					) and age(i.mating_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as pregnancy_status,
				case
					when c.child_name is not null then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when age(i.mating_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as birth_status,
				case 
					when c.child_name is null then 'Sem Cria'
					else c.child_name
				end as child_information,
				i.observation,
				i.created_at
			from natural_matings i
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
						and a.birth_date > i.mating_date
						and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					order by a.birth_date
					limit 1
				) c on true
			where i.user_id = $1 and i.deleted_at is null
		)
		select * from cte i
	`
	orderExpression := " order by "

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "i", 2)
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
	return repositoriesUtil.GetPage[MatingEntry](r.DB, query, sort, 100, args...)
}

func (r *MatingRepository) GetEntriesFoot(
	userId string,
	filter MatingEntryFilter,
) (*MatingFoot, error) {

	statusQuery := `
		with cte as  (
			select
				i.*,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.mating_date
							and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.mating_date
							and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from natural_matings i
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

	return repositoriesUtil.GetOne[MatingFoot](r.DB, query, args...)
}

func (r *MatingRepository) FindEntriesByGroup(userId string, date time.Time) (*[]MatingEntry, error) {

	query := `
        select 
            i.id,
			concat_ws(' - ', b.ring_number, b.name) as bull_name,
            concat_ws(' - ', a.ring_number, a.name) animal_name,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.mating_date
					  and age(t.test_date, i.mating_date) <= interval '340 days'
					  and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.mating_date
					  and age(t.test_date, i.mating_date) <= interval '340 days'
				) and age(i.mating_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = i.animal_id
					  and t.test_date > i.mating_date
					  and age(t.test_date, i.mating_date) <= interval '340 days'
					  and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(i.mating_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case
				when c.child_name is null then 'Sem Cria'
				else child_name
			end as child_information,
            i.observation
        from natural_matings i
            left join animals a on a.id = i.animal_id
            left join animals b on b.id = i.bull_id
			lateral join (
				select
				concat_ws(
					' - ', 
					a.ring_number, 
					coalesce(a.name, a.sex), 
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = i.animal_id
					and  a.birth_date > i.mating_date
					and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
				order by a.birth_date
				limit 1
			) c on true
		where i.user_id = $1 and i.deleted_at is null and i.mating_date = $2
        order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return repositoriesUtil.GetList[MatingEntry](r.DB, query, userId, date)
}

func (r *MatingRepository) GetEntriesByGroupFoot(userId string, date time.Time) (*MatingFoot, error) {
	query := `
		with status as (
			select
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.mating_date
							and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.mating_date
							and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from natural_matings i
			where i.user_id = $1 
				and i.mating_date = $2
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
	return repositoriesUtil.GetOne[MatingFoot](r.DB, query, userId, date)
}

func (r *MatingRepository) FindGroups(userId string) (*[]MatingGroup, error) {
	query := `
		with status as (
			select
				i.mating_date,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.mating_date
							and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = i.animal_id
						  and t.test_date > i.mating_date
						  and age(t.test_date, i.mating_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = i.animal_id
							and a.birth_date > i.mating_date
							and age(a.birth_date, i.mating_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from natural_matings i
			where i.user_id = $1 and i.deleted_at is null
		),
        totals as (
            select 
                mating_date,
				count(*) cow_number,
                count(*) filter (where birth_status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from status i
            group by mating_date
        ),
        rates as (
            select
                mating_date,
                cow_number,
                (birth_success::float / cow_number::float)*100 birth_rate,
                (pregnancy_success::float / cow_number::float)*100 pregnancy_rate
            from totals
        )
        select 
            s.mating_date,
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
		window win as (order by s.mating_date)
        order by s.mating_date desc
    `
	return repositoriesUtil.GetList[MatingGroup](r.DB, query, userId)
}

func (r *MatingRepository) SearchMatingBulls(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select distinct a.id, a.name label
        from animals a 
			join natural_matings i on i.bull_id = a.id
        where i.user_id = $1 and i.deleted_at is null 
        order by a.name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}
