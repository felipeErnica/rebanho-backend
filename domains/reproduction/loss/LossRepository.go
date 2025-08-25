package loss

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type LossRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *LossRepository {
	return &LossRepository{db}
}

func (r *LossRepository) GetLossRate(userId string) (*LossRate, error) {
	query := `
        with date_series as (
            select generate_series(min_date, date_trunc('year', now()), interval '1 year') loss_date
            from (
                select min(date_trunc('year', loss_date)) min_date
                from losses 
                where 
                    age(loss_date) < interval '5 years'
                    and deleted_at is null
                    and user_id = $1
            )
        ),
        birth_count as (
            select
                 s.loss_date,
                 count(b.*) birth_count
            from births b
                join animals c on c.id = b.calf_id
                right join date_series s on s.loss_date = date_trunc('year', c.birth_date)
                    and b.deleted_at is null
                    and b.user_id = $1
            group by 1
        ),
        loss_count as (
            select 
                s.loss_date,
                count(l.*) loss_count
            from losses l
                right join date_series s on s.loss_date = date_trunc('year', l.loss_date)
                    and l.deleted_at is null
                    and l.user_id = $1
            group by 1
        )
        select 
            l.loss_date,
            coalesce((loss_count::float / nullif(birth_count, 0)::float)*100, 0) loss_rate
        from loss_count l join birth_count b using(loss_date)
        order by loss_date
    `

	result, err := repositoriesUtil.GetList[LossRateHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	lossHist := *result
	var current, previous, trend float64

	switch lenght := len(lossHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = lossHist[0].LossRate
		previous = 0
		trend = 0
	default:
		current = lossHist[lenght-1].LossRate
		previous = lossHist[lenght-2].LossRate
		trend = util.CalculatePercentageTrend(current, previous)
	}

	lossRate := &LossRate{
		Current: current,
		Trend:   trend,
		Hist:    lossHist,
	}

	return lossRate, nil
}

func (r *LossRepository) GetLossesHist(userId string) (*[]LossNumbersHist, error) {
	query := `
        with date_series as (
            select generate_series(min_date, date_trunc('month', now()), interval '1 month') gen_date
            from (
                select min(date_trunc('month', loss_date)) min_date
                from losses 
                where 
                    age(loss_date) < interval '5 years'
                    and user_id = $1 
                    and deleted_at is null
            )
        ),
        grouped_losses as (
            select
                date_trunc('month', loss_date) loss_date,
                count(*) loss_numbers
            from losses
            where user_id = $1 and deleted_at is null
            group by 1
        )
        select s.gen_date loss_date, coalesce(l.loss_numbers, 0) loss_numbers
        from grouped_losses l right join date_series s on s.gen_date = l.loss_date
        order by s.gen_date
    `
	return repositoriesUtil.GetList[LossNumbersHist](r.DB, query, userId)
}

func (r *LossRepository) GetMostLossesAnimals(userId string) (*[]MostLossesAnimals, error) {
	query := `
        with losses_count as (
            select 
                animal_id, 
                count(*) losses
            from losses l join animals a on a.id = l.animal_id
            where 
                a.death_date is null 
                and l.deleted_at is null 
                and l.user_id = $1
            group by 1
        ),
        birth_count as (
            select
                mother_id as animal_id,
                count(*) births
            from births b join losses_count l on l.animal_id = b.mother_id
            group by 1
        ),
        total_losses as (
            select count(*) total_losses 
            from losses 
            where 
                date_trunc('year', loss_date) = date_trunc('year', now())
                and deleted_at is null
                and user_id = $1
        ),
        total_births as (
            select count(*) total_births
            from births b join animals c on c.id = b.calf_id
            where 
                date_trunc('year', c.birth_date) = date_trunc('year', now())
                and b.deleted_at is null
                and b.user_id = $1
        ),
        total_loss_rate as (
            select (total_losses::float / total_births::float)*100 loss_rate
            from total_births, total_losses
        ),
        cte as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                l.losses,
                (l.losses::float / b.births::float)*100 loss_rate
            from losses_count l
                join birth_count b using(animal_id)
                left join animals a on a.id = l.animal_id
            order by losses*0.6 + (l.losses::float / b.births::float)*0.4 desc
            limit 10
        )
        select 
            cte.*,
            coalesce(((cte.loss_rate / nullif(t.loss_rate, 0)) - 1)*100, 0) rate_comparison
        from cte, total_loss_rate t
    `
	return repositoriesUtil.GetList[MostLossesAnimals](r.DB, query, userId)
}

func (r *LossRepository) FindPage(
    filter LossFilter,
    cursor string,
    sort string,
    order string,
    userId string,
) (*entity.Page[PregnancyLoss], error) {

    query := `
        select
            l.id,
            l.animal_id,
            concat_ws(' - ', a.ring_number, a.name) animal_name,
            regexp_replace(a.ring_number, '[^1-9]', '', 'g') animal_order,
            l.loss_date,
            l.observation,
            l.created_at
        from losses l join animals a on a.id = l.animal_id
    `

	sortMap := map[string]string{
		"loss_date": "coalesce(l.loss_date, '-infinity')",
		"animal_order":    "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)",
		"animal_name":     "a.name",
	}

    whereExpression := "where l.user_id = $1 and l.deleted_at is null"

    filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "l", 2)
    if err != nil {
        return nil, err
    }

    cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
    if err != nil {
        return nil, err
    }

    cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, "l", cursorArgs, nextParam)
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

    query = query + whereExpression + " order by " + sortExpression
    args := []any{userId}
    args = append(args, cursorArgs...)
    filterArgs := repositoriesUtil.GetFilterArgs(filter)
    args = append(args, filterArgs...)
    return repositoriesUtil.GetPage[PregnancyLoss](r.DB, query, sort, 100, args...)
}

func (r *LossRepository) GetPageFoot(filter LossFilter, userId string) (*LossFooter, error) {
    query := "select count(*) totals from losses l"
    whereExpression := " where l.user_id = $1 and l.deleted_at is null"

    filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "l", 2)
    if err != nil {
        return nil, err
    }

    if filterExpression != "" {
        whereExpression += " and " + filterExpression
    }

    query += whereExpression
    args := []any{userId}
    filterArgs := repositoriesUtil.GetFilterArgs(filter)
    args = append(args, filterArgs...)
    return repositoriesUtil.GetOne[LossFooter](r.DB, query, args...)
}
