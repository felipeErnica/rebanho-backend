package weight

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
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
		WITH weight_gain AS (
			SELECT 
				DATE_TRUNC('month', w.entry_date) entry_date,
				(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) / 
				EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)) AS daily_gain
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		),
		cte AS (
			SELECT 
				entry_date,
				AVG(daily_gain) average_gain
			FROM weight_gain
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 60
		)
		SELECT * FROM cte ORDER BY entry_date
	`
	return util.GetList[AverageWeightGain](r.DB, query, userId)
}

func (r *WeightRepository) GetWeightHist(userId string) (*[]AverageWeight, error) {
	query := `
		WITH cte AS (
			SELECT 
				DATE_TRUNC('month', entry_date) entry_date,
				AVG(weight) average_weight
			FROM weight_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 60
		)
		SELECT * FROM cte ORDER BY entry_date
	`
	return util.GetList[AverageWeight](r.DB, query, userId)
}
func (r *WeightRepository) GetLastWeightGain(userId string) (*CardWeightGain, error) {
	query := `
		WITH weight_gain AS (
			SELECT 
				entry_date,
				(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) / 
				NULLIF(EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)), 0) AS daily_gain
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		),
		cte AS (
			SELECT 
				entry_date,
				AVG(daily_gain) average_gain
			FROM weight_gain
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY entry_date
	`
	result, err := util.GetList[AverageWeightGain](r.DB, query, userId)
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
		WITH cte AS (
			SELECT 
				entry_date,
				AVG(weight) average_weight
			FROM weight_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1 
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY entry_date
	`
	result, err := util.GetList[AverageWeight](r.DB, query, userId)
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
		WITH last_date AS (
			SELECT MAX(entry_date) entry_date FROM weight_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
		),
		stats AS (
			SELECT
				w.id,
				COALESCE(w.weight - LAG(w.weight) OVER win, 0) weight_variation,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) / 
					EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)), 
					0
				) AS weight_gain 
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		)
		SELECT 
			w.id,
			w.animal_id,
			CONCAT_WS(
				' - ',
				a.ring_number,
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) AS animal_info,
			a.birth_date,
			w.entry_date,
			w.weight,
			s.weight_variation,
			s.weight_gain
		FROM last_date l, weight_entries w 
			JOIN stats s ON s.id = w.id
			LEFT JOIN animals a ON a.id = w.animal_id
		WHERE w.entry_date = l.entry_date
			AND w.user_id = $1 
			AND w.deleted_at IS NULL
		ORDER BY COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0)
	`
	return util.GetList[WeightEntry](r.DB, query, userId)
}

func (r *WeightRepository) GetLastGroups(userId string) (*[]WeightGroup, error) {
	query := `
		WITH entries AS (
			SELECT 
				w.entry_date,
				w.weight,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) / 
					EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)), 
					0
				) AS weight_gain
			FROM  weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		),
		cte AS (
			SELECT
				w.entry_date,
				COUNT(w.weight) animals_number,
				AVG(w.weight) average_weight,
				AVG(w.weight_gain) average_gain
			FROM entries w
			GROUP BY 1
		)
		SELECT
			c.*,
			COALESCE( ((c.average_gain / LAG(c.average_gain) OVER win) - 1) * 100, 0) gain_variation,
			COALESCE( ((c.average_weight / LAG(c.average_weight) OVER win) - 1) * 100, 0) weight_variation 
		FROM cte c
		WINDOW win AS (ORDER BY c.entry_date)
		ORDER BY entry_date DESC
		LIMIT 10
	`
	return util.GetList[WeightGroup](r.DB, query, userId)
}

func (r *WeightRepository) GetBestFathers(userId string) (*[]AnimalRating, error) {
	query := `
		WITH gain_tbl AS (
			SELECT 
				w.animal_id,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) / 
					EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)), 
					0
				) AS weight_gain
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		),
		animal_tbl AS (
			SELECT 
				animal_id,
				AVG(weight_gain) animal_gain
			FROM gain_tbl
			GROUP BY 1
		),
		father_tbl AS (
			SELECT
				f.id father_id,
				COUNT(t.animal_id) children_number,
				AVG(t.animal_gain) avg_gain
			FROM animal_tbl t
				LEFT JOIN animals a ON a.id = t.animal_id
				LEFT JOIN animals f ON f.id = a.father_id
			GROUP BY 1
		),
		stats AS (
			SELECT AVG(weight_gain) gn_avg_gain
			FROM gain_tbl
		)
		SELECT 
			CONCAT_WS(' - ', f.ring_number, f.name) animal_name,
			t.avg_gain,
			((t.avg_gain / s.gn_avg_gain) - 1) * 100 gain_trend,
			t.children_number
		FROM stats s, father_tbl t 
			JOIN animals f ON f.id = t.father_id
		WHERE t.children_number >= 10
		ORDER BY t.avg_gain DESC
		LIMIT 10
	`
	return util.GetList[AnimalRating](r.DB, query, userId)
}

func (r *WeightRepository) GetBestMothers(userId string) (*[]AnimalRating, error) {
	query := `
		WITH gain_tbl AS (
			SELECT 
				w.animal_id,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) / 
					EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)), 
					0
				) AS weight_gain
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		),
		animal_tbl AS (
			SELECT 
				animal_id,
				AVG(weight_gain) animal_gain
			FROM gain_tbl
			GROUP BY 1
		),
		mother_tbl AS (
			SELECT
				m.id mother_id,
				COUNT(t.animal_id) children_number,
				AVG(t.animal_gain) avg_gain
			FROM animal_tbl t
				LEFT JOIN animals a ON a.id = t.animal_id
				LEFT JOIN animals m ON m.id = a.mother_id
			GROUP BY 1
		),
		stats AS (
			SELECT AVG(weight_gain) gn_avg_gain
			FROM gain_tbl
		)
		SELECT 
			CONCAT_WS(' - ', m.ring_number, m.name) animal_name,
			t.avg_gain,
			((t.avg_gain / s.gn_avg_gain) - 1) * 100 gain_trend,
			t.children_number
		FROM stats s, mother_tbl t 
			LEFT JOIN animals m ON m.id = t.mother_id
		WHERE t.children_number >= 5
		ORDER BY t.avg_gain DESC
		LIMIT 10
	`
	return util.GetList[AnimalRating](r.DB, query, userId)
}

func (r *WeightRepository) FindEntriesPage(
	sort string,
	order string,
	cursor string,
	filter *WeightFilter,
	userId string,
) (*entity.Page[WeightEntry], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"entry_date":   {Field: "w.entry_date", Order: "desc"},
		"animal_order": {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":  {Field: "coalesce(a.name, '')", Order: "asc"},
		"birth_date":   {Field: "coalesce(a.birth_date, '-infinity')", Order: "asc"},
		"id":           {Field: "w.id", Order: "asc"},
		"created_at":   {Field: "w.created_at", Order: "asc"},
	}

	query := `
		WITH stats AS (
			SELECT
				w.id,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) /
					NULLIF(EXTRACT(day FROM (w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date))), 0),
					0
				) weight_gain,
				COALESCE(w.weight - LAG(w.weight) OVER win, 0) AS weight_variation
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		)
		SELECT
			w.id,
			w.animal_id,
			COALESCE(a.name, '') AS animal_name,
			CONCAT_WS(
				' - ',
				a.ring_number,
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) AS animal_info,
			a.birth_date,
			COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0) AS animal_order,
			CONCAT_WS(' - ', f.ring_number, f.name) AS father_name,
			CONCAT_WS(' - ', m.ring_number, m.name) AS mother_name,
			s.weight_gain,
			s.weight_variation,
			w.entry_date,
			w.weight,
			w.created_at
		FROM weight_entries w 
			JOIN stats s ON s.id = w.id
			LEFT JOIN animals a ON a.id = w.animal_id
			LEFT JOIN animals m ON m.id = a.mother_id
			LEFT JOIN animals f ON f.id = a.father_id
	`

	whereExpression := " WHERE w.user_id = $1 AND w.deleted_at IS NULL"

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "w", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " AND " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression += " AND " + cursorExpression
	}

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}
	orderExpression := " ORDER BY " + sortExpression

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	windowExpression := " WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)"
	query += whereExpression + windowExpression + orderExpression
	return util.GetPage[WeightEntry](r.DB, query, sort, 200, args...)
}

func (r *WeightRepository) GetEntriesPageFoot(filter *WeightFilter, userId string) (*WeightFoot, error) {
	query := `
		SELECT
			w.animal_id,
			w.weight,
			COALESCE(
				(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) /
				NULLIF(EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)), 0),
				0
			) weight_gain
		FROM weight_entries w 
			LEFT JOIN animals a ON a.id = w.animal_id
	`
	whereExpression := " WHERE w.user_id = $1 AND w.deleted_at IS NULL"
	filterExpression, _, err := util.GetFilterExpressions(filter, "w", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " AND " + filterExpression
	}

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	windowExpression := " WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)"
	query += whereExpression + windowExpression
	query = fmt.Sprintf(`
		WITH cte AS (%s)
		SELECT
			COUNT(*) AS animals_num,
			AVG(weight) AS avg_weight,
			AVG(weight_gain) AS avg_gain
		FROM cte
	`, query)

	return util.GetOne[WeightFoot](r.DB, query, args...)
}

func (r *WeightRepository) FindGroups(userId string, order string) (*[]WeightGroup, error) {
	query := `
		WITH entries AS (
			SELECT 
				w.entry_date,
				w.weight,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) / 
					EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)),
					0 
				) weight_gain
			FROM  weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		),
		cte AS (
			SELECT
				w.entry_date,
				COUNT(w.weight) animals_number,
				AVG(w.weight) average_weight,
				AVG(w.weight_gain) average_gain
			FROM entries w
			GROUP BY 1
			ORDER BY 1
		)
		SELECT
			c.*,
			COALESCE(
				((c.average_gain / LAG(c.average_gain) OVER (ORDER BY c.entry_date)) - 1) * 100
			, 0) gain_variation,
			COALESCE(
				((c.average_weight / LAG(c.average_weight) OVER (ORDER BY c.entry_date)) - 1) * 100
			, 0) weight_variation
		FROM cte c
	`
	orderExpression := " ORDER BY c.entry_date " + order
	query += orderExpression
	return util.GetList[WeightGroup](r.DB, query, userId)
}

func (r *WeightRepository) FindEntriesByDate(
	entryDate time.Time,
	userId string,
	order string,
	sort string,
) (*[]WeightEntry, error) {

	sortMap := map[string]util.SortField{
		"animal_order": {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":  {Field: "a.name", Order: "asc"},
		"birth_date":   {Field: "coalesce(a.birth_date, '-infinity')", Order: "desc"},
	}

	query := `
		WITH stats AS (
			SELECT
				w.id,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) /
					NULLIF(EXTRACT(day FROM (w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date))), 0),
					0
				) weight_gain,
				COALESCE(w.weight - LAG(w.weight) OVER win, 0) AS weight_variation
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $2 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		)
		SELECT 
			w.id,
			w.animal_id,
			COALESCE(a.name, '') AS animal_name,
			CONCAT_WS(
				' - ', 
				a.ring_number, 
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) AS animal_info,
			a.birth_date,
			CONCAT_WS(' - ', f.ring_number, f.name) father_name,
			CONCAT_WS(' - ', m.ring_number, m.name) mother_name,
			w.entry_date,
			w.weight,
			s.weight_variation,
			s.weight_gain
		FROM weight_entries w 
			LEFT JOIN animals a ON a.id = w.animal_id
			LEFT JOIN animals f ON f.id = a.father_id
			LEFT JOIN animals m ON m.id = a.mother_id
			JOIN stats s ON s.id = w.id
		WHERE w.entry_date = $1
			AND w.user_id = $2 
			AND w.deleted_at IS NULL
		WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
	`
	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query += " ORDER BY " + sortExpression

	return util.GetList[WeightEntry](r.DB, query, entryDate, userId)
}

func (r *WeightRepository) GetEntriesByDateFoot(entryDate time.Time, userId string) (*WeightFoot, error) {
	query := `
		WITH cte AS (
			SELECT
				w.animal_id,
				w.weight,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) /
					NULLIF(EXTRACT(days FROM w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date)), 0),
					0
				) weight_gain
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.entry_date = $1 AND w.user_id = $2 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		) 
		SELECT 
			COUNT(animal_id) animals_num,
			AVG(weight) avg_weight,
			AVG(weight_gain) avg_gain
		FROM cte
	`
	return util.GetOne[WeightFoot](r.DB, query, entryDate, userId)
}

func (r *WeightRepository) Delete(id string, userId string) *log.APIError {
	query := `
		UPDATE weight_entries
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	err := util.Exec(r.DB, query, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *WeightRepository) Update(entry *WeightEntrySave) (*WeightEntry, *log.APIError) {

	validateErr := validateUpdate(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		UPDATE weight_entries
		SET entry_date = :entry_date,
			weight = :weight
		WHERE id = :id AND user_id = :user_id
	`
	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	selectQuery := `
		WITH stats AS (
			SELECT
				w.id,
				COALESCE(
					(w.weight - COALESCE(LAG(w.weight) OVER win, 38)) /
					NULLIF(EXTRACT(day FROM (w.entry_date - COALESCE(LAG(w.entry_date) OVER win, a.birth_date))), 0),
					0
				) weight_gain,
				COALESCE(w.weight - LAG(w.weight) OVER win, 0) AS weight_variation
			FROM weight_entries w 
				LEFT JOIN animals a ON a.id = w.animal_id
			WHERE w.user_id = $1 AND w.deleted_at IS NULL
			WINDOW win AS (PARTITION BY w.animal_id ORDER BY w.entry_date)
		)
		SELECT
			w.id,
			w.animal_id,
			CONCAT_WS(
				' - ',
				a.ring_number,
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) AS animal_info,
			a.birth_date,
			CONCAT_WS(' - ', f.ring_number, f.name) AS father_name,
			CONCAT_WS(' - ', m.ring_number, m.name) AS mother_name,
			s.weight_gain,
			s.weight_variation,
			w.entry_date,
			w.weight
		FROM weight_entries w 
			JOIN stats s ON s.id = w.id
			LEFT JOIN animals a ON a.id = w.animal_id
			LEFT JOIN animals m ON m.id = a.mother_id
			LEFT JOIN animals f ON f.id = a.father_id
		WHERE id = :id AND user_id = :user_id
	`

	response, err := util.NamedGet(r.DB, selectQuery, WeightEntry{}, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return response, nil
}
