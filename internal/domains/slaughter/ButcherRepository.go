package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type ButcherRepository struct {
	DB *sqlx.DB
}

func newButcherRepository(db *sqlx.DB) *ButcherRepository {
	return &ButcherRepository{db}
}

func (r *ButcherRepository) FindAll(userId string) (*[]ButcherEntry, error) {

	query := `
		WITH entries_stats AS (
			SELECT
				butcher_id,
				COUNT(*) AS animals_number,
				AVG(dead_weight) AS avg_weight,
				AVG(COALESCE(dead_weight / NULLIF(weight * (1 - discount_rate), 0) * 100, 0)) AS avg_rate
			FROM slaughter_entries
			WHERE deleted_at IS NULL AND user_id = $1
			GROUP BY butcher_id
		)
		SELECT 
			b.id,
			b.name,
			b.cnpj,
			b.discount * 100 AS discount,
			b.address,
			s.animals_number,
			s.avg_weight,
			s.avg_rate
		FROM butchers b
			JOIN entries_stats s ON s.butcher_id = b.id
		WHERE b.user_id = $1 AND b.deleted_at IS NULL
		ORDER BY name
	`
	return util.GetList[ButcherEntry](r.DB, query, userId)
}

func (r *ButcherRepository) FindById(id string, userId string) (*ButcherEntry, error) {

	query := `
		SELECT 
			id,
			name,
			cnpj,
			discount * 100 AS discount,
			address
		FROM butchers 
		WHERE id = $1 
			AND user_id = $2 
			AND deleted_at IS NULL
	`
	return util.GetOne[ButcherEntry](r.DB, query, id, userId)
}

func (r *ButcherRepository) Search(userId string) (*[]entity.SearchEntity, error) {

	query := `
		SELECT
			s.id,
			s.name AS label
		FROM butchers s
		WHERE s.user_id = $1 AND s.deleted_at IS NULL
	`
	return util.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *ButcherRepository) FindEntriesPage(
	sort string,
	order string,
	cursor string,
	filter *SlaughterEntryFilter,
	butcherId string,
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
	query += whereExpression + orderExpression

	args := []any{userId, butcherId}
	filterArgs := util.GetFilterArgs(filter)
	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetPage[SlaughterEntry](r.DB, query, sort, 200, args...)
}

func (r *ButcherRepository) FindEntriesPageFoot(
	filter *SlaughterEntryFilter,
	butcherId string,
	userId string,
) (*SlaughterFoot, error) {

	query := `
		SELECT 
			COUNT(s.*) AS animals_number,
			AVG(s.weight) AS avg_weight,
			AVG(s.dead_weight) AS avg_dead_weight,
			AVG((s.dead_weight / NULLIF(weight * (1 - s.discount_rate), 0)) * 100) AS avg_rate
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

func (r *ButcherRepository) Add(entry *ButcherSave) *log.APIError {

	validateErr := validateButcherAdd(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		INSERT INTO butchers (name, cnpj, address, discount, user_id)
		VALUES (:name, :cnpj, :address, CAST(:discount AS float) / 100, :user_id)
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *ButcherRepository) Replace(entry *ButcherSave) *log.APIError {

	query := `
		UPDATE butchers 
		SET name = :name, 
			cnpj = :cnpj, 
			discount = CAST(:discount AS float) / 100, 
			address = :address
		WHERE name = :name
			OR cnpj = :cnpj 
			AND cnpj IS NOT NULL
			AND user_id = :user_id
			AND deleted_at IS NULL
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *ButcherRepository) Update(entry *ButcherSave) (*ButcherEntry, *log.APIError) {

	validateErr := validateButcherUpdate(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		UPDATE butchers 
		SET name = :name, 
			cnpj = :cnpj, 
			discount = CAST(:discount AS float) / 100, 
			address = :address
		WHERE id = :id AND user_id = :user_id
	`

	err = util.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	selectQuery := `
		WITH entries_stats AS (
			SELECT
				butcher_id,
				COUNT(*) AS animals_number,
				AVG(dead_weight) AS avg_weight,
				AVG(COALESCE(dead_weight / NULLIF(weight * (1 - discount_rate), 0) * 100, 0)) AS avg_rate
			FROM slaughter_entries
			GROUP BY butcher_id
		)
		SELECT 
			b.id,
			b.name,
			b.cnpj,
			b.discount * 100 AS discount,
			b.address,
			COALESCE(s.animals_number, 0) AS animals_number,
			COALESCE(s.avg_weight, 0) AS avg_weight,
			COALESCE(s.avg_rate, 0) AS avg_rate
		FROM butchers b 
			LEFT JOIN entries_stats s ON s.butcher_id = b.id
		WHERE b.id = :id AND b.user_id = :user_id
	`

	response, err := util.NamedGetTx(tx, selectQuery, ButcherEntry{}, entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *ButcherRepository) Delete(id string, userId string) *log.APIError {

	validateErr := validateButcherDelete(r.DB, id, userId)
	if validateErr != nil {
		return validateErr
	}

	query := `
		UPDATE butchers
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2
	`

	err := util.Exec(r.DB, query, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
