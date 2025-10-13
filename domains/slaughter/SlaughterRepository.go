package slaughter

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type SlaughterRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *SlaughterRepository {
	return &SlaughterRepository{db}
}

func (r *SlaughterRepository) GetLastPerformance(userId string) (*PerformanceRateCard, error) {
	query := `
		with cte as (
			select 
				entry_date,
				avg((dead_weight / nullif(weight * (1 - discount_rate), 0))*100) performance_rate
			from slaughter_entries s
			where user_id = $1 and deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
	`

	result, err := repositoriesUtil.GetList[PerformanceRateHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	hist := *result
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[0].PerformanceRate
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].PerformanceRate
		previous = hist[lenght-2].PerformanceRate
		trend = util.CalculatePercentageTrend(current, previous)
	}

	card := &PerformanceRateCard{
		Current: current,
		Trend:   trend,
		Hist:    hist,
	}

	return card, nil
}

func (r *SlaughterRepository) GetLastAverageWeight(userId string) (*AverageWeightCard, error) {
	query := `
		with cte as (
			select 
				entry_date,
				avg(weight) avg_weight
			from slaughter_entries 
			where user_id = $1 and deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
	`

	result, err := repositoriesUtil.GetList[AverageWeightHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	hist := *result
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[0].AverageWeight
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].AverageWeight
		previous = hist[lenght-2].AverageWeight
		trend = util.CalculatePercentageTrend(current, previous)
	}

	card := &AverageWeightCard{
		Current: current,
		Trend:   trend,
		Hist:    hist,
	}

	return card, nil
}

func (r *SlaughterRepository) GetLastDeadWeight(userId string) (*AverageWeightCard, error) {
	query := `
		with cte as (
			select 
				entry_date,
				avg(dead_weight) avg_weight
			from slaughter_entries 
			where user_id = $1 and deleted_at is null
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
	`

	result, err := repositoriesUtil.GetList[AverageWeightHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	hist := *result
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[0].AverageWeight
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].AverageWeight
		previous = hist[lenght-2].AverageWeight
		trend = util.CalculatePercentageTrend(current, previous)
	}

	card := &AverageWeightCard{
		Current: current,
		Trend:   trend,
		Hist:    hist,
	}

	return card, nil
}

func (r *SlaughterRepository) GetWeightHist(userId string) (*[]WeightHist, error) {
	query := `
		with cte as (
			select 
				date_trunc('month', s.entry_date) entry_date,
				avg(weight) avg_weight,
				avg(dead_weight) dead_weight
			from slaughter_entries s
			where s.user_id = $1 and s.deleted_at is null
			group by 1
			order by 1 desc
			limit 60
		)
		select * from cte order by 1
	`
	return repositoriesUtil.GetList[WeightHist](r.DB, query, userId)
}

func (r *SlaughterRepository) GetRateHist(userId string) (*[]RateHist, error) {
	query := `
		with cte as (
			select 
				date_trunc('month', s.entry_date) entry_date,
				avg((s.dead_weight / nullif(weight * (1 - s.discount_rate), 0))*100) avg_rate
			from slaughter_entries s
			where s.user_id = $1 and s.deleted_at is null
			group by 1
			order by 1 desc
			limit 60
		)
		select * from cte order by 1
	`
	return repositoriesUtil.GetList[RateHist](r.DB, query, userId)
}
func (r *SlaughterRepository) GetBestFathers(userId string) (*[]TableRatings, error) {
	query := `
		with cte as (
			select 
				a.father_id,
				count(animal_id) animals_number,
				avg(weight) avg_weight,
				avg((s.dead_weight / nullif(weight * (1 - s.discount_rate), 0))*100) avg_rate
			from slaughter_entries s join animals a on a.id = s.animal_id
			where s.user_id = $1 and s.deleted_at is null
			group by 1
		),
		gn_stats as (
			select 
				avg(avg_weight) gn_avg_weight,
				avg(avg_rate) gn_avg_rate
			from cte
			where animals_number >= 10
		)
		select 
			concat_ws(' - ', a.ring_number, a.name) name,
			c.avg_weight,
			((c.avg_weight / nullif(s.gn_avg_weight, 0)) - 1) * 100 weight_comparison,
			c.avg_rate,
			((c.avg_rate / nullif(s.gn_avg_rate, 0)) - 1) * 100 rate_comparison,
			c.animals_number
		from gn_stats s, cte c join animals a on a.id = c.father_id
		where animals_number >= 10
		order by avg_rate desc
		limit 10
	`
	return repositoriesUtil.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetBestMothers(userId string) (*[]TableRatings, error) {
	query := `
		with cte as (
			select 
				a.mother_id,
				count(animal_id) animals_number,
				avg(weight) avg_weight,
				avg((s.dead_weight / nullif(weight * (1 - s.discount_rate), 0))*100) avg_rate
			from slaughter_entries s join animals a on a.id = s.animal_id
			where s.user_id = $1 and s.deleted_at is null
			group by 1
		),
		gn_stats as (
			select 
				avg(avg_weight) gn_avg_weight,
				avg(avg_rate) gn_avg_rate
			from cte
			where animals_number >= 3
		)
		select 
			concat_ws(' - ', a.ring_number, a.name) name,
			c.avg_weight,
			((c.avg_weight / nullif(s.gn_avg_weight, 0)) - 1) * 100 weight_comparison,
			c.avg_rate,
			((c.avg_rate / nullif(s.gn_avg_rate, 0)) - 1) * 100 rate_comparison,
			c.animals_number
		from gn_stats s, cte c join animals a on a.id = c.mother_id
		where c.animals_number >= 3
		order by c.avg_rate desc
		limit 10
	`
	return repositoriesUtil.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetBestSlaughterhouses(userId string) (*[]TableRatings, error) {
	query := `
		with cte as (
			select 
				slaughterhouse_id,
				count(animal_id) animals_number,
				avg(weight) avg_weight,
				avg((dead_weight / nullif(weight * (1 - discount_rate), 0))*100) avg_rate
			from slaughter_entries 
			where user_id = $1 and deleted_at is null
			group by 1
		),
		gn_stats as (
			select 
				avg(avg_weight) gn_avg_weight,
				avg(avg_rate) gn_avg_rate
			from cte
		)
		select 
			s.name,
			c.avg_weight,
			((c.avg_weight / nullif(g.gn_avg_weight, 0)) - 1) * 100 weight_comparison,
			c.avg_rate,
			((c.avg_rate / nullif(g.gn_avg_rate, 0)) - 1) * 100 rate_comparison,
			c.animals_number
		from gn_stats g, cte c join slaughterhouses s on s.id = c.slaughterhouse_id
		where c.avg_rate >= 20
		order by c.avg_rate desc
		limit 10
	`
	return repositoriesUtil.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastGroups(userId string) (*[]SlaughterGroup, error) {
	query := `
		with cte as (
			select 
				entry_date,
				slaughterhouse_id,
				count(animal_id) animals_number,
				avg(dead_weight) avg_weight,
				avg((dead_weight / nullif(weight * (1 - discount_rate), 0))*100) avg_rate
			from slaughter_entries 
			where user_id = $1 and deleted_at is null
			group by 1, 2
		)
		select 
			c.entry_date,
			s.name slaughterhouse,
			c.avg_weight,
			coalesce(((c.avg_weight / lag(nullif(c.avg_weight, 0)) over (order by c.entry_date)) - 1) * 100, 0) weight_variation,
			c.avg_rate,
			coalesce(((c.avg_rate / lag(nullif(c.avg_rate, 0)) over (order by c.entry_date)) - 1) * 100, 0) rate_variation,
			c.animals_number
		from cte c join slaughterhouses s on s.id = c.slaughterhouse_id
		order by c.entry_date desc
		limit 15
	`
	return repositoriesUtil.GetList[SlaughterGroup](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastEntries(userId string) (*[]SlaughterEntry, error) {
	query := `
		with last_date as (
			select max(entry_date) max_date 
			from slaughter_entries
			where user_id = $1 and deleted_at is null
		)
		select 
			concat_ws(
				' - ', 
				a.ring_number, 
				coalesce(a.name, concat_ws(' - ', a.sex, to_char(a.birth_date, 'DD/MM/YYYY')))
			) animal_name,
			h.name slaughterhouse,
			s.entry_date,
			s.weight,
			s.discount_rate,
			(s.weight*(1 - s.discount_rate)) discount_weight,
			s.dead_weight,
			(s.dead_weight /(s.weight*(1 - s.discount_rate))) * 100 performance_rate
		from last_date l, slaughter_entries s
			join animals a on a.id = s.animal_id
			join slaughterhouses h on h.id = s.slaughterhouse_id
		where s.entry_date = l.max_date
			and s.user_id = $1 
			and s.deleted_at is null
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) 
	`
	return repositoriesUtil.GetList[SlaughterEntry](r.DB, query, userId)
}

func (r *SlaughterRepository) FindEntriesPage(
	sort string,
	order string,
	cursor string,
	filter SlaughterEntryFilter,
	userId string,
) (*entity.Page[SlaughterEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"entry_date":       {Field: "s.entry_date", Order: "desc"},
		"animal_order":     {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "coalesce(a.name, '')", Order: "asc"},
		"birth_date":       {Field: "coalesce(a.birth_date, '-infinity')", Order: "desc"},
		"weight":           {Field: "s.weight", Order: "asc"},
		"dead_weight":      {Field: "s.dead_weight", Order: "asc"},
		"performance_rate": {Field: "coalesce(s.dead_weight / nullif(s.weight*(1 - s.discount_rate), 0) * 100, 0)", Order: "asc"},
		"id":               {Field: "s.id", Order: "asc"},
		"created_at":       {Field: "s.created_at", Order: "asc"},
	}

	query := `
		select 
			s.id,
			s.animal_id,
			s.slaughterhouse_id,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			concat_ws(
				' - ', 
				a.ring_number, 
				coalesce(a.name, a.sex),
				to_char(a.birth_date, 'DD/MM/YYYY')
			) animal_name,
			a.birth_date,
			concat_ws(' - ', f.ring_number, f.name) father_name,
			concat_ws(' - ', m.ring_number, m.name) mother_name,
			h.name slaughterhouse,
			s.entry_date,
			s.discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) discount_weight,
			s.dead_weight,
			coalesce(s.dead_weight / nullif(s.weight*(1 - s.discount_rate), 0) * 100, 0) performance_rate,
			s.created_at
		from slaughter_entries s
			join animals a on a.id = s.animal_id
			join slaughterhouses h on h.id = s.slaughterhouse_id
			join animals f on f.id = a.father_id
			join animals m on m.id = a.mother_id
	`

	whereExpression := "where s.user_id = $1 and s.deleted_at is null"

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "s", 2)
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
	orderExpression := " order by " + sortExpression
	query += whereExpression + orderExpression

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[SlaughterEntry](r.DB, query, sort, 200, args...)
}

func (r *SlaughterRepository) GetEntriesPageFoot(filter SlaughterEntryFilter, userId string) (*SlaughterFoot, error) {

	query := `
		select 
			count(s.animal_id) animals_number,
			avg(s.weight) avg_weight,
			avg(s.dead_weight) avg_dead_weight,
			avg((s.dead_weight / nullif(weight * (1 - s.discount_rate), 0))*100) avg_rate
		from slaughter_entries s
	`

	whereExpression := " where s.user_id = $1 and s.deleted_at is null"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "s", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query += whereExpression

	return repositoriesUtil.GetOne[SlaughterFoot](r.DB, query, args...)
}

func (r *SlaughterRepository) FindGroups(userId string) (*[]SlaughterGroup, error) {
	query := `
		with cte as (
			select 
				entry_date,
				slaughterhouse_id,
				count(animal_id) animals_number,
				avg(weight) avg_weight,
				avg((dead_weight / nullif(weight * (1 - discount_rate), 0))*100) avg_rate
			from slaughter_entries 
			where user_id = $1 and deleted_at is null
			group by 1, 2
		)
		select 
			c.entry_date,
			s.name slaughterhouse,
			c.avg_weight,
			((c.avg_weight / lag(nullif(c.avg_weight, 0)) over (order by c.entry_date)) - 1) * 100 weight_variation,
			c.avg_rate,
			((c.avg_rate / lag(nullif(c.avg_rate, 0)) over (order by c.entry_date)) - 1) * 100 rate_variation,
			c.animals_number
		from cte c join slaughterhouses s on s.id = c.slaughterhouse_id
		order by c.entry_date desc
	`

	return repositoriesUtil.GetList[SlaughterGroup](r.DB, query, userId)
}

func (r *SlaughterRepository) FindEntriesByDate(entryDate time.Time, userId string) (*[]SlaughterEntry, error) {
	query := `
		select 
			s.id,
			s.animal_id,
			s.slaughterhouse_id,
			concat_ws(
				' - ', 
				a.ring_number, 
				coalesce(a.name, concat_ws(' - ', a.sex, to_char(a.birth_date, 'DD/MM/YYYY'))) 
			) animal_name,
			s.weight,
			s.discount_rate,
			s.weight * (1 - s.discount_rate) discount_weight,
			s.dead_weight,
			coalesce(s.dead_weight / nullif(s.weight * (1 - s.discount_rate), 0) * 100) performance_rate
		from slaughter_entries s join animals a on a.id = s.animal_id
		where s.entry_date = $1
			and s.user_id = $2 
			and s.deleted_at is null
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) 
	`
	return repositoriesUtil.GetList[SlaughterEntry](r.DB, query, entryDate, userId)
}
