package slaughter

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type SlaughterRepository struct {
	DB *sqlx.DB
}

func NewSlaughterRepository(db *sqlx.DB) *SlaughterRepository {
	return &SlaughterRepository{db}
}

func (r *SlaughterRepository) GetLastPerformance(userId string) (*PerformanceRateCard, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0)) * 100) performance_rate
			FROM slaughter_entries s
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY entry_date
	`

	result, err := util.GetList[PerformanceRateHist](r.DB, query, userId)
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
		WITH cte AS (
			SELECT 
				entry_date,
				AVG(weight) avg_weight
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY entry_date
	`

	result, err := util.GetList[AverageWeightHist](r.DB, query, userId)
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
		WITH cte AS (
			SELECT 
				entry_date,
				AVG(dead_weight) avg_weight
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY entry_date
	`

	result, err := util.GetList[AverageWeightHist](r.DB, query, userId)
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
		WITH cte AS (
			SELECT 
				DATE_TRUNC('month', s.entry_date) entry_date,
				AVG(weight) avg_weight,
				AVG(dead_weight) dead_weight
			FROM slaughter_entries s
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 60
		)
		SELECT * FROM cte ORDER BY 1
	`
	return util.GetList[WeightHist](r.DB, query, userId)
}

func (r *SlaughterRepository) GetRateHist(userId string) (*[]RateHist, error) {
	query := `
		WITH cte AS (
			SELECT 
				DATE_TRUNC('month', s.entry_date) entry_date,
				AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0)) * 100) avg_rate
			FROM slaughter_entries s
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 60
		)
		SELECT * FROM cte ORDER BY 1
	`
	return util.GetList[RateHist](r.DB, query, userId)
}
func (r *SlaughterRepository) GetBestFathers(userId string) (*[]TableRatings, error) {
	query := `
		WITH cte AS (
			SELECT 
				a.father_id,
				COUNT(*) animals_number,
				AVG(weight) avg_weight,
				AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0))*100) avg_rate
			FROM slaughter_entries s JOIN animals a ON a.id = s.animal_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
			GROUP BY 1
			HAVING COUNT(*) >= 10
		),
		gn_stats AS (
			SELECT 
				AVG(avg_weight) gn_avg_weight,
				AVG(avg_rate) gn_avg_rate
			FROM cte
		)
		SELECT 
			CONCAT_WS(' - ', a.ring_number, a.name) name,
			c.avg_weight,
			((c.avg_weight / NULLIF(s.gn_avg_weight, 0)) - 1) * 100 weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(s.gn_avg_rate, 0)) - 1) * 100 rate_comparison,
			c.animals_number
		FROM cte c 
			CROSS JOIN gn_stats s
			JOIN animals a ON a.id = c.father_id
		ORDER BY avg_rate DESC
		LIMIT 10
	`
	return util.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetBestMothers(userId string) (*[]TableRatings, error) {
	query := `
		WITH cte AS (
			SELECT 
				a.mother_id,
				COUNT(*) animals_number,
				AVG(weight) avg_weight,
				AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0))*100) avg_rate
			FROM slaughter_entries s JOIN animals a ON a.id = s.animal_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
			GROUP BY 1
			HAVING COUNT(*) >= 10
		),
		gn_stats AS (
			SELECT 
				AVG(avg_weight) gn_avg_weight,
				AVG(avg_rate) gn_avg_rate
			FROM cte
		)
		SELECT 
			CONCAT_WS(' - ', a.ring_number, a.name) name,
			c.avg_weight,
			((c.avg_weight / NULLIF(s.gn_avg_weight, 0)) - 1) * 100 weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(s.gn_avg_rate, 0)) - 1) * 100 rate_comparison,
			c.animals_number
		FROM cte c 
			CROSS JOIN gn_stats s
			JOIN animals a ON a.id = c.mother_id
		ORDER BY c.avg_rate DESC
		LIMIT 10
	`
	return util.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetBestSlaughterhouses(userId string) (*[]TableRatings, error) {
	query := `
		WITH cte AS (
			SELECT 
				butcher_id,
				COUNT(*) animals_number,
				AVG(weight) avg_weight,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0))*100) avg_rate
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
		),
		gn_stats AS (
			SELECT 
				AVG(avg_weight) gn_avg_weight,
				AVG(avg_rate) gn_avg_rate
			FROM cte
		)
		SELECT 
			s.name,
			c.avg_weight,
			((c.avg_weight / NULLIF(g.gn_avg_weight, 0)) - 1) * 100 weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(g.gn_avg_rate, 0)) - 1) * 100 rate_comparison,
			c.animals_number
		FROM gn_stats g, cte c JOIN butchers s ON s.id = c.butcher_id
		WHERE c.avg_rate >= 20
		ORDER BY c.avg_rate DESC
		LIMIT 10
	`
	return util.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastGroups(userId string) (*[]SlaughterGroup, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date,
				butcher_id,
				COUNT(animal_id) animals_number,
				AVG(dead_weight) avg_weight,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0))*100) avg_rate
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1, 2
		)
		SELECT 
			c.entry_date,
			s.name butcher,
			c.avg_weight,
			COALESCE(((c.avg_weight / LAG(NULLIF(c.avg_weight, 0)) OVER (ORDER BY c.entry_date)) - 1) * 100, 0) weight_variation,
			c.avg_rate,
			COALESCE(((c.avg_rate / LAG(NULLIF(c.avg_rate, 0)) OVER (ORDER BY c.entry_date)) - 1) * 100, 0) rate_variation,
			c.animals_number
		FROM cte c 
			JOIN butchers s ON s.id = c.butcher_id
		ORDER BY c.entry_date DESC
		LIMIT 10
	`
	return util.GetList[SlaughterGroup](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastEntries(userId string) (*[]SlaughterEntry, error) {
	query := `
		WITH last_date AS (
			SELECT MAX(entry_date) max_date 
			FROM slaughter_entries
			WHERE user_id = $1 AND deleted_at IS NULL
		)
		SELECT 
			CONCAT_WS(
				' - ', 
				a.ring_number, 
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) AS animal_info,
			h.name butcher,
			s.entry_date,
			s.weight,
			s.discount_rate * 100 AS discount_rate,
			COALESCE(s.weight * (1 - s.discount_rate), 0) discount_weight,
			s.dead_weight,
			COALESCE(s.dead_weight / (s.weight * (1 - s.discount_rate)), 0) * 100 performance_rate
		FROM slaughter_entries s
			CROSS JOIN last_date l
			LEFT JOIN animals a ON a.id = s.animal_id
			JOIN butchers h ON h.id = s.butcher_id
		WHERE s.entry_date = l.max_date
			AND s.user_id = $1 
			AND s.deleted_at IS NULL
		ORDER BY COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0) 
	`
	return util.GetList[SlaughterEntry](r.DB, query, userId)
}

func (r *SlaughterRepository) FindEntriesPage(
	sort string,
	order string,
	cursor string,
	filter *SlaughterEntryFilter,
	userId string,
) (*entity.Page[SlaughterEntry], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
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
		SELECT 
			s.id,
			s.animal_id,
			s.butcher_id,
			COALESCE(REGEXP_REPLACE(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			a.name AS animal_name, 
			CONCAT_WS(
				' - ', 
				a.ring_number, 
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) AS animal_info,
			a.birth_date,
			CONCAT_WS(' - ', f.ring_number, f.name) father_name,
			CONCAT_WS(' - ', m.ring_number, m.name) mother_name,
			h.name butcher,
			s.entry_date,
			s.discount_rate * 100 AS discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) discount_weight,
			s.dead_weight,
			COALESCE(s.dead_weight / NULLIF(s.weight*(1 - s.discount_rate), 0) * 100, 0) performance_rate,
			s.created_at
		FROM slaughter_entries s
			JOIN butchers h ON h.id = s.butcher_id
			LEFT JOIN animals a ON a.id = s.animal_id
			LEFT JOIN animals f ON f.id = a.father_id
			LEFT JOIN animals m ON m.id = a.mother_id
	`

	whereExpression := "WHERE s.user_id = $1 AND s.deleted_at IS NULL"

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "s", 2)
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
	query += whereExpression + orderExpression

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetPage[SlaughterEntry](r.DB, query, sort, 200, args...)
}

func (r *SlaughterRepository) GetEntriesPageFoot(filter *SlaughterEntryFilter, userId string) (*SlaughterFoot, error) {

	query := `
		SELECT 
			COUNT(s.*) AS animals_number,
			AVG(s.weight) AS avg_weight,
			AVG(s.dead_weight) AS avg_dead_weight,
			AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0)) * 100) AS avg_rate
		FROM slaughter_entries s
	`

	whereExpression := " WHERE s.user_id = $1 AND s.deleted_at IS NULL"

	filterExpression, _, err := util.GetFilterExpressions(filter, "s", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " AND " + filterExpression
	}

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query += whereExpression

	return util.GetOne[SlaughterFoot](r.DB, query, args...)
}

func (r *SlaughterRepository) FindGroups(order string, userId string) (*[]SlaughterGroup, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date,
				butcher_id,
				COUNT(*) AS animals_number,
				AVG(weight) AS avg_weight,
				AVG(dead_weight) AS avg_dead_weight,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0)) * 100) AS avg_rate
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1, 2
		)
		SELECT 
			c.entry_date,
			s.name butcher,
			c.avg_weight,
			c.avg_dead_weight,
			COALESCE((c.avg_weight / LAG(NULLIF(c.avg_weight, 0)) OVER win) - 1, 0) * 100 AS weight_variation, 
			COALESCE((c.avg_dead_weight / LAG(NULLIF(c.avg_dead_weight, 0)) OVER win) - 1, 0) * 100 AS dead_weight_variation,
			c.avg_rate,
			COALESCE((c.avg_rate / LAG(NULLIF(c.avg_rate, 0)) OVER win) - 1, 0) * 100 AS rate_variation,
			c.animals_number
		FROM cte c 
			JOIN butchers s ON s.id = c.butcher_id
		WINDOW win AS (ORDER BY c.entry_date)
		ORDER BY c.entry_date
	`
	query += order
	return util.GetList[SlaughterGroup](r.DB, query, userId)
}

func (r *SlaughterRepository) FindEntriesByDate(
	sort string,
	order string,
	entryDate time.Time,
	userId string,
) (*[]SlaughterEntry, error) {

	sortMap := map[string]util.SortField{
		"animal_order":     {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "coalesce(a.name, '')", Order: "asc"},
		"birth_date":       {Field: "coalesce(a.birth_date, '-infinity')", Order: "desc"},
		"weight":           {Field: "s.weight", Order: "asc"},
		"dead_weight":      {Field: "s.dead_weight", Order: "asc"},
		"performance_rate": {Field: "coalesce(s.dead_weight / nullif(s.weight * (1 - s.discount_rate), 0), 0) * 100", Order: "asc"},
	}

	query := `
		SELECT 
			s.id,
			s.animal_id,
			s.butcher_id,
			a.name AS animal_name,
			CONCAT_WS(
				' - ', 
				a.ring_number, 
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) animal_info,
			CONCAT_WS(' - ', m.ring_number, m.name) AS mother_name,
			CONCAT_WS(' - ', f.ring_number, f.name) AS father_name,
			s.weight,
			s.discount_rate * 100 AS discount_rate,
			s.weight * (1 - s.discount_rate) AS discount_weight,
			s.dead_weight,
			COALESCE(s.dead_weight / NULLIF(s.weight * (1 - s.discount_rate), 0) * 100, 0) AS performance_rate
		FROM slaughter_entries s 
			LEFT JOIN animals a ON a.id = s.animal_id
			LEFT JOIN animals m ON m.id = a.mother_id
			LEFT JOIN animals f ON f.id = a.father_id
		WHERE s.entry_date = $1
			AND s.user_id = $2 
			AND s.deleted_at IS NULL
	`

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	orderExperssion := " ORDER BY " + sortExpression
	query += orderExperssion

	return util.GetList[SlaughterEntry](r.DB, query, entryDate, userId)
}

func (r *SlaughterRepository) GetEntriesByDateFoot(entryDate time.Time, userId string) (*SlaughterFoot, error) {
	query := `
		SELECT 
			COUNT(s.animal_id) animals_number,
			AVG(s.weight) avg_weight,
			AVG(s.dead_weight) avg_dead_weight,
			AVG(COALESCE(s.dead_weight / NULLIF(s.weight * (1 - s.discount_rate), 0), 0) * 100) avg_rate
		FROM slaughter_entries s 
		WHERE s.entry_date = $1
			AND s.user_id = $2 
			AND s.deleted_at IS NULL
	`
	return util.GetOne[SlaughterFoot](r.DB, query, entryDate, userId)
}

func (r *SlaughterRepository) Delete(id string, userId string) *log.APIError {

	tx, err := r.DB.Beginx()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	entryQuery := `
		SELECT
			id,
			animal_id,
			butcher_id,
			entry_date,
			weight,
			dead_weight,
			discount_rate,
			user_id
		FROM slaughter_entries
		WHERE id = $1 AND user_id = $2
	`
	entry, err := util.GetOneTx(tx, entryQuery, &SlaughterEntrySave{}, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	deleteQuery := `
		UPDATE slaughter_entries 
		SET deleted_at = NOW()
		WHERE id = :id AND user_id = :user_id
	`

	err = util.NamedExecTx(tx, deleteQuery, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	animalQuery := `
		UPDATE animals 
		SET death_date = NULL,
		WHERE id = :animal_id
			AND user_id = :user_id
			AND death_date = :entry_date
	`
	err = util.NamedExecTx(tx, animalQuery, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *SlaughterRepository) Update(entry *SlaughterEntrySave) (*SlaughterEntry, *log.APIError) {

	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	oldQuery := `
		SELECT
			id,
			animal_id,
			butcher_id,
			entry_date,
			weight,
			dead_weight,
			discount_rate,
			user_id
		FROM slaughter_entries
		WHERE id = $1 AND user_id = $2
	`
	oldEntry, err := util.GetOneTx(tx, oldQuery, &SlaughterEntrySave{}, entry.Id, entry.UserId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	validateErr := validateUpdate(tx, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		UPDATE slaughter_entries 
		SET entry_date = :entry_date,
			discount_rate = :discount_rate,
			weight = :weight,
			dead_weight = :dead_weight,
			butcher_id = :butcher_id
		WHERE id = :id AND user_id = :user_id
	`

	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	animalQuery := `
		UPDATE animals 
		SET death_date = $1,
		WHERE id = $2
			AND user_id = $3
			AND death_date = $4
	`
	err = util.ExecTx(tx, animalQuery, entry.EntryDate, entry.AnimalId, entry.UserId, oldEntry.EntryDate)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	selectQuery := `
		SELECT 
			s.id,
			s.animal_id,
			s.butcher_id,
			CONCAT_WS(
				' - ', 
				a.ring_number, 
				COALESCE(a.name, a.sex),
				TO_CHAR(a.birth_date, 'DD/MM/YYYY')
			) AS animal_info,
			a.birth_date,
			CONCAT_WS(' - ', f.ring_number, f.name) father_name,
			CONCAT_WS(' - ', m.ring_number, m.name) mother_name,
			h.name butcher,
			s.entry_date,
			s.discount_rate * 100 AS discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) AS discount_weight,
			s.dead_weight,
			COALESCE(s.dead_weight / NULLIF(s.weight*(1 - s.discount_rate), 0) * 100, 0) performance_rate
		FROM slaughter_entries s
			JOIN butchers h ON h.id = s.butcher_id
			LEFT JOIN animals a ON a.id = s.animal_id
			LEFT JOIN animals f ON f.id = a.father_id
			LEFT JOIN animals m ON m.id = a.mother_id
		WHERE id = :id AND user_id = :user_id
	`

	result, err := util.NamedGetTx(tx, selectQuery, SlaughterEntry{}, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return result, nil
}

func (r *SlaughterRepository) Add(entry *SlaughterEntrySave) *log.APIError {

	validateErr := validateAdd(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		INSERT INTO slaughter_entries (
			animal_id, 
			butcher_id, 
			entry_date, 
			discount_rate, 
			weight, 
			dead_weight,
			user_id
		)
		VALUES (
			:animal_id,
			:butcher_id,
			:entry_date,
			:discount_rate,
			:weight,
			:dead_weight,
			:user_id
		)
	`

	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	animalQuery := `
		UPDATE animals
		SET death_date = $1
		WHERE id = $1 AND user_id = $2
	`
	err = util.ExecTx(tx, animalQuery, entry.AnimalId, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *SlaughterRepository) Replace(entry *SlaughterEntrySave) *log.APIError {

	validateErr := validateAdd(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		UPDATE slaughter_entries 
		SET discount_rate = :discount_rate, 
			weight = :weight, 
			dead_weight = :dead_weight
		WHERE entry_date = :entry_date
			AND animal_id = :animal_id
			AND user_id = :user_id
			AND deleted_at IS NULL
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
