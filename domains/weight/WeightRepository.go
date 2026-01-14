package weight

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
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

func (r *WeightRepository) GetWeightGainHist(userId string) (*[]AverageWeightGain, error) {
	query := `
		with weight_gain as (
			select 
				date_trunc('month', w.entry_date) entry_date,
				(w.weight - coalesce(lag(w.weight) over win, 38)) / 
				extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)) as daily_gain
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
		),
		cte as (
			select 
				entry_date,
				avg(daily_gain) average_gain
			from weight_gain
			group by 1
			order by 1 desc
			limit 60
		)
		select * from cte order by entry_date
	`
	return repositoriesUtil.GetList[AverageWeightGain](r.DB, query, userId)
}

func (r *WeightRepository) GetWeightHist(userId string) (*[]AverageWeight, error) {
	query := `
		with cte as (
			select 
				date_trunc('month', entry_date) entry_date,
				avg(weight) average_weight
			from weight_entries
			where user_id = $1 and deleted_at is null
			group by 1
			order by 1 desc
			limit 60
		)
		select * from cte order by entry_date
	`
	return repositoriesUtil.GetList[AverageWeight](r.DB, query, userId)
}
func (r *WeightRepository) GetLastWeightGain(userId string) (*CardWeightGain, error) {
	query := `
		with weight_gain as (
			select 
				entry_date,
				(w.weight - coalesce(lag(w.weight) over win, 38)) / 
				nullif(extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)), 0) as daily_gain
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
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
		),
		stats as (
			select
				w.id,
				coalesce(w.weight - lag(w.weight) over win, 0) weight_variation,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) / 
					extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)), 
					0
				) as weight_gain 
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
		)
		select 
			w.id,
			w.animal_id,
			concat_ws(
				' - ',
				a.ring_number,
				coalesce(a.name, a.sex),
				to_char(a.birth_date, 'DD/MM/YYYY')
			) as animal_info,
			a.birth_date,
			w.entry_date,
			w.weight,
			s.weight_variation,
			s.weight_gain
		from last_date l, weight_entries w 
			join stats s on s.id = w.id
			left join animals a on a.id = w.animal_id
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
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) / 
					extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)), 
					0
				) as weight_gain
			from  weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
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
			coalesce( ((c.average_gain / lag(c.average_gain) over win) - 1) * 100, 0) gain_variation,
			coalesce( ((c.average_weight / lag(c.average_weight) over win) - 1) * 100, 0) weight_variation 
		from cte c
		window win as (order by c.entry_date)
		order by entry_date desc
		limit 10
	`
	return repositoriesUtil.GetList[WeightGroup](r.DB, query, userId)
}

func (r *WeightRepository) GetBestFathers(userId string) (*[]AnimalRating, error) {
	query := `
		with gain_tbl as (
			select 
				w.animal_id,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) / 
					extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)), 
					0
				) as weight_gain
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
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
				f.id father_id,
				count(t.animal_id) children_number,
				avg(t.animal_gain) avg_gain
			from animal_tbl t
				left join animals a on a.id = t.animal_id
				left join animals f on f.id = a.father_id
			group by 1
		),
		stats as (
			select avg(weight_gain) gn_avg_gain
			from gain_tbl
		)
		select 
			concat_ws(' - ', f.ring_number, f.name) animal_name,
			t.avg_gain,
			((t.avg_gain / s.gn_avg_gain) - 1) * 100 gain_trend,
			t.children_number
		from stats s, father_tbl t 
			join animals f on f.id = t.father_id
		where t.children_number >= 10
		order by t.avg_gain desc
		limit 10
	`
	return repositoriesUtil.GetList[AnimalRating](r.DB, query, userId)
}

func (r *WeightRepository) GetBestMothers(userId string) (*[]AnimalRating, error) {
	query := `
		with gain_tbl as (
			select 
				w.animal_id,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) / 
					extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)), 
					0
				) as weight_gain
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
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
				left join animals a on a.id = t.animal_id
				left join animals m on m.id = a.mother_id
			group by 1
		),
		stats as (
			select avg(weight_gain) gn_avg_gain
			from gain_tbl
		)
		select 
			concat_ws(' - ', m.ring_number, m.name) animal_name,
			t.avg_gain,
			((t.avg_gain / s.gn_avg_gain) - 1) * 100 gain_trend,
			t.children_number
		from stats s, mother_tbl t 
			left join animals m on m.id = t.mother_id
		where t.children_number >= 5
		order by t.avg_gain desc
		limit 10
	`
	return repositoriesUtil.GetList[AnimalRating](r.DB, query, userId)
}

func (r *WeightRepository) FindEntriesPage(
	sort string,
	order string,
	cursor string,
	filter *WeightFilter,
	userId string,
) (*entity.Page[WeightEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"entry_date":   {Field: "w.entry_date", Order: "desc"},
		"animal_order": {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":  {Field: "coalesce(a.name, '')", Order: "asc"},
		"birth_date":   {Field: "coalesce(a.birth_date, '-infinity')", Order: "asc"},
		"id":           {Field: "w.id", Order: "asc"},
		"created_at":   {Field: "w.created_at", Order: "asc"},
	}

	query := `
		with stats as (
			select
				w.id,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) /
					nullif(extract(day from (w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date))), 0),
					0
				) weight_gain,
				coalesce(w.weight - lag(w.weight) over win, 0) as weight_variation
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
		)
		select
			w.id,
			w.animal_id,
			coalesce(a.name, '') as animal_name,
			concat_ws(
				' - ',
				a.ring_number,
				coalesce(a.name, a.sex),
				to_char(a.birth_date, 'DD/MM/YYYY')
			) as animal_info,
			a.birth_date,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
			concat_ws(' - ', f.ring_number, f.name) as father_name,
			concat_ws(' - ', m.ring_number, m.name) as mother_name,
			s.weight_gain,
			s.weight_variation,
			w.entry_date,
			w.weight,
			w.created_at
		from weight_entries w 
			join stats s on s.id = w.id
			left join animals a on a.id = w.animal_id
			left join animals m on m.id = a.mother_id
			left join animals f on f.id = a.father_id
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

	windowExpression := " window win as (partition by w.animal_id order by w.entry_date)"
	query += whereExpression + windowExpression + orderExpression
	return repositoriesUtil.GetPage[WeightEntry](r.DB, query, sort, 200, args...)
}

func (r *WeightRepository) GetEntriesPageFoot(filter *WeightFilter, userId string) (*WeightFoot, error) {
	query := `
		select
			w.animal_id,
			w.weight,
			coalesce(
				(w.weight - coalesce(lag(w.weight) over win, 38)) /
				nullif(extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)), 0),
				0
			) weight_gain
		from weight_entries w 
			left join animals a on a.id = w.animal_id
	`
	whereExpression := " where w.user_id = $1 and w.deleted_at is null"
	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "w", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	windowExpression := " window win as (partition by w.animal_id order by w.entry_date)"
	query += whereExpression + windowExpression
	query = fmt.Sprintf(`
		with cte as (%s)
		select
			count(*) as animals_num,
			avg(weight) as avg_weight,
			avg(weight_gain) as avg_gain
		from cte
	`, query)

	return repositoriesUtil.GetOne[WeightFoot](r.DB, query, args...)
}

func (r *WeightRepository) FindGroups(userId string, order string) (*[]WeightGroup, error) {
	query := `
		with entries as (
			select 
				w.entry_date,
				w.weight,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) / 
					extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)),
					0 
				) weight_gain
			from  weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
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
	orderExpression := " order by c.entry_date " + order
	query += orderExpression
	return repositoriesUtil.GetList[WeightGroup](r.DB, query, userId)
}

func (r *WeightRepository) FindEntriesByDate(
	entryDate time.Time,
	userId string,
	order string,
	sort string,
) (*[]WeightEntry, error) {

	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order": {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":  {Field: "a.name", Order: "asc"},
		"birth_date":   {Field: "coalesce(a.birth_date, '-infinity')", Order: "desc"},
	}

	query := `
		with stats as (
			select
				w.id,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) /
					nullif(extract(day from (w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date))), 0),
					0
				) weight_gain,
				coalesce(w.weight - lag(w.weight) over win, 0) as weight_variation
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $2 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
		)
		select 
			w.id,
			w.animal_id,
			coalesce(a.name, '') as animal_name,
			concat_ws(
				' - ', 
				a.ring_number, 
				coalesce(a.name, a.sex), 
				to_char(a.birth_date, 'DD/MM/YYYY')
			) as animal_info,
			a.birth_date,
			concat_ws(' - ', f.ring_number, f.name) father_name,
			concat_ws(' - ', m.ring_number, m.name) mother_name,
			w.entry_date,
			w.weight,
			s.weight_variation,
			s.weight_gain
		from weight_entries w 
			left join animals a on a.id = w.animal_id
			left join animals f on f.id = a.father_id
			left join animals m on m.id = a.mother_id
			join stats s on s.id = w.id
		where w.entry_date = $1
			and w.user_id = $2 
			and w.deleted_at is null
		window win as (partition by w.animal_id order by w.entry_date)
	`
	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query += " order by " + sortExpression

	return repositoriesUtil.GetList[WeightEntry](r.DB, query, entryDate, userId)
}

func (r *WeightRepository) GetEntriesByDateFoot(entryDate time.Time, userId string) (*WeightFoot, error) {
	query := `
		with cte as (
			select
				w.animal_id,
				w.weight,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) /
					nullif(extract(days from w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date)), 0),
					0
				) weight_gain
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.entry_date = $1 and w.user_id = $2 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
		) 
		select 
			count(animal_id) animals_num,
			avg(weight) avg_weight,
			avg(weight_gain) avg_gain
		from cte
	`
	return repositoriesUtil.GetOne[WeightFoot](r.DB, query, entryDate, userId)
}

func (r *WeightRepository) Delete(id string, userId string) *apiError.APIError {
	query := `
		update weight_entries
		set deleted_at = now()
		where id = $1 and user_id = $2
	`
	err := repositoriesUtil.Exec(r.DB, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *WeightRepository) Update(entry *WeightEntrySave) (*WeightEntry, *apiError.APIError) {

	validateErr := validateUpdate(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		update weight_entries
		set entry_date = :entry_date,
			weight = :weight
		where id = :id and user_id = :user_id
	`
	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		with stats as (
			select
				w.id,
				coalesce(
					(w.weight - coalesce(lag(w.weight) over win, 38)) /
					nullif(extract(day from (w.entry_date - coalesce(lag(w.entry_date) over win, a.birth_date))), 0),
					0
				) weight_gain,
				coalesce(w.weight - lag(w.weight) over win, 0) as weight_variation
			from weight_entries w 
				left join animals a on a.id = w.animal_id
			where w.user_id = $1 and w.deleted_at is null
			window win as (partition by w.animal_id order by w.entry_date)
		)
		select
			w.id,
			w.animal_id,
			concat_ws(
				' - ',
				a.ring_number,
				coalesce(a.name, a.sex),
				to_char(a.birth_date, 'DD/MM/YYYY')
			) as animal_info,
			a.birth_date,
			concat_ws(' - ', f.ring_number, f.name) as father_name,
			concat_ws(' - ', m.ring_number, m.name) as mother_name,
			s.weight_gain,
			s.weight_variation,
			w.entry_date,
			w.weight
		from weight_entries w 
			join stats s on s.id = w.id
			left join animals a on a.id = w.animal_id
			left join animals m on m.id = a.mother_id
			left join animals f on f.id = a.father_id
		where id = :id and user_id = :user_id
	`

	response, err := repositoriesUtil.NamedGet(r.DB, selectQuery, WeightEntry{}, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil
}
