package embryoTransfer

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type TransferRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TransferRepository {
	return &TransferRepository{db}
}

func (r *TransferRepository) GetBirthRateStats(userId string) (*CardEntry, error) {
	query := `
        with totals as (
            select 
                transfer_date,
                count(*) as total,
                count(*) filter (where exists (
					select 1 
					from animals a
					where a.mother_id = et.donor_id
					  and a.birth_date > et.transfer_date
					  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
				)) as birth_success
            from embryo_transfer et
			where et.user_id = $1 and et.deleted_at is null
			group by 1
            order by 1 desc
            limit 10
        )
        select 
            transfer_date,
			coalesce(birth_success::float / nullif(total, 0), 0) * 100 as birth_rate
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

func (r *TransferRepository) GetPregnancyRateStats(userId string) (*CardEntry, error) {
	query := `
		with insemination_status as (
			select
				et.transfer_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status
			from embryo_transfer et
			where et.user_id = $1 and et.deleted_at is null 
		),
		cte as (
			select
				t.transfer_date,
				count(t.*) as total,
				count(t.*) filter (where t.pregnancy_status = 'SUCCESS') as pregnancy_success
			from insemination_status t
			group by 1
			order by 1 desc
			limit 10
		)
		select 
			transfer_date,
			coalesce(pregnancy_success::float / nullif(total, 0), 0) * 100 as pregnancy_rate
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

func (r *TransferRepository) GetAnimalsNumber(userId string) (*CardEntry, error) {
	query := `
		with cte as (
			select
				transfer_date,
				count(*) as animals_number
			from embryo_transfer
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

func (r *TransferRepository) GetTransferHist(userId string) (*[]TransferHist, error) {
	query := `
        with cte as (
			select 
				et.transfer_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
			where et.user_id = $1 and et.deleted_at is null
		),
        totals as (
            select
                transfer_date,
                count(*) animals_number,
                count(*) filter (where birth_status = 'SUCCESS') births_number,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancies_number
            from cte
            group by 1
            order by 1 desc
            limit 30
        )
        select * from totals order by transfer_date
    `
	return repositoriesUtil.GetList[TransferHist](r.DB, query, userId)
}

func (r *TransferRepository) GetFutureBirths(userId string) (*[]FutureBirths, error) {
	query := `
		with upcoming_births as (
			select (t.test_date + 310 - t.pregnancy_time * interval '1 day') as birth_forecast
			from embryo_transfer et
			join pregnancy_tests t
				on t.animal_id = et.receiver_id
				and t.test_date > et.transfer_date
				and age(t.test_date, et.transfer_date) <= interval '340 days'
				and t.pregnancy_status = 'SUCCESS'
			where et.user_id = $1
			  and et.deleted_at is null
			  and not exists (
				  select 1
				  from animals a
				  where a.mother_id = et.donor_id
					and a.birth_date > et.transfer_date
					and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					and not exists (
						select 1
						from pregnancy_tests t
						where t.animal_id = a.mother_id
							and t.test_date between et.transfer_date and a.birth_date
							and t.pregnancy_status = 'FAILED'
					)
			  )
			  and t.test_date + (310 - t.pregnancy_time * interval '1 day') >= now()  
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

func (r *TransferRepository) GetBestBull(userId string) (*[]BestAnimals, error) {
	query := `
		with status as (
			select
				concat_ws(' - ', a.ring_number, a.name) animal_name,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
                left join animals a on et.bull_id = a.id 
			where et.user_id = $1 and et.deleted_at is null
		),
        totals as (
            select
                s.animal_name,
                count(s.*) total,
                count(s.*) filter (where s.birth_status = 'SUCCESS') birth_success,
                count(s.*) filter (where s.pregnancy_status = 'SUCCESS') pregnancy_success
            from status s
            group by 1
        ),
        rates as (
            select 
                animal_name,
                total,
                coalesce(birth_success::float / nullif(total, 0), 0) * 100 birth_rate,
                coalesce(pregnancy_success::float / nullif(total, 0), 0) * 100 pregnancy_rate
            from totals
        )
        select
			animal_name,
			total,
			birth_rate,
			pregnancy_rate,
			coalesce(birth_rate / nullif(avg(birth_rate) over (), 0) - 1, 0) * 100 as birth_comparison_rate,
			coalesce(pregnancy_rate / nullif(avg(pregnancy_rate) over (), 0) - 1, 0) * 100 as pregnancy_comparison_rate
		from rates
		order by birth_rate desc;
    `
	return repositoriesUtil.GetList[BestAnimals](r.DB, query, userId)
}

func (r *TransferRepository) GetBestDonors(userId string) (*[]BestAnimals, error) {
	query := `
		with status as (
			select
				concat_ws(' - ', a.ring_number, a.name) animal_name,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
                left join animals a on a.id = et.donor_id
			where et.user_id = $1 and et.deleted_at is null
		),
        totals as (
            select
                s.animal_name,
                count(s.*) total,
                count(s.*) filter (where s.birth_status = 'SUCCESS') birth_success,
                count(s.*) filter (where s.pregnancy_status = 'SUCCESS') pregnancy_success
            from status s
            group by 1
        ),
        rates as (
            select 
                animal_name,
                total,
                coalesce(birth_success::float / nullif(total, 0), 0) * 100 birth_rate,
                coalesce(pregnancy_success::float / nullif(total, 0), 0) * 100 pregnancy_rate
            from totals
        )
        select
			animal_name,
			total,
			birth_rate,
			pregnancy_rate,
			coalesce(birth_rate / nullif(avg(birth_rate) over (), 0) - 1, 0) * 100 as birth_comparison_rate,
			coalesce(pregnancy_rate / nullif(avg(pregnancy_rate) over (), 0) - 1, 0) * 100 as pregnancy_comparison_rate
		from rates
		order by birth_rate desc;
    `
	return repositoriesUtil.GetList[BestAnimals](r.DB, query, userId)
}

func (r *TransferRepository) GetBestReceivers(userId string) (*[]BestAnimals, error) {
	query := `
		with status as (
			select
				concat_ws(' - ', a.ring_number, a.name) animal_name,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
                left join animals a on a.id = et.receiver_id
			where et.user_id = $1 and et.deleted_at is null
		),
        totals as (
            select
                s.animal_name,
                count(s.*) total,
                count(s.*) filter (where s.birth_status = 'SUCCESS') birth_success,
                count(s.*) filter (where s.pregnancy_status = 'SUCCESS') pregnancy_success
            from status s
            group by 1
        ),
        rates as (
            select 
                animal_name,
                total,
                coalesce(birth_success::float / nullif(total, 0), 0) * 100 birth_rate,
                coalesce(pregnancy_success::float / nullif(total, 0), 0) * 100 pregnancy_rate
            from totals
        )
        select
			animal_name,
			total,
			birth_rate,
			pregnancy_rate,
			coalesce(birth_rate / nullif(avg(birth_rate) over (), 0) - 1, 0) * 100 as birth_comparison_rate,
			coalesce(pregnancy_rate / nullif(avg(pregnancy_rate) over (), 0) - 1, 0) * 100 as pregnancy_comparison_rate
		from rates
		order by birth_rate desc;
    `
	return repositoriesUtil.GetList[BestAnimals](r.DB, query, userId)
}

func (r *TransferRepository) GetLastGroups(userId string) (*[]TransferGroup, error) {
	query := `
		with insemination_data as (
			select
				et.transfer_date,
				case
					when exists (
						select 1 
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1 	
						from animals a
						where a.mother_id = et.donor_id
						  and a.birth_date > et.transfer_date
						  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
			where et.user_id = $1 and et.deleted_at is null
		),
		daily_stats as (
			select
				transfer_date,
				count(*) as cow_number,
				count(*) filter (where birth_status = 'SUCCESS') as birth_success,
				count(*) filter (where pregnancy_status = 'SUCCESS') as pregnancy_success
			from insemination_data
			group by transfer_date
		),
		rates as (
			select
				transfer_date,
				cow_number,
				coalesce(birth_success::float * 100 / nullif(cow_number, 0), 0) as birth_rate,
				coalesce(pregnancy_success::float * 100 / nullif(cow_number, 0), 0) as pregnancy_rate
			from daily_stats
		)
		select
			transfer_date,
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
		window win as (order by transfer_date)
		order by transfer_date desc
		limit 5;
    `
	return repositoriesUtil.GetList[TransferGroup](r.DB, query, userId)
}

func (r *TransferRepository) GetLastEntries(userId string) (*LastEntry, error) {

	lastDateQuery := `
		select max(transfer_date) max_date
		from embryo_transfer 
		where deleted_at is null and user_id = $1
	`

	var lastDate time.Time
	err := repositoriesUtil.GetPrimitive(r.DB, lastDateQuery, &lastDate, userId)
	if err != nil {
		return nil, err
	}

	query := `
		select 
			et.id,
			et.transfer_date,
			et.bull_id,
			concat_ws(' - ', r.ring_number, r.name) as receiver_info,
			concat_ws(' - ', d.ring_number, d.name) as donor_info,
			b.name as bull_name,
			case
				when exists (
					select 1 
					from animals a
					where a.mother_id = et.donor_id
					  and a.birth_date > et.transfer_date
					  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
				) then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
					  and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
				) and age(et.transfer_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when exists (
					select 1 
					from animals a
					where a.mother_id = et.donor_id
					  and a.birth_date > et.transfer_date
					  and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
				) then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
					  and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(et.transfer_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status
		from embryo_transfer et
			left join animals r on r.id = et.receiver_id
			left join animals d on d.id = et.donor_id
			left join animals b on b.id = et.bull_id
		where et.user_id = $1 
			and et.transfer_date = $2
			and et.deleted_at is null
		order by coalesce(regexp_replace(r.ring_number, '[^0-9]', '', 'g')::int, 0);
    `
	result, err := repositoriesUtil.GetList[EmbryoTransfer](r.DB, query, userId, lastDate)
	if err != nil {
		return nil, err
	}

	lastEntry := &LastEntry{
		TransferDate: lastDate,
		Entries:      *result,
	}

	return lastEntry, nil
}

func (r *TransferRepository) FindEntriesPage(
	userId string,
	filter TransferEntryFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[EmbryoTransfer], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"receiver_order": {Field: "et.receiver_order", Order: "asc"},
		"donor_order":    {Field: "et.donor_order", Order: "asc"},
		"receiver_name":  {Field: "et.receiver_name", Order: "asc"},
		"donor_name":     {Field: "et.donor_name", Order: "asc"},
		"transfer_date":  {Field: "coalesce(et.transfer_date, '-infinity')", Order: "asc"},
		"id":             {Field: "et.id", Order: "asc"},
		"created_at":     {Field: "et.created_at", Order: "asc"},
	}

	query := `
        with cte as (
			select 
				et.id,
				et.receiver_id,
				et.donor_id,
				coalesce(regexp_replace(r.ring_number, '[^0-9]', '', 'g')::int, 0) as receiver_order,
				coalesce(regexp_replace(d.ring_number, '[^0-9]', '', 'g')::int, 0) as donor_order,
				r.name as receiver_name,
				d.name as donor_name,
				concat_ws(' - ', r.ring_number, r.name) as receiver_info,
				concat_ws(' - ', d.ring_number, d.name) as donor_info,
				et.transfer_date,
				et.bull_id,
				b.name as bull_name,
				case
					when c.child_name is not null then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					when not exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
					) and age(et.transfer_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as pregnancy_status,
				case
					when c.child_name is not null then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'FAILED'
					) then 'FAILED'
					when age(et.transfer_date) < interval '340 days' then 'STAND_BY'
					else 'FAILED'
				end as birth_status,
				case 
					when c.child_name is null then 'Sem Cria'
					else c.child_name
				end as child_information,
				et.observation,
				et.created_at
			from embryo_transfer et
				left join animals r on r.id = et.receiver_id
				left join animals d on d.id = et.donor_id
				left join animals b on b.id = et.bull_id
				left join lateral (
					select
					concat_ws(
						' - ', 
						a.ring_number, 
						coalesce(a.name, a.sex), 
						to_char(a.birth_date, 'DD/MM/YYYY')
					) as child_name
					from animals a
					where a.mother_id = et.donor_id
						and a.birth_date > et.transfer_date
						and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					order by a.birth_date
					limit 1
				) c on true
			where et.user_id = $1 and et.deleted_at is null
		)
		select * from cte et
	`
	orderExpression := " order by "

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "et", 2)
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
	return repositoriesUtil.GetPage[EmbryoTransfer](r.DB, query, sort, 100, args...)
}

func (r *TransferRepository) GetEntriesFoot(
	userId string,
	filter TransferEntryFilter,
) (*TransferFoot, error) {

	statusQuery := `
		with cte as  (
			select
				et.animal_id,
				et.bull_id,
				et.transfer_date,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
			where et.user_id = $1 and et.deleted_at is null
		)
		select pregnancy_status, birth_status
		from cte et
	`

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "et", 2)
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

	return repositoriesUtil.GetOne[TransferFoot](r.DB, query, args...)
}

func (r *TransferRepository) FindEntriesByGroup(userId string, date time.Time) (*[]EmbryoTransfer, error) {

	query := `
        select 
            et.id,
			et.bull_id,
			et.receiver_id,
			et.donor_id,
			concat_ws(' - ', b.ring_number, b.name) as bull_name,
            concat_ws(' - ', r.ring_number, r.name) receiver_info,
            concat_ws(' - ', d.ring_number, d.name) donor_info,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
					  and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
				) and age(et.transfer_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
					  and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(et.transfer_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case
				when c.child_name is null then 'Sem Cria'
				else child_name
			end as child_information,
            et.observation
        from embryo_transfer et
            left join animals r on r.id = et.receiver_id
            left join animals d on d.id = et.donor_id
            left join animals b on b.id = et.bull_id
			left join lateral (
				select
				concat_ws(
					' - ', 
					a.ring_number, 
					coalesce(a.name, a.sex), 
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = et.donor_id
					and  a.birth_date > et.transfer_date
					and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
				order by a.birth_date
				limit 1
			) c on true
		where et.user_id = $1 and et.deleted_at is null and et.transfer_date = $2
        order by coalesce(regexp_replace(r.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return repositoriesUtil.GetList[EmbryoTransfer](r.DB, query, userId, date)
}

func (r *TransferRepository) GetEntriesByGroupFoot(userId string, date time.Time) (*TransferFoot, error) {
	query := `
		with status as (
			select
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
			where et.user_id = $1 
				and et.transfer_date = $2
				and et.deleted_at is null
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
            coalesce(birth_success::float / nullif(totals, 0), 0) * 100 average_birth_rate,
            coalesce(pregnancy_success::float / nullif(totals, 0), 0) * 100 average_pregnancy_rate
        from counting
    `
	return repositoriesUtil.GetOne[TransferFoot](r.DB, query, userId, date)
}

func (r *TransferRepository) FindGroups(userId string) (*[]TransferGroup, error) {
	query := `
		with status as (
			select
				et.transfer_date,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
			where et.user_id = $1 and et.deleted_at is null
		),
        totals as (
            select 
                transfer_date,
				count(*) cow_number,
                count(*) filter (where birth_status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from status s
            group by transfer_date
        ),
        rates as (
            select
                transfer_date,
                cow_number,
                coalesce(birth_success::float / nullif(cow_number, 0), 0) * 100 birth_rate,
                coalesce(pregnancy_success::float / nullif(cow_number, 0), 0) * 100 pregnancy_rate
            from totals
        )
        select 
            s.transfer_date,
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
		window win as (order by s.transfer_date)
        order by s.transfer_date desc
    `
	return repositoriesUtil.GetList[TransferGroup](r.DB, query, userId)
}

func (r *TransferRepository) SearchTransferBulls(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, name as label
        from animals 
        where is_transfer_bull = true
			and user_id = $1
			and deleted_at is null 
        order by name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *TransferRepository) SearchNonTransferBulls(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, name as label
        from animals 
        where is_transfer_bull = false
			and sex = 'M'
			and animal_type = 'REPRODUCTION_ANIMAL'
			and user_id = $1
			and deleted_at is null 
        order by name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *TransferRepository) UpdateAsTransferBulls(id string, userId string) *apiError.APIError {
	query := `
		update animals
		set is_transfer_bull = true
		where id = $1 and user_id = $2
    `
	err := repositoriesUtil.Exec(r.DB, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) SearchEmbryoDonors(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, name as label
        from animals 
        where is_embryo_donor = true
			and user_id = $1
			and deleted_at is null 
        order by name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *TransferRepository) SearchNonEmbryoDonors(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, name as label
        from animals 
        where is_embryo_donor = false
			and sex = 'F'
			and animal_type in ('DAIRY_ANIMAL', 'REPRODUCTION_ANIMAL')
			and user_id = $1
			and deleted_at is null 
        order by name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *TransferRepository) UpdateAsEmbryoDonors(id string, userId string) *apiError.APIError {
	query := `
		update animals
		set is_embryo_donor = true
		where id = $1 and user_id = $2
    `
	err := repositoriesUtil.Exec(r.DB, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) AddTransfer(entry *EmbryoTransferSave) *apiError.APIError {

	validateErr := validateAdd(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		insert into embryo_transfer (
			receiver_id, 
			donor_id, 
			bull_id, 
			transfer_date, 
			observation, 
			user_id
		)
		values (
			:receiver_id, 
			:donor_id, 
			:bull_id, 
			:transfer_date, 
			:observation, 
			:user_id
		)
    `

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) Replace(entry *EmbryoTransferSave) *apiError.APIError {

	query := `
		update embryo_transfer 
		set donor_id = :donor_id, 
			bull_id = :bull_id, 
			observation = :observation, 
		where receiver_id = :receiver_id 
			and transfer_date = :transfer_date
			and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) Delete(id string, userId string) *apiError.APIError {

	oldQuery := `
		select
			id,
			donor_id,
			receiver_id,
			bull_id,
			transfer_date,
			observation,
			user_id
		from embryo_transfer
		where id = $1 and user_id = $2
	`

	oldEntry, err := repositoriesUtil.GetOne[EmbryoTransferSave](r.DB, oldQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	validateErr := validateDelete(r.DB, oldEntry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		update embryo_transfer
		set deleted_at = now()
		where id = $1 and user_id = $2
	`

	err = repositoriesUtil.Exec(r.DB, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *TransferRepository) Update(entry *EmbryoTransferSave) (*EmbryoTransfer, *apiError.APIError) {

	validateErr := validateUpdate(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		update embryo_transfer
		set donor_id = :donor_id,
			receiver_id = :receiver_id,
			bull_id = :bull_id,
			transfer_date = :transfer_date,
			observation = :observation
		where id = :id and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		select 
			et.id,
			et.receiver_id,
			et.donor_id,
			concat_ws(' - ', r.ring_number, r.name) as receiver_info,
			concat_ws(' - ', d.ring_number, d.name) as donor_info,
			et.transfer_date,
			et.bull_id,
			b.name as bull_name,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
					  and t.pregnancy_status = 'SUCCESS'
				) then 'SUCCESS'
				when not exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
				) and age(et.transfer_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as pregnancy_status,
			case
				when c.child_name is not null then 'SUCCESS'
				when exists (
					select 1 
					from pregnancy_tests t
					where t.animal_id = et.receiver_id
					  and t.test_date > et.transfer_date
					  and age(t.test_date, et.transfer_date) <= interval '340 days'
					  and t.pregnancy_status = 'FAILED'
				) then 'FAILED'
				when age(et.transfer_date) < interval '340 days' then 'STAND_BY'
				else 'FAILED'
			end as birth_status,
			case 
				when c.child_name is null then 'Sem Cria'
				else c.child_name
			end as child_information,
			et.observation
		from embryo_transfer et
			left join animals r on r.id = et.receiver_id
			left join animals d on d.id = et.donor_id
			left join animals b on b.id = et.bull_id
			left join lateral (
				select
				concat_ws(
					' - ', 
					a.ring_number, 
					coalesce(a.name, a.sex), 
					to_char(a.birth_date, 'DD/MM/YYYY')
				) as child_name
				from animals a
				where a.mother_id = et.donor_id
					and a.birth_date > et.transfer_date
					and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
				order by a.birth_date
				limit 1
			) c on true
		where et.id = $1 and et.user_id = $2
	`

	response, err := repositoriesUtil.GetOne[EmbryoTransfer](r.DB, selectQuery, entry.Id, entry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil

}

func (r *TransferRepository) UpdateGroup(transferDate time.Time, entry *TransferGroup) (*TransferGroup, *apiError.APIError) {

	validateErr := validateUpdateGroups(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		update embryo_transfer
		set transfer_date = $1
		where transfer_date = $2 and user_id = $3 and deleted_at is null
	`

	err := repositoriesUtil.Exec(r.DB, query, entry.TransferDate, transferDate, entry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	returnQuery := `
		with status as (
			select
				et.transfer_date,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					when exists (
						select 1 
						from pregnancy_tests t
						where t.animal_id = et.receiver_id
						  and t.test_date > et.transfer_date
						  and age(t.test_date, et.transfer_date) <= interval '340 days'
						  and t.pregnancy_status = 'SUCCESS'
					) then 'SUCCESS'
					else 'FAILED'
				end as pregnancy_status,
				case
					when exists (
						select 1
						from animals a
						where a.mother_id = et.donor_id
							and a.birth_date > et.transfer_date
							and age(a.birth_date, et.transfer_date) between interval '240 days' and interval '340 days'
					) then 'SUCCESS'
					else 'FAILED'
				end as birth_status
			from embryo_transfer et
			where et.user_id = $1 
				and et.transfer_date = $2
				and et.deleted_at is null
		),
        totals as (
            select 
                transfer_date,
				count(*) cow_number,
                count(*) filter (where birth_status = 'SUCCESS') birth_success,
                count(*) filter (where pregnancy_status = 'SUCCESS') pregnancy_success
            from status s
            group by transfer_date
        ),
        rates as (
            select
                transfer_date,
                cow_number,
                coalesce(birth_success::float / nullif(cow_number, 0), 0) * 100 birth_rate,
                coalesce(pregnancy_success::float / nullif(cow_number, 0), 0) * 100 pregnancy_rate
            from totals
        )
        select 
            s.transfer_date,
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
		window win as (order by s.transfer_date)
    `

	response, err := repositoriesUtil.GetOne[TransferGroup](r.DB, returnQuery, entry.UserId, entry.TransferDate)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *TransferRepository) DeleteGroup(transferDate time.Time, userId string) *apiError.APIError {

	query := `
		delete from embryo_transfer
		where transfer_date = $1 and user_id = $2
	`

	err := repositoriesUtil.Exec(r.DB, query, transferDate, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}
