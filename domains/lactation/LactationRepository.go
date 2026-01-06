package lactation

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/apiError"
	pastureEntries "github.com/felipeErnica/rebanho-backend/domains/farm-area/pasture-entries"
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

func (r *LactationRepository) GetLastMilk(userId string) (*CardContainer, error) {
	query := `
        with cte as (
			select
				l.entry_date,
				sum(l.quantity) as total_milk
			from milk_entries l
			where l.user_id = $1 and l.deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
    `
	result, err := repositoriesUtil.GetList[TotalMilkEntry](r.DB, query, userId)
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

func (r *LactationRepository) GetYearMilk(userId string) (*CardContainer, error) {
	query := `
        with cte as (
			select
				date_trunc('year', l.entry_date) as entry_date,
				sum(l.quantity) as total_milk
			from milk_entries l
			where l.user_id = $1 and l.deleted_at is null
			group by 1
			order by 1 desc
			limit 30
		)
		select * from cte order by entry_date
    `
	result, err := repositoriesUtil.GetList[TotalMilkEntry](r.DB, query, userId)
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

func (r *LactationRepository) GetMilkProduction(userId string) (*[]MilkProductionEntry, error) {
	query := `
        with cte as (
			select
				date_trunc('month', l.entry_date) as entry_date,
				sum(l.quantity) as total_milk,
				count(l.animal_id) animals_number
			from milk_entries l
			where l.user_id = $1 and l.deleted_at is null
			group by 1
			order by 1 desc
			limit 60
		)
		select * from cte order by entry_date
    `
	return repositoriesUtil.GetList[MilkProductionEntry](r.DB, query, userId)
}

func (r *LactationRepository) GetLastAverageMilk(userId string) (*CardContainer, error) {
	query := `
        with cte as (
			select
				l.entry_date,
				avg(l.quantity) as avg_milk
			from milk_entries l
			where l.user_id = $1 and l.deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
    `

	result, err := repositoriesUtil.GetList[AverageMilkEntry](r.DB, query, userId)
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
		current = averageHist[0].AverageMilk
		previous = 0
		trend = 0
	default:
		current = averageHist[lenght-1].AverageMilk
		previous = averageHist[lenght-2].AverageMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	averageMilk := &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}

	return averageMilk, nil
}

func (r *LactationRepository) GetYearAverageMilk(userId string) (*CardContainer, error) {
	query := `
        with cte as (
			select
				date_trunc('year', l.entry_date) as entry_date,
				avg(l.quantity) as avg_milk
			from milk_entries l
			where l.user_id = $1 and l.deleted_at is null
			group by 1
			order by 1 desc
			limit 30
		)
		select * from cte order by entry_date
    `

	result, err := repositoriesUtil.GetList[AverageMilkEntry](r.DB, query, userId)
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
		current = averageHist[0].AverageMilk
		previous = 0
		trend = 0
	default:
		current = averageHist[lenght-1].AverageMilk
		previous = averageHist[lenght-2].AverageMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	averageMilk := &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}

	return averageMilk, nil
}

func (r *LactationRepository) GetLastLactating(userId string) (*CardContainer, error) {
	query := `
        with cte as (
            select
                entry_date,
                count(*) animals_number
            from milk_entries
            where user_id = $1 and deleted_at is null
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

func (r *LactationRepository) GetLastDry(userId string) (*CardContainer, error) {
	query := `
		with cte as (
            select
                entry_date,
                count(*) animals_number
            from milk_entries m
				join lactations l on l.animal_id = m.animal_id
					and m.entry_date = l.end_date
					and l.deleted_at is null
            where m.user_id = $1 and m.deleted_at is null
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
		trend = current - previous
	}

	averageAnimals := &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}

	return averageAnimals, nil
}

func (r *LactationRepository) GetDairyTypes(userId string) (*DairyTypes, error) {

	query := `
		with cte as (
			select
				a.id,
				case
					when l.id is null then false
					else true
				end as is_lactating
			from animals a
				left join lactations l on l.animal_id = a.id
					and l.end_date is null
					and l.deleted_at is null
			where a.user_id = $1 
				and a.deleted_at is null
				and a.animal_type = 'DAIRY_ANIMAL'
				and a.death_date is null
		)
		select
			count(*) filter (where is_lactating = false) as dry,
			count(*) filter (where is_lactating = true) as lactating
		from cte
	`

	result, err := repositoriesUtil.GetOne[DairyTypes](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *LactationRepository) GetBestAnimals(userId string) (*[]AnimalsRating, error) {
	query := `
        with lac_prod as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
				l.end_date,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_prod s using (id)
            where exists (
				select 1
				from animals a
				where a.id = l.animal_id and a.death_date is null
			) 
        ),
        lac_stats as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) lac_num,
                avg(l.period) filter (where l.end_date is not null) avg_period,
                avg(l.avg_prod) filter (where l.end_date is not null) avg_prod,
                avg(l.total) filter (where l.end_date is not null) avg_total,
                avg(l.lac_interval) filter (where l.lac_interval is not null) avg_interval
            from lac_tbl l 
				join animals a on a.id = l.animal_id
            group by 1
			having count(l.*) >= 3
        ),
        gn_stats as (
            select
                avg(avg_period) gn_avg_period,
                avg(avg_prod) gn_avg_prod,
                avg(avg_total) gn_avg_total,
                avg(avg_interval) gn_avg_interval,
				stddev(avg_total) dev_total,
				stddev(avg_interval) dev_interval
            from lac_stats
        ),
		cte as (
			select
				l.*,
				((l.avg_period / nullif(s.gn_avg_period, 0)) - 1)*100 as period_rate,
				((l.avg_prod / nullif(s.gn_avg_prod, 0)) - 1)*100 as prod_rate,
				((l.avg_total / nullif(s.gn_avg_total, 0)) - 1)*100 as total_rate,
				((l.avg_interval / nullif(s.gn_avg_interval, 0)) - 1)*100 as interval_rate,
				(l.avg_total - gn_avg_total ) / nullif(s.dev_total, 0) as total_score,
				(l.avg_interval - gn_avg_interval) / nullif(s.dev_interval, 0) as interval_score
			from lac_stats l
				cross join gn_stats s
		)
		select cte.*
		from cte
		where (total_score * 0.6 - interval_score * 0.4) > 0
		order by (total_score * 0.6 - interval_score * 0.4) desc
		limit 10
    `
	return repositoriesUtil.GetList[AnimalsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetWorstAnimals(userId string) (*[]AnimalsRating, error) {
	query := `
        with lac_prod as (
            select
                l.id,
                avg(m.quantity) avg_prod
            from lactations l
                join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
				l.end_date,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_prod s using (id)
            where exists (
				select 1
				from animals a
				where a.id = l.animal_id and a.death_date is null
			)
        ),
        lac_stats as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) lac_num,
                avg(l.period) filter (where l.end_date is not null) avg_period,
                avg(l.avg_prod) filter (where l.end_date is not null)avg_prod,
                avg(l.total) filter (where l.end_date is not null) avg_total,
                avg(l.lac_interval) filter (where l.lac_interval is not null) avg_interval
            from lac_tbl l 
				join animals a on a.id = l.animal_id
            group by 1
			having count(l.*) >= 3
        ),
        gn_stats as (
            select
                avg(avg_period) gn_avg_period,
                avg(avg_prod) gn_avg_prod,
                avg(avg_total) gn_avg_total,
                avg(avg_interval) gn_avg_interval,
				stddev(avg_total) dev_total,
				stddev(avg_interval) dev_interval
            from lac_stats
        ),
		cte as (
			select
				l.*,
				((l.avg_period / nullif(s.gn_avg_period, 0)) - 1)*100 as period_rate,
				((l.avg_prod / nullif(s.gn_avg_prod, 0)) - 1)*100 as prod_rate,
				((l.avg_total / nullif(s.gn_avg_total, 0)) - 1)*100 as total_rate,
				((l.avg_interval / nullif(s.gn_avg_interval, 0)) - 1)*100 as interval_rate,
				(l.avg_total - gn_avg_total ) / nullif(s.dev_total, 0) as total_score,
				(l.avg_interval - gn_avg_interval) / nullif(s.dev_interval, 0) as interval_score
			from lac_stats l
				cross join gn_stats s
		)
		select cte.*
		from cte
		where (-total_score * 0.6 + interval_score * 0.4) > 0
		order by (-total_score * 0.6 + interval_score * 0.4) desc
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
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
				l.end_date,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) filter (where l.end_date is not null) avg_period,
                avg(l.avg_prod) filter (where l.end_date is not null) avg_prod,
                avg(l.total) filter (where l.end_date is not null) avg_total,
                avg(l.lac_interval) filter (where l.lac_interval is not null) avg_interval
            from lac_tbl l 
            group by 1
        ),
        mother_stats as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
				count(cte.*) children_number,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.mother_id
					and f.death_date is null
            group by 1
			having count(cte.*) >= 3
        ),
        gn_stats as (
            select
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_total) dev_total
            from cte_animals
        ),
		cte as (
			select
				m.*,
				((m.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((m.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((m.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((m.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(m.avg_total - s.avg_total) / nullif(s.dev_total, 0) as total_score,
				(m.avg_interval - s.avg_interval) / nullif(s.dev_interval, 0) as interval_score
			from mother_stats m
				cross join gn_stats s
		)
		select *
		from cte
		where (total_score * 0.6 - interval_score * 0.4) > 0
		order by (total_score * 0.6 - interval_score * 0.4) desc
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
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
				l.end_date,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) filter (where l.end_date is not null) avg_period,
                avg(l.avg_prod) filter (where l.end_date is not null) avg_prod,
                avg(l.total) filter (where l.end_date is not null) avg_total,
                avg(l.lac_interval) filter (where l.lac_interval is not null) avg_interval
            from lac_tbl l 
            group by 1
        ),
        mother_stats as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
				count(cte.*) children_number,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.mother_id
					and f.death_date is null
            group by 1
			having count(cte.*) >= 3
        ),
        gn_stats as (
            select
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_total) dev_total
            from cte_animals
        ),
		cte as (
			select
				m.*,
				((m.avg_period / nullif(s.avg_period, 0)) - 1)*100 period_rate,
				((m.avg_prod / nullif(s.avg_prod, 0)) - 1)*100 prod_rate,
				((m.avg_total / nullif(s.avg_total, 0)) - 1)*100 total_rate,
				((m.avg_interval / nullif(s.avg_interval, 0)) - 1)*100 interval_rate,
				(m.avg_total - s.avg_total) / nullif(s.dev_total, 0) as total_score,
				(m.avg_interval - s.avg_interval) / nullif(s.dev_interval, 0) as interval_score
			from mother_stats m
				cross join gn_stats s
		)
		select *
		from cte
		where (-total_score * 0.6 + interval_score * 0.4) > 0
		order by (-total_score * 0.6 + interval_score * 0.4) desc
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
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
				l.end_date,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total,
                extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) filter (where l.end_date is not null) avg_period,
                avg(l.avg_prod) filter (where l.end_date is not null) avg_prod,
                avg(l.total) filter (where l.end_date is not null) avg_total,
                avg(l.lac_interval) filter (where l.lac_interval is not null) avg_interval
            from lac_tbl l 
            group by 1
        ),
        mother_stats as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
				count(cte.*) children_number,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.father_id
					and f.death_date is null
            group by 1
			having count(cte.*) >= 5
        ),
        gn_stats as (
            select
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_total) dev_total
            from cte_animals
        ),
		cte as (
			select
				m.*,
				((m.avg_period / nullif(s.avg_period, 0)) - 1) * 100 period_rate,
				((m.avg_prod / nullif(s.avg_prod, 0)) - 1) * 100 prod_rate,
				((m.avg_total / nullif(s.avg_total, 0)) - 1) * 100 total_rate,
				((m.avg_interval / nullif(s.avg_interval, 0)) - 1) * 100 interval_rate,
				(m.avg_total - s.avg_total) / nullif(s.dev_total, 0) as total_score,
				(m.avg_interval - s.avg_interval) / nullif(s.dev_interval, 0) as interval_score
			from mother_stats m
				cross join gn_stats s
		)
		select *
		from cte
		where (total_score * 0.6 - interval_score * 0.4) > 0
		order by (total_score * 0.6 - interval_score * 0.4) desc
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
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),
        lac_tbl as (
            select
                l.animal_id,
                s.avg_prod,
				l.end_date,
                extract(days from l.end_date - l.start_date) + 1 period,
                (extract(days from l.end_date - l.start_date) + 1) * s.avg_prod total,
                extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) lac_interval
            from lactations l
                join lac_stats s using (id)
        ),
        cte_animals as (
            select
               	l.animal_id,
                count(l.*) lac_num,
                avg(l.period) filter (where l.end_date is not null) avg_period,
                avg(l.avg_prod) filter (where l.end_date is not null) avg_prod,
                avg(l.total) filter (where l.end_date is not null) avg_total,
                avg(l.lac_interval) filter (where l.lac_interval is not null) avg_interval
            from lac_tbl l 
            group by 1
        ),
        mother_stats as (
            select
                concat_ws(' - ', f.ring_number, f.name) parent_name,
				count(cte.*) children_number,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval
            from cte_animals cte
				join animals a on a.id = cte.animal_id
				join animals f on 
					f.id = a.father_id
					and f.death_date is null
            group by 1
			having count(cte.*) >= 5
        ),
        gn_stats as (
            select
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total,
                avg(avg_interval) avg_interval,
                stddev(avg_interval) dev_interval,
                stddev(avg_total) dev_total
            from cte_animals
        ),
		cte as (
			select
				m.*,
				((m.avg_period / nullif(s.avg_period, 0)) - 1) * 100 period_rate,
				((m.avg_prod / nullif(s.avg_prod, 0)) - 1) * 100 prod_rate,
				((m.avg_total / nullif(s.avg_total, 0)) - 1) * 100 total_rate,
				((m.avg_interval / nullif(s.avg_interval, 0)) - 1) * 100 interval_rate,
				(m.avg_total - s.avg_total) / nullif(s.dev_total, 0) as total_score,
				(m.avg_interval - s.avg_interval) / nullif(s.dev_interval, 0) as interval_score
			from mother_stats m
				cross join gn_stats s
			where m.avg_interval > 0 
		)
		select *
		from cte
		where (-total_score * 0.6 + interval_score * 0.4) > 0
		order by (-total_score * 0.6 + interval_score * 0.4) desc
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
			m.id,
			concat_ws(' - ', a.ring_number, a.name) as animal_info,
			coalesce(p.name, 'Sem Pasto') as pasture_name,
			m.entry_date,
			m.quantity
		from milk_entries m 
			cross join max_tbl
			join animals a on a.id = m.animal_id
			left join pasture_entries pe on pe.animal_id = m.animal_id
				and m.entry_date >= pe.entry_date
				and m.entry_date < coalesce(pe.exit_date, now())
				and pe.deleted_at is null
			left join pastures p on p.id = pe.pasture_id
		where m.user_id = $1 
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
		limit 5
	`
	return repositoriesUtil.GetList[LactationGroup](r.DB, query, userId)
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
		"animal_order":     {Field: "cte.animal_order", Order: "asc"},
		"name":             {Field: "cte.name", Order: "asc"},
		"start_date":       {Field: "cte.start_date", Order: "asc"},
		"end_date":         {Field: "coalesce(cte.end_date, '-infinity')", Order: "asc"},
		"calf_birth_date":  {Field: "coalesce(cte.calf_birth_date, -infinity)", Order: "asc"},
		"avg_production":   {Field: "coalesce(cte.avg_production, 0)", Order: "asc"},
		"lac_period":       {Field: "cte.lac_period", Order: "asc"},
		"total_production": {Field: "coalesce(cte.total_production, 0)", Order: "asc"},
		"lac_interval":     {Field: "coalesce(cte.lac_interval, 0)", Order: "asc"},
		"id":               {Field: "cte.id", Order: "asc"},
		"created_at":       {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        with lac_stats as (
            select
                l.id,
                avg(coalesce(m.quantity, 0)) avg_prod,
				max(entry_date) max_date,
				max(coalesce(m.quantity, 0)) peak
            from lactations l
                left join milk_entries m on 
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
				a.name,
				concat_ws(' - ', a.ring_number, a.name) as animal_name,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
				l.calf_id,
				c.birth_date calf_birth_date,
				case
					when l.calf_id is null then 'Sem Bezerro'
					when c.name is not null then format(
						'%s (%s)',
						concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY')),
						concat_ws(' - ', c.ring_number, c.name)
					)
					else concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY'))
				end as calf_info,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) lac_period,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod total_production,
				extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
				s.peak,
				l.observation,
				l.created_at
			from lactations l
				join lac_stats s using (id)
				join animals a on a.id = l.animal_id
				left join animals c on c.id = l.calf_id
				left join animals cm on cm.id = c.mother_id
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

	whereExpression := repositoriesUtil.GetWhereExpression(filterExpression, cursorExpression)

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

func (r *LactationRepository) FindById(id string, userId string) (*LactationHist, error) {

	query := `
        with lac_stats as (
            select
                l.id,
                avg(coalesce(m.quantity, 0)) avg_prod,
				max(entry_date) max_date,
				max(coalesce(m.quantity, 0)) peak
            from lactations l
                left join milk_entries m on 
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
				a.name,
				concat_ws(' - ', a.ring_number, a.name) as animal_name,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
				l.calf_id,
				c.birth_date calf_birth_date,
				case
					when l.calf_id is null then 'Sem Bezerro'
					when c.name is not null then format(
						'%s (%s)',
						concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY')),
						concat_ws(' - ', c.ring_number, c.name)
					)
					else concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY'))
				end as calf_info,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) lac_period,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod total_production,
				extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
				s.peak,
				l.observation,
				l.created_at
			from lactations l
				join lac_stats s using (id)
				join animals a on a.id = l.animal_id
				left join animals c on c.id = l.calf_id
				left join animals cm on cm.id = c.mother_id
			where l.id = $1
				and l.user_id = $2
				and l.deleted_at is null
		)
		select * from cte
    `
	return repositoriesUtil.GetOne[LactationHist](r.DB, query, id, userId)
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
				extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) lac_interval,
				s.peak,
				l.observation
			from lactations l
				join lac_stats s using (id)
				left join animals c on c.id = l.calf_id
			where l.user_id = $1 and l.deleted_at is null
		),
		lac as (%s)
		select 
			count(*) as total_lacs,
			avg(lac_interval) filter (where lac_interval is not null) avg_lac_interval,
			avg(lac_period) filter (where end_date is not null) avg_lac_period,
			avg(total_production) filter (where end_date is not null) avg_total_production,
			avg(avg_production) filter (where end_date is not null) avg_production,
			avg(peak) filter (where end_date is not null) avg_peak
		from lac
	`, lacQuery)

	return repositoriesUtil.GetOne[LactationHistFoot](r.DB, mainQuery, args...)
}

func (r *LactationRepository) GetLactationEntries(lacId string) (*[]MilkEntry, error) {
	query := `
		select
			m.id,
			m.animal_id,
			concat_ws(' - ', a.ring_number, a.name) as animal_info,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
			coalesce(p.name, 'Sem Pasto') as pasture_name,
			m.entry_date,
			m.quantity
		from milk_entries m
			join animals a on a.id = m.animal_id
			join lactations l on 
				l.id = $1
				and m.entry_date >= l.start_date
				and m.entry_date <=	coalesce(l.end_date, now())
				and m.animal_id = l.animal_id
			left join pasture_entries pe on
				pe.animal_id = m.animal_id
				and pe.entry_date <= m.entry_date
				and m.entry_date <= coalesce(pe.exit_date, now())
			left join pastures p on p.id = pe.pasture_id
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

func (r *LactationRepository) SearchLactatingAnimals(userId string) (*[]entity.SearchEntity, error) {
	query := `
		select
			l.id,
			format(
				'%s (Data de Início: %s)',
				concat_ws(' - ', a.ring_number, a.name), 
				to_char(l.start_date, 'DD/MM/YYYY')
			) as label
		from lactations l
			join animals a on a.id = l.animal_id
				and a.death_date is null
		where l.deleted_at is null
			and l.end_date is null
			and l.user_id = $1
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *LactationRepository) SearchDryAnimals(userId string) (*[]entity.SearchEntity, error) {
	query := `
		select
			a.id,
			concat_ws(' - ', a.ring_number, a.name) as label
		from animals a
		where a.animal_type = 'DAIRY_ANIMAL'
			and a.death_date is null
			and not exists (
				select 1
				from lactations l
				where l.animal_id = a.id
					and l.end_date is null
			)
			and a.deleted_at is null
			and a.user_id = $1
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *LactationRepository) FindLacAnimalsPage(
	filter LactationHistFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[LactationHist], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":     {Field: "cte.animal_order", Order: "asc"},
		"name":             {Field: "cte.name", Order: "asc"},
		"start_date":       {Field: "cte.start_date", Order: "asc"},
		"end_date":         {Field: "coalesce(cte.end_date, '-infinity')", Order: "asc"},
		"calf_birth_date":  {Field: "coalesce(cte.calf_birth_date, -infinity)", Order: "asc"},
		"avg_production":   {Field: "coalesce(cte.avg_production, 0)", Order: "asc"},
		"lac_period":       {Field: "cte.lac_period", Order: "asc"},
		"total_production": {Field: "coalesce(cte.total_production, 0)", Order: "asc"},
		"lac_interval":     {Field: "coalesce(cte.lac_interval, 0)", Order: "asc"},
		"id":               {Field: "cte.id", Order: "asc"},
		"created_at":       {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        with lac_stats as (
            select
                l.id,
                avg(coalesce(m.quantity, 0)) avg_prod,
				max(entry_date) max_date,
				max(coalesce(m.quantity, 0)) peak
            from lactations l
                left join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),

		lac_cte as (
			select
				l.id,
				l.animal_id,
				l.calf_id,
				l.start_date,
				l.end_date,
				extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
				l.observation,
				l.created_at
			from lactations l
			where l.user_id = $1 and l.deleted_at is null
		),

		cte as (
			select
				l.id,
				l.animal_id,
				a.name,
				concat_ws(' - ', a.ring_number, a.name) as animal_name,
				coalesce(nullif(regexp_replace(a.ring_number, '[^0-9]', '', 'g'), '')::int, 0) as animal_order,
				l.calf_id,
				c.birth_date calf_birth_date,
				case
					when l.calf_id is null then 'Sem Bezerro'
					when c.name is not null then format(
						'%s (%s)',
						concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY')),
						concat_ws(' - ', c.ring_number, c.name)
					)
					else concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY'))
				end as calf_info,
				l.start_date,
				l.end_date,
				coalesce(s.avg_prod, 0) as avg_production,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) as lac_period,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) * coalesce(s.avg_prod, 0) as total_production,
				l.lac_interval,
				coalesce(s.peak, 0) as peak,
				l.observation,
				l.created_at
			from lac_cte l
				join animals a on a.id = l.animal_id
				left join lac_stats s on s.id = l.id
				left join animals c on c.id = l.calf_id
				left join animals cm on cm.id = c.mother_id
			where l.end_date is null
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

	whereExpression := repositoriesUtil.GetWhereExpression(filterExpression, cursorExpression)

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

func (r *LactationRepository) GetLacAnimalsFoot(filter LactationHistFilter, userId string) (*LactationHistFoot, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(coalesce(m.quantity, 0)) avg_prod,
				max(entry_date) max_date,
				max(coalesce(m.quantity, 0)) peak
            from lactations l
                left join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),

		lac_cte as (
			select
				l.id,
				l.animal_id,
				l.calf_id,
				l.start_date,
				l.end_date,
				extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
				l.observation,
				l.created_at
			from lactations l
			where l.user_id = $1 and l.deleted_at is null
		),

		cte as (
			select
				l.id,
				l.animal_id,
				a.name,
				concat_ws(' - ', a.ring_number, a.name) as animal_name,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
				l.calf_id,
				c.birth_date calf_birth_date,
				case
					when l.calf_id is null then 'Sem Bezerro'
					when c.name is not null then format(
						'%s (%s)',
						concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY')),
						concat_ws(' - ', c.ring_number, c.name)
					)
					else concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY'))
				end as calf_info,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) lac_period,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod total_production,
				l.lac_interval,
				s.peak,
				l.observation,
				l.created_at
			from lac_cte l
				join animals a on a.id = l.animal_id
				left join lac_stats s on s.id = l.id
				left join animals c on c.id = l.calf_id
				left join animals cm on cm.id = c.mother_id
			where l.end_date is null
		)

		select 
			count(*) as total_lacs,
			avg(lac_period) as avg_lac_period,
			avg(avg_production) as avg_production,
			avg(total_production) as avg_total_production,
			avg(lac_interval) as avg_lac_interval,
			avg(peak) as avg_peak
		from cte
	`

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		query += " where " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[LactationHistFoot](r.DB, query, args...)
}
func (r *LactationRepository) FindDryAnimalsPage(
	filter LactationHistFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[LactationHist], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":     {Field: "cte.animal_order", Order: "asc"},
		"name":             {Field: "cte.name", Order: "asc"},
		"start_date":       {Field: "cte.start_date", Order: "asc"},
		"end_date":         {Field: "coalesce(cte.end_date, '-infinity')", Order: "asc"},
		"calf_birth_date":  {Field: "coalesce(cte.calf_birth_date, -infinity)", Order: "asc"},
		"avg_production":   {Field: "coalesce(cte.avg_production, 0)", Order: "asc"},
		"lac_period":       {Field: "cte.lac_period", Order: "asc"},
		"total_production": {Field: "coalesce(cte.total_production, 0)", Order: "asc"},
		"lac_interval":     {Field: "coalesce(cte.lac_interval, 0)", Order: "asc"},
		"id":               {Field: "cte.id", Order: "asc"},
		"created_at":       {Field: "cte.created_at", Order: "asc"},
	}

	query := `
        with lac_stats as (
            select
                l.id,
                avg(coalesce(m.quantity, 0)) avg_prod,
				max(entry_date) max_date,
				max(coalesce(m.quantity, 0)) peak
            from lactations l
                left join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),

		lac_cte as (
			select distinct on (l.animal_id)
				l.id,
				l.animal_id,
				l.calf_id,
				l.start_date,
				l.end_date,
				extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
				l.observation,
				l.created_at
			from lactations l 
			where l.user_id = $1 and l.deleted_at is null
			order by l.animal_id, l.start_date desc
		),

		cte as (
			select
				coalesce(l.id, a.id) as id,
				coalesce(l.animal_id, a.id) as animal_id,
				coalesce(a.name, '') as name,
				concat_ws(' - ', a.ring_number, a.name) as animal_name,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
				l.calf_id,
				c.birth_date calf_birth_date,
				case
					when l.calf_id is null then 'Sem Bezerro'
					when c.name is not null then format(
						'%s (%s)',
						concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY')),
						concat_ws(' - ', c.ring_number, c.name)
					)
					else concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY'))
				end as calf_info,
				l.start_date,
				l.end_date,
				coalesce(s.avg_prod, 0) as avg_production,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) as lac_period,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) * coalesce(s.avg_prod, 0) as total_production,
				coalesce(l.lac_interval, 0) as lac_interval,
				coalesce(s.peak, 0) as peak,
				l.observation,
				coalesce(l.created_at, a.created_at) as created_at
			from animals a
				left join lac_cte l on a.id = l.animal_id
				left join lac_stats s on s.id = l.id
				left join animals c on c.id = l.calf_id
				left join animals cm on cm.id = c.mother_id
			where a.user_id = $1 
				and a.deleted_at is null
				and a.animal_type = 'DAIRY_ANIMAL'
				and a.is_outside_animal = false
				and not exists (
					select 1
					from lactations lac
					where lac.animal_id = a.id
						and lac.deleted_at is null
						and lac.end_date is null
				)
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

	whereExpression := repositoriesUtil.GetWhereExpression(filterExpression, cursorExpression)

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

func (r *LactationRepository) GetDryAnimalsFoot(filter LactationHistFilter, userId string) (*LactationHistFoot, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                avg(coalesce(m.quantity, 0)) avg_prod,
				max(entry_date) max_date,
				max(coalesce(m.quantity, 0)) peak
            from lactations l
                left join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
            group by 1
        ),

		lac_cte as (
			select distinct on (l.animal_id)
				l.id,
				l.animal_id,
				l.calf_id,
				l.start_date,
				l.end_date,
				extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
				l.observation,
				l.created_at
			from lactations l 
			where l.user_id = $1 and l.deleted_at is null
			order by l.animal_id, l.start_date desc
		),

		cte as (
			select
				l.id,
				l.animal_id,
				a.name,
				concat_ws(' - ', a.ring_number, a.name) as animal_name,
				coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
				l.calf_id,
				c.birth_date calf_birth_date,
				case
					when l.calf_id is null then 'Sem Bezerro'
					when c.name is not null then format(
						'%s (%s)',
						concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY')),
						concat_ws(' - ', c.ring_number, c.name)
					)
					else concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY'))
				end as calf_info,
				l.start_date,
				l.end_date,
				s.avg_prod avg_production,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) lac_period,
				coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod total_production,
				l.lac_interval,
				s.peak,
				l.observation,
				l.created_at
			from animals a
				left join lac_cte l on a.id = l.animal_id
				left join lac_stats s on s.id = l.id
				left join animals c on c.id = l.calf_id
				left join animals cm on cm.id = c.mother_id
			where a.user_id = $1 
				and a.deleted_at is null
				and a.animal_type = 'DAIRY_ANIMAL'
				and a.is_outside_animal = false
				and not exists (
					select 1
					from lactations lac
					where lac.animal_id = a.id
						and lac.deleted_at is null
						and lac.end_date is null
				)
		)
		select 
			count(*) as total_lacs,
			avg(lac_period) as avg_lac_period,
			avg(avg_production) as avg_production,
			avg(total_production) as avg_total_production,
			avg(lac_interval) as avg_lac_interval,
			avg(peak) as avg_peak
		from cte
	`

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		query += " where " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[LactationHistFoot](r.DB, query, args...)
}

func (r *LactationRepository) SearchNewLactationCalf(userId string) (*[]entity.SearchEntity, error) {
	query := `
		select
			a.id,
			concat_ws(' - ', a.ring_number, a.sex, to_char(a.birth_date, 'DD/MM/YYYY')) as label
		from animals a
		where a.death_date is null
			and a.animal_type = 'OFFSPRING'
			and a.deleted_at is null
			and not exists (
				select 1
				from lactations l
				where l.deleted_at is null
					and l.user_id = $1 
					and l.calf_id = a.id
			)
			and a.user_id = $1
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0), a.birth_date
	`
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *LactationRepository) SearchLactationCalf(userId string) (*[]entity.SearchEntity, error) {
	query := `
		select
			a.id,
			case
				when a.name is not null then format(
					'%s (%s)',
					concat_ws(' - ', m.ring_number, a.sex, to_char(a.birth_date, 'DD/MM/YYYY')),
					concat_ws(' - ', a.ring_number, a.name)
				)
				else concat_ws(' - ', m.ring_number, a.sex, to_char(a.birth_date, 'DD/MM/YYYY'))
			end as label
		from animals a
			join animals m on m.id = a.mother_id
				and m.animal_type = 'DAIRY_ANIMAL'
		where a.deleted_at is null
			and a.animal_type <> 'OUTSIDE_ANIMAL'
			and a.user_id = $1
		order by coalesce(regexp_replace(m.ring_number, '[^0-9]', '', 'g')::int, 0), a.birth_date
	`
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *LactationRepository) AddLactation(entry *AddLactationStruct) *apiError.APIError {

	if entry.PastureId != nil {

		pastureEntry := pastureEntries.PastureEntry{
			PastureId: *entry.PastureId,
			AnimalId:  entry.AnimalId,
			EntryDate: entry.StartDate,
			UserId:    entry.UserId,
		}

		pastureRepository := pastureEntries.NewRepository(r.DB)
		err := pastureRepository.TransferEntry(&pastureEntry)
		if err != nil {
			return err
		}

		if entry.CalfId != nil {
			calfEntry := pastureEntries.PastureEntry{
				PastureId: *entry.PastureId,
				AnimalId:  *entry.CalfId,
				EntryDate: entry.StartDate,
				UserId:    entry.UserId,
			}

			err := pastureRepository.TransferCalfEntry(&calfEntry)
			if err != nil {
				return err
			}

		}
	}

	lacEntry := LactationHist{
		CalfId:      entry.CalfId,
		AnimalId:    entry.AnimalId,
		StartDate:   &entry.StartDate,
		EndDate:     entry.EndDate,
		Observation: entry.Observation,
		UserId:      entry.UserId,
	}

	validateErr := validateAddLacation(r.DB, lacEntry)
	if validateErr != nil {
		return validateErr
	}

	insertQuery := `
		insert into lactations (animal_id, calf_id, start_date, user_id)
		values (:animal_id, :calf_id, :start_date, :user_id)
	`

	err := repositoriesUtil.NamedExec(r.DB, insertQuery, &lacEntry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *LactationRepository) EndLactation(entry *AddLactationStruct) *apiError.APIError {

	if entry.PastureId != nil {

		pastureEntry := pastureEntries.PastureEntry{
			PastureId: *entry.PastureId,
			AnimalId:  entry.AnimalId,
			EntryDate: entry.StartDate,
			UserId:    entry.UserId,
		}

		pastureRepository := pastureEntries.NewRepository(r.DB)
		err := pastureRepository.TransferEntry(&pastureEntry)
		if err != nil {
			return err
		}

		if entry.CalfId != nil {
			calfEntry := pastureEntries.PastureEntry{
				PastureId: *entry.PastureId,
				AnimalId:  *entry.CalfId,
				EntryDate: entry.StartDate,
				UserId:    entry.UserId,
			}

			err := pastureRepository.TransferCalfEntry(&calfEntry)
			if err != nil {
				return err
			}

		}
	}

	lacEntry := LactationHist{
		Id:          entry.Id,
		EndDate:     entry.EndDate,
		Observation: entry.Observation,
	}

	validateErr := validateUpdateLacation(r.DB, lacEntry)
	if validateErr != nil {
		return validateErr
	}

	insertQuery := `
		update lactations
		set end_date = :end_date,
			observation = :observation
		where id = :id
	`

	err := repositoriesUtil.NamedExec(r.DB, insertQuery, &lacEntry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *LactationRepository) UpdateLactation(entry *LactationHist) (*LactationHist, *apiError.APIError) {

	validateErr := validateUpdateLacation(r.DB, *entry)
	if validateErr != nil {
		return nil, validateErr
	}

	insertQuery := `
		update lactations
		set calf_id = :calf_id,
			start_date = :start_date,
			end_date = :end_date,
			observation = :observation
		where id = :id and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, insertQuery, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		with lac_stats as (
            select
                avg(coalesce(m.quantity, 0)) as avg_prod,
				max(entry_date) as max_date,
				max(coalesce(m.quantity, 0)) as peak
            from lactations l
                left join milk_entries m on 
                    l.animal_id = m.animal_id
                    and l.start_date <= m.entry_date
                    and coalesce(l.end_date, now()) >= m.entry_date
                    and m.deleted_at is null
                    and m.user_id = $1
			where l.id = $1
        )
		select
			l.id,
			l.animal_id,
			concat_ws(' - ', a.ring_number, a.name) as animal_name,
			l.calf_id,
			c.birth_date calf_birth_date,
			case
				when l.calf_id is null then 'Sem Bezerro'
				when c.name is not null then format(
					'%s (%s)',
					concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY')),
					concat_ws(' - ', c.ring_number, c.name)
				)
				else concat_ws(' - ', cm.ring_number, c.sex, to_char(c.birth_date, 'DD/MM/YYYY'))
			end as calf_info,
			l.start_date,
			l.end_date,
			s.avg_prod as avg_production,
			coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) as lac_period,
			coalesce(extract(days from coalesce(l.end_date, s.max_date) - l.start_date) + 1, 0) * s.avg_prod as total_production,
			extract(days from l.start_date - lag(l.end_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
			s.peak,
			l.observation
		from lactations l
			cross join lac_stats s
			join animals a on a.id = l.animal_id
			left join animals c on c.id = l.calf_id
			left join animals cm on cm.id = c.mother_id
		where l.id = $1
	`

	response, err := repositoriesUtil.GetOne[LactationHist](r.DB, selectQuery, entry.Id)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *LactationRepository) DeleteLactation(id string, userId string) *apiError.APIError {

	tx, err := r.DB.Beginx()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		update lactations
		set deleted_at = now()
		where id = $1 and user_id = $2
	`
	err = repositoriesUtil.ExecTx(tx, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	entriesQuery := `
		update milk_entries m
		set deleted_at = now()
		from lactations l
		where l.id = $1
			and m.animal_id = l.animal_id
			and m.entry_date between l.start_date and coalesce(l.end_date, now());
			and m.user_id = $2
	`
	err = repositoriesUtil.ExecTx(tx, entriesQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}
