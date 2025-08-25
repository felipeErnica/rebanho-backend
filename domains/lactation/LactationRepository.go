package lactation

import (
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
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total
            from lactations l
                join lac_stats s using (id)
                join animals a on a.id = l.animal_id
        ),
        cte as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total
            from lac_tbl l join animals a on a.id = l.animal_id
            group by 1
        ),
        cte_stats as (
            select
                max(lac_num) max_lac_num,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total
            from cte
        ),
        results as (
            select
                cte.*,
                ((cte.avg_period / s.avg_period) - 1)*100 period_rate,
                ((cte.avg_prod / s.avg_prod) - 1)*100 prod_rate,
                ((cte.avg_total / s.avg_total) - 1)*100 total_rate
            from cte, cte_stats s
            where cte.lac_num > 1
            order by (cte.lac_num / s.max_lac_num)*0.4 + (cte.avg_total / s.avg_total)*0.6  desc
            limit 10
        )
        select * from results order by lac_num desc, avg_total desc
    `
    return repositoriesUtil.GetList[AnimalsRating](r.DB, query, userId)
}

func (r *LactationRepository) GetWorstAnimals(userId string) (*[]AnimalsRating, error) {
	query := `
        with lac_stats as (
            select
                l.id,
                max(entry_date) max_date,
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
                (extract(days from l.end_date - l.start_date) + 1)*s.avg_prod total
            from lactations l
                join lac_stats s using (id)
                join animals a on a.id = l.animal_id
        ),
        cte as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) lac_num,
                avg(l.period) avg_period,
                avg(l.avg_prod) avg_prod,
                avg(l.total) avg_total
            from lac_tbl l join animals a on a.id = l.animal_id
            group by 1
        ),
        cte_stats as (
            select
                max(lac_num) max_lac_num,
                avg(avg_period) avg_period,
                avg(avg_prod) avg_prod,
                avg(avg_total) avg_total
            from cte
        ),
        results as (
            select
                cte.*,
                ((cte.avg_period / s.avg_period) - 1)*100 period_rate,
                ((cte.avg_prod / s.avg_prod) - 1)*100 prod_rate,
                ((cte.avg_total / s.avg_total) - 1)*100 total_rate
            from cte, cte_stats s
            where cte.lac_num > 1
            order by (cte.lac_num / s.max_lac_num)*0.4 + (1 - (cte.avg_total / s.avg_total))*0.6 desc
            limit 10
        )
        select * from results order by lac_num desc, avg_total
    `
    return repositoriesUtil.GetList[AnimalsRating](r.DB, query, userId)
}
