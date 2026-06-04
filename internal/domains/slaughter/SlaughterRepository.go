package slaughter

import (
	"fmt"

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
			WHERE group_id = :group_id
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
				g.entry_date AS date,
				AVG((s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0)) ) AS value
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
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
				g.entry_date AS date,
				AVG(s.weight) AS value
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
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
				g.entry_date AS date,
				AVG(g.dead_weight) AS value
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = g.group_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
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
				DATE_TRUNC('month', g.entry_date) entry_date,
				AVG(s.weight) avg_weight,
				AVG(s.dead_weight) dead_weight
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
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
				DATE_TRUNC('month', g.entry_date) AS date,
				AVG((s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0))) AS value
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
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
				COUNT(*) AS animals_number,
				AVG(weight) AS avg_weight,
				AVG((s.dead_weight / NULLIF(weight * (1 - g.discount), 0))) AS avg_rate
			FROM slaughter_entries s 
			JOIN animals a ON a.id = s.animal_id
			JOIN slaughter_groups g ON g.id = s.group_id
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
				COUNT(*) AS animals_number,
				AVG(weight) AS avg_weight,
				AVG((s.dead_weight / NULLIF(weight * (1 - g.discount), 0))) AS avg_rate
			FROM slaughter_entries s 
			JOIN animals a ON a.id = s.animal_id
			JOIN slaughter_groups g ON g.id = s.group_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
			GROUP BY 1
			HAVING COUNT(*) >= 10
		),
		gn_stats AS (
			SELECT 
				AVG(avg_weight) AS gn_avg_weight,
				AVG(avg_rate) AS gn_avg_rate
			FROM cte
		)
		SELECT 
			CONCAT_WS(' - ', a.tag, a.name) AS name,
			c.avg_weight,
			((c.avg_weight / NULLIF(s.gn_avg_weight, 0)) - 1) AS weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(s.gn_avg_rate, 0)) - 1) AS rate_comparison,
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
				g.butcher_id,
				COUNT(*) AS animals_number,
				AVG(s.weight) AS avg_weight,
				AVG((s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0))) AS avg_rate
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
		),
		gn_stats AS (
			SELECT 
				AVG(avg_weight) AS gn_avg_weight,
				AVG(avg_rate) AS gn_avg_rate
			FROM cte
		)
		SELECT 
			s.name,
			c.avg_weight,
			((c.avg_weight / NULLIF(g.gn_avg_weight, 0)) - 1) AS weight_comparison,
			c.avg_rate,
			((c.avg_rate / NULLIF(g.gn_avg_rate, 0)) - 1) AS rate_comparison,
			c.animals_number
		FROM gn_stats g
		CROSS JOIN cte c 
		JOIN butchers s ON s.id = c.butcher_id
		WHERE c.avg_rate >= 20
		ORDER BY c.avg_rate DESC
		LIMIT 10
	`
	return util.GetList[TableRatings](r.DB, query, userId)
}

func (r *SlaughterRepository) GetLastEntries(userId string) (*[]SlaughterDB, error) {
	query := `

		WITH last_date AS (
			SELECT MAX(entry_date) AS max_date
			FROM slaughter_groups
			WHERE user_id = $1 AND deleted_at IS NULL
		),

		last_groups AS (
			SELECT 
				g.id AS group_id,
				g.entry_date AS group_date,
				g.discount AS group_discount
				
				g.butcher_id,
				b.name AS butcher_name
			FROM slaughter_groups g
			CROSS JOIN last_date ld
			JOIN butchers b ON b.id = g.butcher_id
			WHERE g.entry_date = ld.max_date
				AND g.user_id = $1 
				AND g.deleted_at IS NULL
		)

		SELECT 
			s.id,
			s.weight,
			s.weight * (1 - g.group_discount) discount_weight,
			s.dead_weight,
			s.dead_weight / NULLIF(s.weight*(1 - s.discount_rate), 0) AS performance_rate,

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

			g.group_id,
			g.group_date,
			g.group_discount,

			g.butcher_id,
			g.butcher_name
		FROM slaughter_entries s
		CROSS JOIN last_group l
		LEFT JOIN animals a ON a.id = s.animal_id
		LEFT JOIN animals m ON m.id = a.mother_id
		LEFT JOIN animals f ON f.id = a.father_id
		WHERE s.user_id = $1 AND s.deleted_at IS NULL
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
		"entry_date":       {Field: "cte.group_date", Order: "desc"},
		"animal_order":     {Field: "COALESCE(NULLIF(regexp_replace(cte.animal_tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "COALESCE(cte.animal_name, '')", Order: "asc"},
		"birth_date":       {Field: "COALESCE(cte.animal_birth, '-infinity')", Order: "desc"},
		"weight":           {Field: "cte.weight", Order: "asc"},
		"dead_weight":      {Field: "cte.dead_weight", Order: "asc"},
		"performance_rate": {Field: "COALESCE(cte.dead_weight / NULLIF(cte.weight*(1 - cte.group_discount), 0) , 0)", Order: "asc"},
		"id":               {Field: "cte.id", Order: "asc"},
		"created_at":       {Field: "cte.created_at", Order: "asc"},
	}

	query := `
		WITH cte AS (
			SELECT 
				s.id,
				s.weight,
				s.weight * (1 - g.discount) AS discount_weight,
				s.dead_weight,
				s.dead_weight / NULLIF(s.weight*(1 - g.discount), 0) AS performance_rate,

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

				s.group_id,
				g.entry_date AS group_date,
				g.discount AS group_discount,

				g.butcher_id,
				b.name butcher_name,

				s.created_at
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
			JOIN butchers b ON b.id = g.butcher_id
			LEFT JOIN animals a ON a.id = s.animal_id
			LEFT JOIN animals f ON f.id = a.father_id
			LEFT JOIN animals m ON m.id = a.mother_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
		)
		SELECT * FROM cte
	`

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression, cursorExpression)

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
		WITH cte AS (
			SELECT 
				s.id,
				s.weight,
				s.weight * (1 - g.discount) AS discount_weight,
				s.dead_weight,
				s.dead_weight / NULLIF(s.weight*(1 - g.discount), 0) AS performance_rate,

				s.animal_id,
				a.sex AS animal_sex,
				a.birth_date AS animal_birth,

				a.father_id,
				a.mother_id,

				s.group_id,
				g.entry_date AS group_date,
				g.discount AS group_discount,
				g.butcher_id,
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
			JOIN butchers b ON b.id = g.butcher_id
			LEFT JOIN animals a ON a.id = s.animal_id
			WHERE s.user_id = $1 AND s.deleted_at IS NULL
		)

		SELECT 
			COUNT(s.*) AS animals_number,
			AVG(cte.weight) AS avg_weight,
			AVG(cte.dead_weight) AS avg_dead_weight,
			AVG((cte.dead_weight / NULLIF(cte.weight * (1 - cte.group_discount), 0)) ) AS avg_rate
		FROM cte
	`

	filterExpression, _, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression)

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
		"entry_date":       {Field: "cte.group_date", Order: "desc"},
		"animal_order":     {Field: "COALESCE(NULLIF(regexp_replace(cte.animal_tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "COALESCE(cte.animal_name, '')", Order: "asc"},
		"birth_date":       {Field: "COALESCE(cte.animal_birth, '-infinity')", Order: "desc"},
		"weight":           {Field: "cte.weight", Order: "asc"},
		"dead_weight":      {Field: "cte.dead_weight", Order: "asc"},
		"performance_rate": {Field: "COALESCE(cte.dead_weight / NULLIF(cte.weight*(1 - cte.group_discount), 0) , 0)", Order: "asc"},
		"id":               {Field: "s.id", Order: "asc"},
		"created_at":       {Field: "s.created_at", Order: "asc"},
	}

	query := `
		WITH cte AS (
			SELECT 
				s.id,
				s.weight,
				s.weight * (1 - g.discount) discount_weight,
				s.dead_weight,
				s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0) AS performance_rate,

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

				s.group_id,
				g.entry_date AS group_date,
				g.discount AS group_discount,

				g.butcher_id,
				h.name butcher_name,

				s.created_at
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
			JOIN butchers h ON h.id = g.butcher_id
			LEFT JOIN animals a ON a.id = s.animal_id
			LEFT JOIN animals f ON f.id = a.father_id
			LEFT JOIN animals m ON m.id = a.mother_id
			WHERE s.user_id = $1 
				AND g.butcher_id = $2
				AND s.deleted_at IS NULL
		)
	`


	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "cte", 3)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression, cursorExpression)

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
		WITH cte AS (
			SELECT
				s.weight,
				s.dead_weight,
				s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0) AS performance_rate,
				s.animal_id,
				a.birth_date AS animal_birth,
				a.father_id,
				a.mother_id,
				g.entry_date AS group_date,
				g.butcher_id
			FROM slaughter_entries s 
			JOIN slaughter_groups g ON g.id = s.group_id
			WHERE s.user_id = $1 AND g.butcher_id = $2 AND s.deleted_at IS NULL
		)

		SELECT 
			COUNT(cte.*) AS animals_number,
			AVG(cte.weight) AS avg_weight,
			AVG(cte.dead_weight) AS avg_dead_weight,
			AVG(cte.performance_rate) AS avg_rate
		FROM cte 
	`

	filterExpression, _, err := util.GetFilterExpressions(filter, "cte", 3)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression)

	args := []any{userId, butcherId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query += whereExpression

	return util.GetOne[SlaughterFoot](r.DB, query, args...)
}

func (r *SlaughterRepository) FindEntriesByGroup(
	sort string,
	order string,
	filter *SlaughterFilter,
	groupId string,
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
				s.weight,
				s.weight * (1 - g.discount) AS discount_weight,
				s.dead_weight,
				(s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0)) AS performance_rate,

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
			FROM slaughter_entries s 
			JOIN slaughter_groups g ON g.id = s.group_id
			JOIN butchers h ON h.id = s.butcher_id
			LEFT JOIN animals a ON a.id = s.animal_id
			LEFT JOIN animals m ON m.id = a.mother_id
			LEFT JOIN animals f ON f.id = a.father_id
			WHERE s.user_id = $1 
				AND s.deleted_at IS NULL
				AND s.group_id = $2
		)
		SELECT * FROM cte
	`

	filterExp, _, err := util.GetFilterExpressions(filter, "cte", 3)
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
	args := []any{userId, groupId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetList[SlaughterDB](r.DB, query, args...)
}

func (r *SlaughterRepository) GetEntriesFoot(filter *SlaughterFilter, groupId string, userId string) (*SlaughterFoot, error) {
	query := `
		WITH cte AS (
			SELECT 
				g.entry_date,
				s.animal_id,
				s.weight,
				s.dead_weight,
				(s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0)) AS rate
			FROM slaughter_entries s
			JOIN slaughter_groups g ON g.id = s.group_id
			WHERE s.user_id = $1 
				AND s.group_id = $2
				AND s.deleted_at IS NULL
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

	query := `
		UPDATE animals a
		SET death_date = NULL
		FROM slaughter_entries se
		WHERE se.id = $1
			AND a.id = se.animal_id 
			AND se.user_id = $2
	`
	err = util.ExecTx(tx, query, id, userId)
	if err != nil {
		return err
	}

	deleteQuery := `
		UPDATE slaughter_entries 
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`
	err = util.ExecTx(tx, deleteQuery, id, userId)
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
		UPDATE animals a
		SET death_date = NULL
		FROM (%s) b
		WHERE a.id = b.animal_id AND a.user_id = b.user_id
	`, batchQuery)

	err = util.ExecTx(tx, query, args...)
	if err != nil {
		return err
	}

	query = fmt.Sprintf(`
		UPDATE slaughter_entries s
		SET deleted_at = NOW()
		FROM (%s) b
		WHERE s.id = b.id AND s.user_id = b.user_id
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
		SET group_id = :group_id,
			weight = :weight,
			dead_weight = :dead_weight,
		WHERE id = :id AND user_id = :user_id
	`

	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, err
	}

	query = `
		UPDATE animals a
		SET death_date = g.entry_date
		FROM slaughter_groups g
		WHERE a.id = :animal_id 
			AND g.id = :group_id 
			AND a.user_id = :user_id
	`
	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, err
	}

	query = `
		SELECT 
			s.id,
			s.weight,
			s.weight * (1 - g.discount) AS discount_weight,
			s.dead_weight,
			s.dead_weight / NULLIF(s.weight * (1 - g.discount), 0) AS performance_rate,

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

			g.entry_date AS group_date,
			g.discount AS group_discount
		FROM slaughter_entries s
		JOIN slaughter_groups g ON g.id = s.group_id
		JOIN butchers b ON b.id = g.butcher_id
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

func (r *SlaughterRepository) UpdateBatch(batch *[]SlaughterSave) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	//Altera as datas de morte dos animais antes de modificar o abate
	query := `
		UPDATE animals a
		SET death_date = g.entry_date
		FROM slaughter_groups g
		WHERE a.id = :animal_id 
			AND g.id = :group_id
			AND a.user_id = :user_id
	`

	err = util.NamedExecTx(tx, query, batch)
	if err != nil {
		return err
	}

	query = `
		UPDATE slaughter_entries 
		SET group_id = :group_id,
			weight = :weight,
			dead_weight = :dead_weight
		WHERE id = :id AND user_id = :user_id
	`

	err = util.NamedExecTx(tx, query, batch)
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

	//Altera as datas de morte dos animais antes de adicionar
	query := `
		UPDATE animals a
		SET death_date = g.entry_date
		FROM slaughter_groups g
		WHERE a.id = :animal_id 
			AND g.id = :group_id
			AND a.user_id = :user_id
	`
	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return err
	}

	if entry.Overwrite {
		query = `
			UPDATE slaughter_entries 
			SET weight = :weight, 
				dead_weight = :dead_weight
				WHERE group_id = :group_id
				AND animal_id = :animal_id
				AND user_id = :user_id
				AND deleted_at IS NULL
		`
	} else {
		query = `
			INSERT INTO slaughter_entries (
				animal_id, 
				group_id,
				weight, 
				dead_weight,
				user_id
			)
			VALUES (
				:animal_id,
				:group_id,
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

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
