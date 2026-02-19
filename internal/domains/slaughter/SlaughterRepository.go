package slaughter

import (
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type SlaughterRepository struct {
	DB *sqlx.DB
}

type SaveValidation struct {
	Exists   bool `db:"slaughter_exists"`
	HasDeath bool `db:"has_death"`
}

type SaveValidationBatch struct {
	Exists   int `db:"slaughter_exists"`
	HasDeath int `db:"has_death"`
}

func NewRepository(db *sqlx.DB) *SlaughterRepository {
	return &SlaughterRepository{db}
}

func (r *SlaughterRepository) CheckSave(entry *SlaughterSave) (*SaveValidation, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM slaughter_entries
			WHERE entry_date = :entry_date
				AND animal_id = :animal_id
				AND id IS DISTINCT FROM :id
				AND user_id = :user_id
				AND deleted_at IS NULL
		) AS slaughter_exists,
		EXISTS (
			SELECT 1
			FROM animals
			WHERE id = :animal_id
				AND death_date IS NOT NULL
				AND death_date IS DISTINCT FROM :entry_date
				AND user_id = :user_id
				AND deleted_at IS NULL
		) AS has_death
	`
	validate, err := util.NamedGet(r.DB, query, SaveValidation{}, entry)
	if err != nil {
		return nil, err
	}

	return validate, nil
}

func (r *SlaughterRepository) GetLastPerformance(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date AS date,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0)) ) AS value
			FROM slaughter_entries s
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY date
	`
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastAverageWeight(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date AS date,
				AVG(weight) AS value
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY date
	`
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastDeadWeight(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date AS date,
				AVG(dead_weight) AS value
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY date
	`
	return util.GetList[util.GraphData](r.DB, query, userId)
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

func (r *SlaughterRepository) GetRateHist(userId string) (*[]util.GraphData, error) {
	query := `
		WITH cte AS (
			SELECT 
				DATE_TRUNC('month', s.entry_date) AS date,
				AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0)) ) AS value
			FROM slaughter_entries s
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 60
		)
		SELECT * FROM cte ORDER BY 1
	`
	return util.GetList[util.GraphData](r.DB, query, userId)
}
func (r *SlaughterRepository) GetBestFathers(userId string) (*[]TableRatings, error) {
	query := `
		WITH cte AS (
			SELECT 
				a.father_id,
				COUNT(*) animals_number,
				AVG(weight) avg_weight,
				AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0))) avg_rate
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
			CONCAT_WS(' - ', a.tag, a.name) name,
			c.avg_weight,
			((c.avg_weight / NULLIF(s.gn_avg_weight, 0)) - 1)  weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(s.gn_avg_rate, 0)) - 1)  rate_comparison,
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
				AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0))) avg_rate
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
			CONCAT_WS(' - ', a.tag, a.name) name,
			c.avg_weight,
			((c.avg_weight / NULLIF(s.gn_avg_weight, 0)) - 1)  weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(s.gn_avg_rate, 0)) - 1)  rate_comparison,
			c.animals_number
		FROM cte c 
			CROSS JOIN gn_stats s
			JOIN animals a ON a.id = c.mother_id
		ORDER BY c.avg_rate DESC
		LIMIT 10
	`
	return util.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetBestButchers(userId string) (*[]TableRatings, error) {
	query := `
		WITH cte AS (
			SELECT 
				butcher_id,
				COUNT(*) animals_number,
				AVG(weight) avg_weight,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0))) avg_rate
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
			((c.avg_weight / NULLIF(g.gn_avg_weight, 0)) - 1)  weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(g.gn_avg_rate, 0)) - 1)  rate_comparison,
			c.animals_number
		FROM gn_stats g, cte c JOIN butchers s ON s.id = c.butcher_id
		WHERE c.avg_rate >= 20
		ORDER BY c.avg_rate DESC
		LIMIT 10
	`
	return util.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastGroups(userId string) (*[]SlaughterGroupDB, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date,
				butcher_id,
				COUNT(animal_id) animals_number,
				AVG(dead_weight) avg_weight,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0))) avg_rate
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1, 2
		)
		SELECT 
			c.entry_date,
			c.avg_weight,
			((c.avg_weight / LAG(NULLIF(c.avg_weight, 0)) OVER (ORDER BY c.entry_date)) - 1)  AS weight_variation,
			c.avg_rate,
			((c.avg_rate / LAG(NULLIF(c.avg_rate, 0)) OVER (ORDER BY c.entry_date)) - 1)  AS rate_variation,
			c.animals_number,

			b.id AS butcher_id,
			b.name AS butcher_name,
			b.discount AS butcher_discount
		FROM cte c 
		JOIN butchers b ON b.id = c.butcher_id
		ORDER BY c.entry_date DESC
		LIMIT 10
	`
	return util.GetList[SlaughterGroupDB](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastEntries(userId string) (*[]SlaughterDB, error) {
	query := `
		WITH last_date AS (
			SELECT MAX(entry_date) max_date 
			FROM slaughter_entries
			WHERE user_id = $1 AND deleted_at IS NULL
		)
		SELECT 
			s.id,
			s.animal_id,
			a.tag AS animal_tag, 
			a.name AS animal_name, 
			a.sex AS animal_sex,
			a.birth_date AS animal_birth,
			COALESCE(NULLIF(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g'), '')::int, 0) animal_order,

			a.father_id,
			f.tag AS father_tag,
			f.name AS father_name,

			a.mother_id,
			m.tag AS mother_tag,
			m.name AS mother_name,

			s.butcher_id,
			h.name butcher_name,

			s.entry_date,
			s.discount_rate AS discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) discount_weight,
			s.dead_weight,
			s.dead_weight / NULLIF(s.weight*(1 - s.discount_rate), 0) AS performance_rate
		FROM slaughter_entries s
		CROSS JOIN last_date l
		LEFT JOIN animals a ON a.id = s.animal_id
		LEFT JOIN animals m ON m.id = a.mother_id
		LEFT JOIN animals f ON f.id = a.father_id
		JOIN butchers h ON h.id = s.butcher_id
		WHERE s.entry_date = l.max_date
			AND s.user_id = $1 
			AND s.deleted_at IS NULL
		ORDER BY COALESCE(NULLIF(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g'), '')::int, 0) 
	`
	return util.GetList[SlaughterDB](r.DB, query, userId)
}

func (r *SlaughterRepository) FindPage(
	sort string,
	order string,
	cursor string,
	limit int,
	filter *SlaughterFilter,
	userId string,
) (*[]SlaughterDB, error) {

	sortMap := map[string]util.SortField{
		"entry_date":       {Field: "s.entry_date", Order: "desc"},
		"animal_order":     {Field: "COALESCE(NULLIF(regexp_replace(a.tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "COALESCE(a.name, '')", Order: "asc"},
		"birth_date":       {Field: "COALESCE(a.birth_date, '-infinity')", Order: "desc"},
		"weight":           {Field: "s.weight", Order: "asc"},
		"dead_weight":      {Field: "s.dead_weight", Order: "asc"},
		"performance_rate": {Field: "COALESCE(s.dead_weight / NULLIF(s.weight*(1 - s.discount_rate), 0) , 0)", Order: "asc"},
		"id":               {Field: "s.id", Order: "asc"},
		"created_at":       {Field: "s.created_at", Order: "asc"},
	}

	query := `
		SELECT 
			s.id,
			s.animal_id,
			a.tag AS animal_tag, 
			a.name AS animal_name, 
			a.sex AS animal_sex,
			a.birth_date AS animal_birth,
			COALESCE(NULLIF(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g'), '')::int, 0) animal_order,

			a.father_id,
			f.tag AS father_tag,
			f.name AS father_name,

			a.mother_id,
			m.tag AS mother_tag,
			m.name AS mother_name,

			s.butcher_id,
			h.name butcher_name,

			s.entry_date,
			s.discount_rate AS discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) discount_weight,
			s.dead_weight,
			s.dead_weight / NULLIF(s.weight*(1 - s.discount_rate), 0) AS performance_rate,
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
	query += whereExpression + orderExpression + fmt.Sprintf(" LIMIT %d", limit)

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetList[SlaughterDB](r.DB, query, args...)
}

func (r *SlaughterRepository) GetPageFoot(filter *SlaughterFilter, userId string) (*SlaughterFoot, error) {

	query := `
		SELECT 
			COUNT(s.*) AS animals_number,
			AVG(s.weight) AS avg_weight,
			AVG(s.dead_weight) AS avg_dead_weight,
			AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0)) ) AS avg_rate
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

func (r *SlaughterRepository) FindButcherPage(
	sort string,
	order string,
	cursor string,
	filter *SlaughterFilter,
	butcherId string,
	limit int,
	userId string,
) (*[]SlaughterDB, error) {

	sortMap := map[string]util.SortField{
		"entry_date":       {Field: "s.entry_date", Order: "desc"},
		"animal_order":     {Field: "COALESCE(NULLIF(regexp_replace(a.tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "COALESCE(a.name, '')", Order: "asc"},
		"birth_date":       {Field: "COALESCE(a.birth_date, '-infinity')", Order: "desc"},
		"weight":           {Field: "s.weight", Order: "asc"},
		"dead_weight":      {Field: "s.dead_weight", Order: "asc"},
		"performance_rate": {Field: "COALESCE(s.dead_weight / NULLIF(s.weight*(1 - s.discount_rate), 0) , 0)", Order: "asc"},
		"id":               {Field: "s.id", Order: "asc"},
		"created_at":       {Field: "s.created_at", Order: "asc"},
	}

	query := `
		SELECT 
			s.id,
			s.animal_id,
			a.tag AS animal_tag, 
			a.name AS animal_name, 
			a.sex AS animal_sex,
			a.birth_date AS animal_birth,
			COALESCE(NULLIF(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g'), '')::int, 0) animal_order,

			a.father_id,
			f.tag AS father_tag,
			f.name AS father_name,

			a.mother_id,
			m.tag AS mother_tag,
			m.name AS mother_name,

			s.butcher_id,
			h.name butcher_name,

			s.entry_date,
			s.discount_rate AS discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) discount_weight,
			s.dead_weight,
			s.dead_weight / NULLIF(s.weight * (1 - s.discount_rate), 0) AS performance_rate,
			s.created_at
		FROM slaughter_entries s
		JOIN butchers h ON h.id = s.butcher_id
		LEFT JOIN animals a ON a.id = s.animal_id
		LEFT JOIN animals f ON f.id = a.father_id
		LEFT JOIN animals m ON m.id = a.mother_id
	`

	whereExpression := `
		WHERE s.user_id = $1 
			AND s.butcher_id = $2
			AND s.deleted_at IS NULL
		`

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "s", 3)
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
	query += whereExpression + orderExpression + fmt.Sprintf(" LIMIT %d", limit)

	args := []any{userId, butcherId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetList[SlaughterDB](r.DB, query, args...)
}

func (r *SlaughterRepository) GetButcherPageFoot(
	filter *SlaughterFilter,
	butcherId string,
	userId string,
) (*SlaughterFoot, error) {

	query := `
		SELECT 
			COUNT(s.*) AS animals_number,
			AVG(s.weight) AS avg_weight,
			AVG(s.dead_weight) AS avg_dead_weight,
			AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0)) ) AS avg_rate
		FROM slaughter_entries s
	`

	whereExpression := " WHERE s.user_id = $1 AND s.butcher_id = $2 AND s.deleted_at IS NULL"

	filterExpression, _, err := util.GetFilterExpressions(filter, "s", 3)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " AND " + filterExpression
	}

	args := []any{userId, butcherId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query += whereExpression

	return util.GetOne[SlaughterFoot](r.DB, query, args...)
}

func (r *SlaughterRepository) FindGroups(order string, userId string) (*[]SlaughterGroupDB, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date,
				butcher_id,
				COUNT(*) AS animals_number,
				AVG(weight) AS avg_weight,
				AVG(dead_weight) AS avg_dead_weight,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0)) ) AS avg_rate
			FROM slaughter_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1, 2
		)
		SELECT 
			c.entry_date,
			c.avg_weight,
			c.avg_dead_weight,
			COALESCE((c.avg_weight / LAG(NULLIF(c.avg_weight, 0)) OVER win) - 1, 0)  AS weight_variation, 
			COALESCE((c.avg_dead_weight / LAG(NULLIF(c.avg_dead_weight, 0)) OVER win) - 1, 0)  AS dead_weight_variation,
			c.avg_rate,
			COALESCE((c.avg_rate / LAG(NULLIF(c.avg_rate, 0)) OVER win) - 1, 0)  AS rate_variation,
			c.animals_number,

			c.butcher_id,
			b.name butcher_name,
			b.discount butcher_discount
		FROM cte c 
			JOIN butchers b ON b.id = c.butcher_id
		WINDOW win AS (ORDER BY c.entry_date)
		ORDER BY c.entry_date
	`
	query += order
	return util.GetList[SlaughterGroupDB](r.DB, query, userId)
}

func (r *SlaughterRepository) FindEntries(
	sort string,
	order string,
	filter *SlaughterFilter,
	userId string,
) (*[]SlaughterDB, error) {

	sortMap := map[string]util.SortField{
		"animal_order":     {Field: "COALESCE(NULLIF(regexp_replace(cte.animal_tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "COALESCE(cte.animal_name, '')", Order: "asc"},
		"birth_date":       {Field: "COALESCE(cte.animal_birth, '-infinity')", Order: "desc"},
		"weight":           {Field: "cte.weight", Order: "asc"},
		"dead_weight":      {Field: "cte.dead_weight", Order: "asc"},
		"performance_rate": {Field: "COALESCE(cte.performance_rate)", Order: "asc"},
	}

	query := `
		WITH cte AS (
			SELECT 
				s.id,
				s.animal_id,
				a.tag AS animal_tag, 
				a.name AS animal_name, 
				a.sex AS animal_sex,
				a.birth_date AS animal_birth,
				COALESCE(NULLIF(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g'), '')::int, 0) AS animal_order,

				a.father_id,
				f.tag AS father_tag,
				f.name AS father_name,

				a.mother_id,
				m.tag AS mother_tag,
				m.name AS mother_name,

				s.butcher_id,
				h.name butcher_name,

				s.entry_date,
				s.discount_rate AS discount_rate,
				s.weight,
				s.weight * (1 - s.discount_rate) AS discount_weight,
				s.dead_weight,
				(s.dead_weight / NULLIF(s.weight * (1 - s.discount_rate), 0)) AS performance_rate
			FROM slaughter_entries s 
			JOIN butchers h ON h.id = s.butcher_id
			LEFT JOIN animals a ON a.id = s.animal_id
			LEFT JOIN animals m ON m.id = a.mother_id
			LEFT JOIN animals f ON f.id = a.father_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
		)
		SELECT * FROM cte
	`

	filterExp, _, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExp)
	orderExperssion := " ORDER BY " + sortExpression
	query += whereExpression + orderExperssion
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetList[SlaughterDB](r.DB, query, args...)
}

func (r *SlaughterRepository) GetEntriesFoot(filter *SlaughterFilter, userId string) (*SlaughterFoot, error) {
	query := `
		WITH cte AS (
			SELECT 
				s.entry_date,
				s.animal_id,
				s.weight,
				s.dead_weight,
				(s.dead_weight / NULLIF(s.weight * (1 - s.discount_rate), 0))  AS rate
			FROM slaughter_entries s
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
		)

		SELECT 
			COUNT(s.animal_id) animals_number,
			AVG(s.weight) avg_weight,
			AVG(s.dead_weight) avg_dead_weight,
			AVG(s.rate) avg_rate
		FROM cte s 
	`
	filterExp, _, err := util.GetFilterExpressions(filter, "s", 2)
	if err != nil {
		return nil, err
	}
	whereExp := util.GetWhereExpression(filterExp)
	query += whereExp

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	return util.GetOne[SlaughterFoot](r.DB, query, args...)
}

func (r *SlaughterRepository) Delete(id string, userId string) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	entryQuery := `
		SELECT
			id,
			animal_id,
			user_id
		FROM slaughter_entries
		WHERE id = $1 AND user_id = $2
	`
	entry, err := util.GetOneTx(tx, entryQuery, &SlaughterSave{}, id, userId)
	if err != nil {
		return err
	}

	deleteQuery := `
		UPDATE slaughter_entries 
		SET deleted_at = NOW()
		WHERE id = :id AND user_id = :user_id
	`
	err = util.NamedExecTx(tx, deleteQuery, entry)
	if err != nil {
		return err
	}

	animalQuery := `
		UPDATE animals 
		SET death_date = NULL
		WHERE id = :animal_id AND user_id = :user_id
	`
	err = util.NamedExecTx(tx, animalQuery, entry)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *SlaughterRepository) DeleteBatch(ids []string, userId string) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	userExpress := "user_id = $1 AND deleted_at IS NULL"
	idExpress, _ := util.GetInExpression(ids, "id", 2)
	whereExpress := util.GetWhereExpression(userExpress, idExpress)

	batchQuery := `
		SELECT
			id,
			animal_id,
			user_id
		FROM slaughter_entries
	`
	batchQuery += whereExpress
	args := []any{userId}
	args = util.GetSliceArgs(args, ids)

	query := fmt.Sprintf(`
		UPDATE slaughter_entries 
		SET deleted_at = NOW()
		FROM (%s) batch_entries be
		WHERE id = be.id AND user_id = be.user_id
	`, batchQuery)

	err = util.ExecTx(tx, query, args...)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		UPDATE animals 
		SET death_date = NULL
		FROM (%s) AS batch_entries be
		WHERE id = be.animal_id AND user_id = be.user_id
	`, batchQuery)

	err = util.ExecTx(tx, query, args...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func (r *SlaughterRepository) Update(entry *SlaughterSave) (*SlaughterDB, error) {

	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

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
		return nil, err
	}

	query = `
		UPDATE animals 
		SET death_date = :entry_date
		WHERE id = :animal_id AND user_id = :user_id
	`
	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, err
	}

	query = `
		SELECT 
			s.id,
			s.animal_id,
			a.tag AS animal_tag, 
			a.name AS animal_name, 
			a.sex AS animal_sex,
			a.birth_date AS animal_birth,

			a.father_id,
			f.tag AS father_tag,
			f.name AS father_name,

			a.mother_id,
			m.tag AS mother_tag,
			m.name AS mother_name,

			s.butcher_id,
			b.name AS butcher_name,

			s.entry_date,
			s.discount_rate AS discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) AS discount_weight,
			s.dead_weight,
			s.dead_weight / NULLIF(s.weight * (1 - s.discount_rate), 0) AS performance_rate
		FROM slaughter_entries s
		JOIN butchers b ON b.id = s.butcher_id
		LEFT JOIN animals a ON a.id = s.animal_id
		LEFT JOIN animals f ON f.id = a.father_id
		LEFT JOIN animals m ON m.id = a.mother_id
		WHERE s.id = :id AND s.user_id = :user_id
	`

	result, err := util.NamedGetTx(tx, query, SlaughterDB{}, entry)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *SlaughterRepository) UpdateBatch(batch *SlaughterSaveBatch) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()
	args := []any{batch.UserId}
	animalArgs := []any{batch.UserId}
	valueTbl := make([]string, len(batch.Entries))
	animalValueTbl := make([]string, len(batch.Entries))

	for i, entry := range batch.Entries {
		rowArgs := []any{
			entry.Id,
			entry.AnimalId,
			entry.EntryDate,
			entry.DiscountRate,
			entry.Weight,
			entry.DeadWeight,
			entry.ButcherId,
		}
		valueTbl[i] = util.GenerateRowExpression(i + 2, rowArgs) 
		args = append(args, rowArgs...)

		animalRowArgs := []any{entry.AnimalId, entry.EntryDate}
		animalValueTbl[i] = util.GenerateRowExpression(i + 2, animalRowArgs)
		animalArgs = append(animalArgs, animalRowArgs...)
	}

	query := fmt.Sprintf(`
		WITH batch_entries(id, animal_id, entry_date, discount_rate, weight, dead_weight, butcher_id) AS (
			VALUES(%s)
		)
		UPDATE slaughter_entries 
		SET entry_date = be.entry_date,
			discount_rate = be.discount_rate,
			weight = be.weight,
			dead_weight = be.dead_weight,
			butcher_id = be.butcher_id
		FROM batch_entries be
		WHERE id = be.id AND user_id = $1
	`, strings.Join(valueTbl, ","))

	err = util.ExecTx(tx, query, args...)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		WITH batch_entries(animal_id, entry_date) AS (
			VALUES(%s)
		)
		UPDATE animals 
		SET death_date = be.entry_date
		FROM batch_entries be
		WHERE id = be.animal_id AND user_id = $1
	`, strings.Join(animalValueTbl, ","))

	err = util.ExecTx(tx, query, animalArgs...)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (r *SlaughterRepository) Add(entry *SlaughterSave) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	var query string
	if entry.Overwrite {
		query = `
			UPDATE slaughter_entries 
			SET discount_rate = :discount_rate, 
				weight = :weight, 
				dead_weight = :dead_weight
			WHERE entry_date = :entry_date
				AND animal_id = :animal_id
				AND user_id = :user_id
				AND deleted_at IS NULL
		`
	} else {
		query = `
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
	}

	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return err
	}

	animalQuery := `
		UPDATE animals
		SET death_date = :entry_date
		WHERE id = :animal_id 
			AND user_id = :user_id
	`
	err = util.NamedExecTx(tx, animalQuery, entry)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
