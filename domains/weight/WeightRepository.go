package weight

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type WeightRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *WeightRepository {
	return &WeightRepository{db}
}

func (r *WeightRepository) GetYearWeightGain(userId string) (*CardWeightGain, error) {
	query := `
		with weight_gain as (
			select 
				date_trunc('year', w.entry_date) entry_date,
				(w.weight - 38) / extract(days from w.entry_date - a.birth_date) daily_gain
			from weight_entries w join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
		),
		cte as (
			select 
				entry_date,
				avg(daily_gain) average_gain
			from weight_gain
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
	`
	result, err := repositoriesUtil.GetList[AverageWeightGain](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	gainHist := *result
	var current, previous, trend float64

	switch lenght := len(gainHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = gainHist[0].AverageGain
		previous = 0
		trend = 0
	default:
		current = gainHist[lenght-1].AverageGain
		previous = gainHist[lenght-2].AverageGain
		trend = util.CalculatePercentageTrend(current, previous)
	}

	cardInfo := &CardWeightGain{
		Current: current,
		Trend:   trend,
		Hist:    gainHist,
	}
	return cardInfo, nil
}

func (r *WeightRepository) GetYearWeight(userId string) (*CardWeight, error) {
	query := `
		with cte as (
			select 
				date_trunc('year', entry_date) entry_date,
				avg(weight) average_weight
			from weight_entries
			where user_id = $1 and deleted_at is null
			group by 1 
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
	`
	result, err := repositoriesUtil.GetList[AverageWeight](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	gainHist := *result
	var current, previous, trend float64

	switch lenght := len(gainHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = gainHist[0].AverageWeight
		previous = 0
		trend = 0
	default:
		current = gainHist[lenght-1].AverageWeight
		previous = gainHist[lenght-2].AverageWeight
		trend = util.CalculatePercentageTrend(current, previous)
	}

	cardInfo := &CardWeight{
		Current: current,
		Trend:   trend,
		Hist:    gainHist,
	}
	return cardInfo, nil
}

func (r *WeightRepository) GetLastWeightGain(userId string) (*CardWeightGain, error) {
	query := `
		with weight_gain as (
			select 
				entry_date,
				(w.weight - 38) / extract(days from w.entry_date - a.birth_date) daily_gain
			from weight_entries w join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
		),
		cte as (
			select 
				entry_date,
				avg(daily_gain) average_gain
			from weight_gain
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
	`
	result, err := repositoriesUtil.GetList[AverageWeightGain](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	gainHist := *result
	var current, previous, trend float64

	switch lenght := len(gainHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = gainHist[0].AverageGain
		previous = 0
		trend = 0
	default:
		current = gainHist[lenght-1].AverageGain
		previous = gainHist[lenght-2].AverageGain
		trend = util.CalculatePercentageTrend(current, previous)
	}

	cardInfo := &CardWeightGain{
		Current: current,
		Trend:   trend,
		Hist:    gainHist,
	}
	return cardInfo, nil
}

func (r *WeightRepository) GetLastWeight(userId string) (*CardWeight, error) {
	query := `
		with cte as (
			select 
				entry_date,
				avg(weight) average_weight
			from weight_entries
			where user_id = $1 and deleted_at is null
			group by 1 
			order by 1 desc
			limit 10
		)
		select * from cte order by entry_date
	`
	result, err := repositoriesUtil.GetList[AverageWeight](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	gainHist := *result
	var current, previous, trend float64

	switch lenght := len(gainHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = gainHist[0].AverageWeight
		previous = 0
		trend = 0
	default:
		current = gainHist[lenght-1].AverageWeight
		previous = gainHist[lenght-2].AverageWeight
		trend = util.CalculatePercentageTrend(current, previous)
	}

	cardInfo := &CardWeight{
		Current: current,
		Trend:   trend,
		Hist:    gainHist,
	}
	return cardInfo, nil
}

func (r *WeightRepository) GetLastEntries(userId string) (*[]WeightEntry, error) {
	query := `
		with last_date as (
			select max(entry_date) entry_date from weight_entries 
			where user_id = $1 and deleted_at is null
		)
		select 
			w.id,
			w.animal_id,
			concat_ws(' - ', a.ring_number, coalesce(a.name, to_char(a.birth_date, 'DD/MM/YYYY')), a.sex) animal_name,
			a.birth_date,
			w.entry_date,
			w.weight,
			coalesce(w.weight - lag(w.weight) over (partition by w.animal_id order by w.entry_date), 0) weight_variation,
			(w.weight - 38) / extract(days from w.entry_date - a.birth_date) weight_gain
		from last_date l, weight_entries w join animals a on a.id = w.animal_id
		where w.entry_date = l.entry_date
			and w.user_id = $1 
			and w.deleted_at is null
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return repositoriesUtil.GetList[WeightEntry](r.DB, query, userId)
}

func (r *WeightRepository) GetLastGroups(userId string) (*[]WeightGroup, error) {
	query := `
		with entries as (
			select 
				w.entry_date,
				w.weight,
				(w.weight - 38) / extract(days from w.entry_date - a.birth_date) weight_gain
			from  weight_entries w join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
		),
		cte as (
			select
				w.entry_date,
				count(w.weight) animals_number,
				avg(w.weight) average_weight,
				avg(w.weight_gain) average_gain
			from entries w
			group by 1
		)
		select
			c.*,
			coalesce(
				((c.average_gain / lag(c.average_gain) over (order by c.entry_date)) - 1) * 100
			, 0) gain_variation,
			coalesce(
				((c.average_weight / lag(c.average_weight) over (order by c.entry_date)) - 1) * 100
			, 0) weight_variation
		from cte c
		order by entry_date desc
		limit 5
	`
	return repositoriesUtil.GetList[WeightGroup](r.DB, query, userId)
}

func (r *WeightRepository) GetBestFathers(userId string) (*[]AnimalRating, error) {
	query := `
		with gain_tbl as (
			select 
				w.animal_id,
				(w.weight - 38)/extract(days from w.entry_date - a.birth_date) weight_gain
			from weight_entries w join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
		),
		animal_tbl as (
			select 
				animal_id,
				avg(weight_gain) animal_gain
			from gain_tbl
			group by 1
		),
		father_tbl as (
			select
				m.id mother_id,
				count(t.animal_id) children_number,
				avg(t.animal_gain) avg_gain
			from animal_tbl t
				join animals a on a.id = t.animal_id
				join animals m on m.id = a.father_id
			group by 1
		),
		stats as (
			select 
				avg(avg_gain) gn_avg_gain,
				avg(children_number) avg_children,
				stddev(avg_gain) dev_gain,
				stddev(children_number) dev_children
			from father_tbl
		),
		z_score as (
			select
				t.*,
				(t.children_number - s.avg_children) / s.dev_children z_children,
				(t.avg_gain - s.gn_avg_gain) / s.dev_gain z_gain
			from father_tbl t, stats s
		)
		select 
			concat_ws(' - ', m.ring_number, m.name) animal_name,
			t.avg_gain,
			((t.avg_gain / s.gn_avg_gain) - 1) * 100 gain_trend,
			t.children_number
		from stats s, z_score t join animals m on m.id = t.mother_id
		order by (
			case
				when z_children < 0 then z_gain*0.3 + z_children*0.7 
				else z_gain + z_children
			end 
		) desc
		limit 10
	`
	return repositoriesUtil.GetList[AnimalRating](r.DB, query, userId)
}

func (r *WeightRepository) GetBestMothers(userId string) (*[]AnimalRating, error) {
	query := `
		with gain_tbl as (
			select 
				w.animal_id,
				(w.weight - 38)/extract(days from w.entry_date - a.birth_date) weight_gain
			from weight_entries w join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
		),
		animal_tbl as (
			select 
				animal_id,
				avg(weight_gain) animal_gain
			from gain_tbl
			group by 1
		),
		mother_tbl as (
			select
				m.id mother_id,
				count(t.animal_id) children_number,
				avg(t.animal_gain) avg_gain
			from animal_tbl t
				join animals a on a.id = t.animal_id
				join animals m on m.id = a.mother_id
			group by 1
		),
		stats as (
			select 
				avg(avg_gain) gn_avg_gain,
				avg(children_number) avg_children,
				stddev(avg_gain) dev_gain,
				stddev(children_number) dev_children
			from mother_tbl
		),
		z_score as (
			select
				t.*,
				(t.children_number - s.avg_children) / s.dev_children z_children,
				(t.avg_gain - s.gn_avg_gain) / s.dev_gain z_gain
			from mother_tbl t, stats s
		)
		select 
			concat_ws(' - ', m.ring_number, m.name) animal_name,
			t.avg_gain,
			((t.avg_gain / s.gn_avg_gain) - 1) * 100 gain_trend,
			t.children_number
		from stats s, z_score t join animals m on m.id = t.mother_id
		order by (
			case
				when z_children < 0 then z_gain*0.3 + z_children*0.7 
				else z_gain + z_children
			end 
		) desc
		limit 10
	`
	return repositoriesUtil.GetList[AnimalRating](r.DB, query, userId)
}

func (r *WeightRepository) FindEntriesPage(
	sort string,
	order string,
	cursor string,
	filter WeightFilter,
	userId string,
) (*entity.Page[WeightEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"entry_date":   {Field: "w.entry_date", Order: "desc"},
		"animal_order": {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":  {Field: "a.name", Order: "asc"},
		"birth_date":   {Field: "coalesce(a.birth_date, '-infinite')", Order: "asc"},
		"id":           {Field: "w.id", Order: "asc"},
		"created_at":   {Field: "w.created_at", Order: "asc"},
	}

	query := `
		select
			w.id,
			w.animal_id,
			concat_ws(' - ', a.ring_number, coalesce(a.name, to_char(a.birth_date, 'DD/MM/YYYY')), a.sex) animal_name,
			a.birth_date,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			w.entry_date,
			w.weight,
			coalesce((w.weight - 38) / nullif(extract(days from w.entry_date - a.birth_date), 0), 0) weight_gain,
			coalesce(w.weight - lag(w.weight) over (partition by w.animal_id order by w.entry_date), 0) weight_variation,
			w.created_at
		from weight_entries w join animals a on a.id = w.animal_id
	`

	whereExpression := " where w.user_id = $1 and w.deleted_at is null"

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "w", 2)
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

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	query += whereExpression + orderExpression
	return repositoriesUtil.GetPage[WeightEntry](r.DB, query, sort, 200, args...)
}

func (r *WeightRepository) FindGroups(userId string) (*[]WeightGroup, error) {
	query := `
		with entries as (
			select 
				w.entry_date,
				w.weight,
				(w.weight - 38) / extract(days from w.entry_date - a.birth_date) weight_gain
			from  weight_entries w join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
		),
		cte as (
			select
				w.entry_date,
				count(w.weight) animals_number,
				avg(w.weight) average_weight,
				avg(w.weight_gain) average_gain
			from entries w
			group by 1
			order by 1
		)
		select
			c.*,
			coalesce(
				((c.average_gain / lag(c.average_gain) over (order by c.entry_date)) - 1) * 100
			, 0) gain_variation,
			coalesce(
				((c.average_weight / lag(c.average_weight) over (order by c.entry_date)) - 1) * 100
			, 0) weight_variation
		from cte c
	`
	return repositoriesUtil.GetList[WeightGroup](r.DB, query, userId)
}

func (r *WeightRepository) FindEntriesByDate(entryDate time.Time, userId string) (*[]WeightEntry, error) {
	query := `
		select 
			w.id,
			w.animal_id,
			concat_ws(' - ', a.ring_number, a.name, a.sex) animal_name,
			a.birth_date,
			w.entry_date,
			w.weight,
			coalesce(w.weight - lag(w.weight) over (partition by w.animal_id order by w.entry_date), 0) weight_variation,
			(w.weight - 38) / extract(days from w.entry_date - a.birth_date) weight_gain
		from weight_entries w join animals a on a.id = w.animal_id
		where w.entry_date = $1
			and w.user_id = $2 
			and w.deleted_at is null
	`
	return repositoriesUtil.GetList[WeightEntry](r.DB, query, entryDate, userId)
}
