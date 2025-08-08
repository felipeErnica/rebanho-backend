package insemination

import (
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type InseminationRepository struct {
	Db *sqlx.DB
}

func NewEntryRepository(db *sqlx.DB) *InseminationRepository {
	return &InseminationRepository{db}
}

func (r *InseminationRepository) GetBirthRateStats(userId string) (*BirthRateStats, error) {
	query := `
        with totals as (
            select distinct
                date_trunc('month', g.insemination_date) date_month,
                count(*) over (order by date_trunc('month', g.insemination_date)) total,
                count(*) filter (where status = 'SUCCESS') over (order by date_trunc('month', g.insemination_date)) successful
            from insemination_entries i
                left join insemination_groups g on g.id = i.group_id 
            where 
                age(g.insemination_date) > interval '11 months'
                and i.user_id = $1
            order by 1
            limit 10
        )
        select 
            date_month,
            (successful::float/total::float)*100 birth_rate
        from totals
    `
	result, err := repositoriesUtil.GetList[BirthRateHist](r.Db, query, userId)
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

func (r *InseminationRepository) GetInseminationStats(userId string) (*[]InseminationHist, error) {
	query := `
        with totals as (
            select
                date_trunc('month', g.insemination_date) date_month,
                count(*) total,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries i
                left join insemination_groups g on g.id = i.group_id 
            where 
                age(g.insemination_date) > interval '11 months'
                and i.user_id = $1
            group by 1
            order by 1
            limit 60
        )
        select 
            date_month,
            total,
            (successful::float/total::float)*100 birth_rate
        from totals
    `
	return repositoriesUtil.GetList[InseminationHist](r.Db, query, userId)
}

func (r *InseminationRepository) GetBestBull(userId string) (*[]InseminationBulls, error) {
	query := `
        with totals as (
            select
                bull.name bull_name,
                count(*) total,
                count(*) filter (where status = 'SUCCESS') successful
            from insemination_entries i
                left join insemination_groups g on g.id = i.group_id 
                left join animals bull on g.bull_id = bull.id 
            where 
                age(g.insemination_date) > interval '11 months'
                and i.user_id = $1
            group by 1
        ),
        cte as (
            select 
                bull_name,
                total,
                (successful::float/total::float)*100 birth_rate
            from totals
        ),
        average as (select avg(birth_rate) average_rate from cte) 
        select
            cte.*,
            ((cte.birth_rate / avg.average_rate) - 1)*100 comparison_rate
        from cte, average avg
        order by cte.birth_rate desc
    `
	return repositoriesUtil.GetList[InseminationBulls](r.Db, query, userId)
}
