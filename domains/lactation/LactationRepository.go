package lactation

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type LactationRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *LactationRepository {
	return &LactationRepository{db}
}

func (r *LactationRepository) GetYearlyMilk(userId string) (*CardContainer, error) {
	query := `
        select
            date_trunc('year', l.entry_date) entry_date,
            sum(l.quantity) total_milk
        from milk_entries l
        where l.user_id = $1 and l.deleted_at is null
        group by 1
        order by 1
        limit 10
    `
	result, err := repositoriesUtil.GetList[YearProductionHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	averageHist := *result
	var current, previous, trend float64

	switch lenght := len(averageHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = averageHist[0].TotalMilk
		previous = 0
		trend = 0
	default:
		current = averageHist[lenght-1].TotalMilk
		previous = averageHist[lenght-2].TotalMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	averageMilk := &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}

	return averageMilk, nil
}

func (r *LactationRepository) GetMonthMilk(userId string) (*CardContainer, error) {
	query := `
        with month_sum as (
            select
                entry_date,
                sum(quantity) total_milk
            from milk_entries 
            where user_id = $1 and deleted_at is null
            group by 1
            order by 1
        ),
        cte as (
            select 
                date_trunc('month', entry_date) entry_date,
                sum(total_milk) total_milk
            from month_sum
            group by 1
            order by 1 desc
            limit 10
        )
        select * from cte order by entry_date
    `
	result, err := repositoriesUtil.GetList[MonthMilkHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	averageHist := *result
	var current, previous, trend float64

	switch lenght := len(averageHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = averageHist[0].TotalMilk
		previous = 0
		trend = 0
	default:
		current = averageHist[lenght-1].TotalMilk
		previous = averageHist[lenght-2].TotalMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	averageMilk := &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}

	return averageMilk, nil
}

func (r *LactationRepository) GetAnimalsAverage(userId string) (*CardContainer, error) {
	query := `
        with animals_sum as (
            select
                entry_date,
                count(*) animals_number
            from milk_entries
            where user_id = $1 and deleted_at is null
            group by 1
        ),
        cte as (
            select
                date_trunc('month', entry_date) entry_date,
                max(animals_number) animals_number
            from animals_sum
            group by 1
            order by 1 desc
            limit 10
        )
        select * from cte order by entry_date
    `

	result, err := repositoriesUtil.GetList[AnimalsAverageHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	averageHist := *result
	var current, previous, trend float64

	switch lenght := len(averageHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = averageHist[0].AnimalsNumber
		previous = 0
		trend = 0
	default:
		current = averageHist[lenght-1].AnimalsNumber
		previous = averageHist[lenght-2].AnimalsNumber
		trend = util.CalculatePercentageTrend(current, previous)
	}

	averageAnimals := &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}

	return averageAnimals, nil
}

func (r *LactationRepository) GetMilkProduction(userId string) (*[]MilkProductionHist, error) {
	query := `
        with animals_sum as (
            select
                entry_date,
                count(*) animals_number,
                sum(quantity) total_milk
            from milk_entries
            where user_id = $1 and deleted_at is null
            group by 1
        ),
        cte as (
            select
                date_trunc('month', entry_date) entry_date,
                max(animals_number) animals_number,
                sum(total_milk) total_milk
            from animals_sum
            group by 1
            order by 1 desc
            limit 60
        )
        select * from cte order by entry_date
    `
	return repositoriesUtil.GetList[MilkProductionHist](r.DB, query, userId)
}

func (r *LactationRepository) GetBestAnimals(userId string) (*[]AnimalsRating, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and l.end_date >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
                join animals a on 
					a.id = l.animal_id
					and a.death_date is null
        ),
        cte as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total,
                avg(l.lac_interval) avg_interval
            from lac_tbl l join animals a on a.id = l.animal_id
            group by 1
        ),
        cte_stats as (
            select
                avg(lac_num) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(lac_num) dev_lac,
                stddev(avg_total) dev_total
            from cte
        ),
		cte_scores as (
			select
				cte.*,
				((cte.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((cte.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((cte.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((cte.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(cte.lac_num - s.avg_lac)/ nullif(s.dev_lac, 0) as z_lac,
				(cte.avg_total - s.avg_total)/ nullif(s.dev_total, 0) as z_total,
				(cte.avg_interval - s.avg_interval)/ nullif(s.dev_interval, 0) as z_interval
			from cte, cte_stats s
			where cte.avg_interval <> 0
		)
		select *
		from cte_scores
		order by (
			case 
				when z_total < 0 and -z_interval < 0 then z_total - z_interval 
				else (z_lac*0.2 - z_interval*0.4 + z_total*0.4)
			end
		) desc
		limit 10
    `
	return repositoriesUtil.GetList[AnimalsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetWorstAnimals(userId string) (*[]AnimalsRating, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.end_date is not null
                    and l.start_date <= m.entry_date
                    and l.end_date >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
                join animals a on 
					a.id = l.animal_id
					and a.death_date is null
        ),
        cte as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total,
                avg(l.lac_interval) avg_interval
            from lac_tbl l join animals a on a.id = l.animal_id
            group by 1
        ),
        cte_stats as (
            select
                avg(lac_num) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(lac_num) dev_lac,
                stddev(avg_total) dev_total
            from cte
        ),
		cte_scores as (
			select
				cte.*,
				((cte.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((cte.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((cte.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((cte.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(cte.lac_num - s.avg_lac)/ nullif(s.dev_lac, 0) as z_lac,
				(cte.avg_total - s.avg_total)/ nullif(s.dev_total, 0) as z_total,
				(cte.avg_interval - s.avg_interval)/ nullif(s.dev_interval, 0) as z_interval
			from cte, cte_stats s
			where cte.avg_interval <> 0
		)
		select *
		from cte_scores
		order by (
			case 
				when z_total > 0 and -z_interval > 0 then - z_total + z_interval 
				else (z_lac*0.2 + z_interval*0.4 - z_total*0.4)
			end
		) desc
		limit 10
    `
	return repositoriesUtil.GetList[AnimalsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetBestMothers(userId string) (*[]ParentsRating, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and l.end_date >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
                join animals a on 
					a.id = l.animal_id
					and a.death_date is null
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total,
                avg(l.lac_interval) avg_interval
            from lac_tbl l 
            group by 1
        ),
        cte as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
                avg(lac_num) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.mother_id
					and f.death_date is null
			where lac_num > 1
            group by 1
        ),
        cte_stats as (
            select
                avg(avg_lac) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_lac) dev_lac,
                stddev(avg_total) dev_total
            from cte
        ),
		cte_scores as (
			select
				cte.*,
				((cte.avg_lac / nullif(s.avg_lac, 0) ) - 1)*100 lac_rate,
				((cte.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((cte.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((cte.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((cte.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(cte.avg_lac - s.avg_lac)/ nullif(s.dev_lac, 0) as z_lac,
				(cte.avg_total - s.avg_total)/ nullif(s.dev_total, 0) as z_total,
				(cte.avg_interval - s.avg_interval)/ nullif(s.dev_interval, 0) as z_interval
			from cte, cte_stats s
		)
		select *
		from cte_scores
		order by (
			case 
				when z_total < 0 and -z_interval < 0 then z_total - z_interval 
				else (z_lac*0.2 - z_interval*0.4 + z_total*0.4)
			end
		) desc
		limit 10
    `
	return repositoriesUtil.GetList[ParentsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetWorstMothers(userId string) (*[]ParentsRating, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and l.end_date >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
                join animals a on 
					a.id = l.animal_id
					and a.death_date is null
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total,
                avg(l.lac_interval) avg_interval
            from lac_tbl l 
            group by 1
        ),
        cte as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
                avg(lac_num) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.mother_id
					and f.death_date is null
			where lac_num > 1
            group by 1
        ),
        cte_stats as (
            select
                avg(avg_lac) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_lac) dev_lac,
                stddev(avg_total) dev_total
            from cte
        ),
		cte_scores as (
			select
				cte.*,
				((cte.avg_lac / nullif(s.avg_lac, 0) ) - 1)*100 lac_rate,
				((cte.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((cte.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((cte.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((cte.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(cte.avg_lac - s.avg_lac)/ nullif(s.dev_lac, 0) as z_lac,
				(cte.avg_total - s.avg_total)/ nullif(s.dev_total, 0) as z_total,
				(cte.avg_interval - s.avg_interval)/ nullif(s.dev_interval, 0) as z_interval
			from cte, cte_stats s
		)
		select *
		from cte_scores
		order by (
			case 
				when z_total > 0 and -z_interval > 0 then -z_total + z_interval 
				else (z_lac*0.2 + z_interval*0.4 - z_total*0.4)
			end
		) desc
		limit 10
    `
	return repositoriesUtil.GetList[ParentsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetBestFathers(userId string) (*[]ParentsRating, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and l.end_date >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
                join animals a on 
					a.id = l.animal_id
					and a.death_date is null
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total,
                avg(l.lac_interval) avg_interval
            from lac_tbl l 
            group by 1
        ),
        cte as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
                avg(lac_num) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.father_id
					and f.death_date is null
			where lac_num > 1
            group by 1
        ),
        cte_stats as (
            select
                avg(avg_lac) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_lac) dev_lac,
                stddev(avg_total) dev_total
            from cte
        ),
		cte_scores as (
			select
				cte.*,
				((cte.avg_lac / nullif(s.avg_lac, 0) ) - 1)*100 lac_rate,
				((cte.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((cte.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((cte.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((cte.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(cte.avg_lac - s.avg_lac)/ nullif(s.dev_lac, 0) as z_lac,
				(cte.avg_total - s.avg_total)/ nullif(s.dev_total, 0) as z_total,
				(cte.avg_interval - s.avg_interval)/ nullif(s.dev_interval, 0) as z_interval
			from cte, cte_stats s
		)
		select *
		from cte_scores
		order by (
			case 
				when z_total < 0 and -z_interval < 0 then z_total - z_interval 
				else (z_lac*0.2 - z_interval*0.4 + z_total*0.4)
			end
		) desc
		limit 10
    `
	return repositoriesUtil.GetList[ParentsRating](r.DB, query, userId)
}
func (r *LactationRepository) GetWorstFathers(userId string) (*[]ParentsRating, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and l.end_date >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
                join animals a on 
					a.id = l.animal_id
					and a.death_date is null
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total,
                avg(l.lac_interval) avg_interval
            from lac_tbl l 
            group by 1
        ),
        cte as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
                avg(lac_num) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.father_id
					and f.death_date is null
			where lac_num > 1
            group by 1
        ),
        cte_stats as (
            select
                avg(avg_lac) avg_lac,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_lac) dev_lac,
                stddev(avg_total) dev_total
            from cte
        ),
		cte_scores as (
			select
				cte.*,
				((cte.avg_lac / nullif(s.avg_lac, 0) ) - 1)*100 lac_rate,
				((cte.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((cte.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((cte.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((cte.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(cte.avg_lac - s.avg_lac)/ nullif(s.dev_lac, 0) as z_lac,
				(cte.avg_total - s.avg_total)/ nullif(s.dev_total, 0) as z_total,
				(cte.avg_interval - s.avg_interval)/ nullif(s.dev_interval, 0) as z_interval
			from cte, cte_stats s
		)
		select *
		from cte_scores
		order by (
			case 
				when z_total > 0 and -z_interval > 0 then -z_total + z_interval 
				else (z_lac*0.2 + z_interval*0.4 - z_total*0.4)
			end
		) desc
		limit 10
    `
	return repositoriesUtil.GetList[ParentsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetLastEntries(userId string) (*[]MilkEntry, error) {
	query := `
		with max_tbl as (
			select max(entry_date) max_date 
			from milk_entries 
			where user_id = $1 and deleted_at is null
		)
		select 
			concat_ws(' - ', a.ring_number, a.name) animal_name,
			m.entry_date,
			m.quantity
		from max_tbl, milk_entries m 
			join animals a on a.id = m.animal_id
		where 
			m.user_id = $1 
			and m.deleted_at is null 
			and m.entry_date = max_date
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return repositoriesUtil.GetList[MilkEntry](r.DB, query, userId)
}

func (r *LactationRepository) GetLastGroups(userId string) (*[]LactationGroup, error) {
	query := `
		with cte as (
			select 
				entry_date,
				count(*) animals_number,
				sum(quantity) total_milk,
				avg(quantity) avg_milk
			from milk_entries
			where user_id = $1 and deleted_at is null
			group by 1
		)
		select 
			cte.*,
			coalesce(animals_number - lag(animals_number) over (order by entry_date), 0) number_difference,
			coalesce(((total_milk / lag(total_milk) over (order by entry_date)) - 1)*100, 0) total_rate,
			coalesce(((avg_milk / lag(avg_milk) over (order by entry_date)) - 1)*100, 0) avg_rate
		from cte
		order by entry_date desc
		limit 10
	`
	return repositoriesUtil.GetList[LactationGroup](r.DB, query, userId)
}

func (r *LactationRepository) FindGroupsPage(
	filter LactationGroupFilter,
	order string,
	cursor string,
	userId string,
) (*entity.Page[LactationGroup], error) {

	sortMap := map[string]repositoriesUtil.SortField{"entry_date": {Field: "cte.entry_date", Order: "asc"}}

	query := `
		with cte as (
			select 
				entry_date,
				count(*) animals_number,
				sum(quantity) total_milk,
				avg(quantity) avg_milk
			from milk_entries
			where user_id = $1 and deleted_at is null
			group by 1
		)
		select 
			cte.*,
			coalesce(animals_number - lag(animals_number) over (order by entry_date), 0) number_difference,
			coalesce(((total_milk / lag(total_milk) over (order by entry_date)) - 1)*100, 0) total_rate,
			coalesce(((avg_milk / lag(avg_milk) over (order by entry_date)) - 1)*100, 0) avg_rate
		from cte
    `
	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, "entry_date", order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	var whereExpression string
	if filterExpression != "" {
		whereExpression = "where " + filterExpression
	}

	if cursorExpression != "" {
		if whereExpression != "" {
			whereExpression += " and " + cursorExpression
		} else {
			whereExpression = "where " + cursorExpression
		}
	}

	query += whereExpression + " order by cte.entry_date " + order
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[LactationGroup](r.DB, query, "entry_date", 100, args...)
}

func (r *LactationRepository) FindEntriesPage(
	filter MilkEntryFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[MilkEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)

	sortMap := map[string]repositoriesUtil.SortField{
		"animal_name":  {Field: "a.name", Order: "asc"},
		"animal_order": {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"entry_date":   {Field: "m.entry_date", Order: "desc"},
		"quantity":     {Field: "m.quantity", Order: "asc"},
		"id":           {Field: "m.id", Order: "asc"},
		"created_at":   {Field: "m.created_at", Order: "asc"},
	}

	query := `
		select
			m.id,
			m.animal_id,
			concat_ws(' - ', a.ring_number, a.name) animal_name,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			m.entry_date,
			m.quantity
		from milk_entries m
			join animals a on a.id = m.animal_id
    `

	whereExpression := "where m.user_id = $1 and m.deleted_at is null"

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "m", 2)
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

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression += " and " + cursorExpression
	}

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query += whereExpression + " order by " + sortExpression
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[MilkEntry](r.DB, query, sort, 100, args...)
}

func (r *LactationRepository) GetEntriesPageFoot(filter MilkEntryFilter, userId string) (*MilkEntryFoot, error) {
	query := `
		select
			count(*) animals_number,
			sum(quantity) total_milk,
			avg(quantity) avg_milk
		from milk_entries m
    `
	whereExpression := "where m.user_id = $1 and m.deleted_at is null"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "m", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query = query + whereExpression
	return repositoriesUtil.GetOne[MilkEntryFoot](r.DB, query, args...)
}

func (r *LactationRepository) GetGroupEntries(userId string, entryDate time.Time) (*[]MilkEntry, error) {

	query := `
		select
			m.id,
			m.animal_id,
			concat_ws(' - ', a.ring_number, a.name) animal_name,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			m.entry_date,
			m.quantity
		from milk_entries m
			join animals a on a.id = m.animal_id
		where m.user_id = $1 and m.deleted_at is null and m.entry_date = $2
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
	return repositoriesUtil.GetList[MilkEntry](r.DB, query, userId, entryDate)
}

func (r *LactationRepository) GetGroupEntriesFoot(userId string, entryDate time.Time) (*MilkEntryFoot, error) {
	query := `
		select
			count(*) animals_number,
			sum(quantity) total_milk,
			avg(quantity) avg_milk
		from milk_entries m
		where m.user_id = $1 and deleted_at is null and m.entry_date = $2
    `
	return repositoriesUtil.GetOne[MilkEntryFoot](r.DB, query, userId, entryDate)
}

func (r *LactationRepository) FindLactationPage(
	filter LactationHistFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[LactationHist], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":    {Field: "cte.animal_order", Order: "asc"},
		"animal_name":     {Field: "cte.animal_name", Order: "asc"},
		"start_date":      {Field: "cte.start_date", Order: "asc"},
		"end_date":        {Field: "cte.end_date", Order: "asc"},
		"calf_birth_date": {Field: "cte.calf_birth_date", Order: "asc"},
		"id":              {Field: "cte.id", Order: "asc"},
		"created_at":      {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        with lac_stats as (
            select
                l.id,
                avg(m.quantity) avg_prod,
				max(entry_date) max_date,
				max(m.quantity) peak
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
		cte as (
			select
				l.id,
				l.animal_id,
				concat_ws(' - ', a.ring_number, a.name) animal_name,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
				l.calf_id,
				c.birth_date calf_birth_date,
				concat_ws(' - ', to_char(c.birth_date, 'DD/MM/YYYY'), c.sex, cf.name) calf_info,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1 lac_period,
				(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1)*s.avg_prod total_production,
				extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval,
				s.peak,
				l.created_at
			from lactations l
				join lac_stats s using (id)
				join animals a on a.id = l.animal_id
				left join animals c on c.id = l.calf_id
				left join animals cf on cf.id = c.father_id
			where l.user_id = $1 and l.deleted_at is null
		)
		select * from cte
    `

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression := ""
	if filterExpression != "" {
		whereExpression = "where " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression += " and " + cursorExpression
		if filterExpression == "" {
			whereExpression = "where " + cursorExpression
		}
	}

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	orderBy := " order by " + sortExpression
	query += whereExpression + orderBy
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	return repositoriesUtil.GetPage[LactationHist](r.DB, query, sort, 100, args...)
}

func (r *LactationRepository) GetLactationPageFoot(filter LactationHistFilter, userId string) (*LactationHistFoot, error) {

	lacQuery := "select * from cte"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		lacQuery += " where " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	mainQuery := fmt.Sprintf(`
        with lac_stats as (
            select
                l.id,
				max(m.entry_date) max_date,
                avg(m.quantity) avg_prod,
				max(m.quantity) peak
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
		cte as (
			select
				l.animal_id,
				c.birth_date calf_birth_date,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1 lac_period,
				(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1)*s.avg_prod total_production,
				extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) lac_interval,
				s.peak
			from lactations l
				join lac_stats s using (id)
				left join animals c on c.id = l.calf_id
			where l.user_id = $1 and l.deleted_at is null
		),
		lac as (%s)
		select 
			count(*) as total_lacs,
			avg(coalesce(lac_interval, 0)) avg_lac_interval,
			avg(lac_period) avg_lac_period,
			avg(total_production) avg_total_production,
			avg(avg_production) avg_production,
			avg(peak) avg_peak
		from lac
	`, lacQuery)

	return repositoriesUtil.GetOne[LactationHistFoot](r.DB, mainQuery, args...)
}

func (r *LactationRepository) GetLactationEntries(lacId string) (*[]MilkEntry, error) {
	query := `
		select
			m.id,
			m.animal_id,
			concat_ws(' - ', a.ring_number, a.name) animal_name,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			m.entry_date,
			m.quantity
		from milk_entries m
			join animals a on a.id = m.animal_id
			join lactations l on 
				l.id = $1
				and m.entry_date >= l.start_date
				and m.entry_date <=	coalesce(l.end_date, now())
				and m.animal_id = l.animal_id
		where m.deleted_at is null
		order by m.entry_date
    `
	return repositoriesUtil.GetList[MilkEntry](r.DB, query, lacId)
}

func (r *LactationRepository) GetLactationEntriesFoot(lacId string) (*MilkEntryFoot, error) {
	query := `
		select
			count(*) animals_number,
			(extract(days from max(entry_date) - min(entry_date)) + 1)*avg(quantity) total_milk,
			avg(quantity) avg_milk
		from milk_entries m
			join lactations l on 
				l.id = $1
				and m.entry_date >= l.start_date
				and m.entry_date <=	coalesce(l.end_date, now())
				and m.animal_id = l.animal_id
		where m.deleted_at is null
    `
	return repositoriesUtil.GetOne[MilkEntryFoot](r.DB, query, lacId)
}
