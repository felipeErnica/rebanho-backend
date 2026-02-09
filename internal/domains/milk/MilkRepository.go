package milk

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type MilkRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *MilkRepository {
	return &MilkRepository{db}
}

type SaveMilkValidation struct {
	HasLac             bool `db:"has_lac"`
	EntryExist         bool `db:"entry_exist"`
	IsDifferentPasture bool `db:"is_different_pasture"`
}

func (r *MilkRepository) CheckGroupUpdateConflicts(entry LactationGroupSave) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM milk_entries
			WHERE entry_date = :entry_date
				AND user_id = :user_id
				AND deleted_at IS NULL
		)
	`
	var exists bool
	err := util.NamedPrimitive(r.DB, query, &exists, entry)
	return exists, err
}

func (r *MilkRepository) CheckMilkEntryConflicts(entry *MilkEntrySave) (*SaveMilkValidation, error) {
	query := `
		SELECT 
			EXISTS (
				SELECT 1
				FROM lactations l
				WHERE l.animal_id = :animal_id
					AND l.start_date <= :entry_date
					AND :entry_date <= COALESCE(end_date, NOW())
					AND l.deleted_at IS NULL
					AND l.user_id = :user_id
			) AS has_lac,
			(
				SELECT id
				FROM milk_entries m
				WHERE m.animal_id = :animal_id
					AND m.entry_date = :entry_date
					AND id IS DISTINCT FROM :id
					AND m.user_id = :user_id
					AND m.deleted_at IS NULL
			) AS entry_exist,
			(
				SELECT COALESCE(pasture_id <> :pasture_id, FALSE)
				FROM pasture_entries
				WHERE animal_id = :animal_id 
					AND exit_date IS NULL
					AND user_id = :user_id
					AND deleted_at IS NULL
			) AS is_different_pasture 
	`
	return util.NamedGet(r.DB, query, SaveMilkValidation{}, entry)
}

func (r *MilkRepository) FindGroupsPage(
	filter *LactationGroupFilter,
	order string,
	cursor string,
	limit int,
	userId string,
) (*[]LactationGroup, error) {

	sortMap := map[string]util.SortField{
		"entry_date": {Field: "cte.entry_date", Order: "asc"},
	}

	query := `
		WITH cte AS (
			SELECT 
				entry_date,
				COUNT(*) animals_number,
				SUM(quantity) total_milk,
				AVG(quantity) avg_milk
			FROM milk_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
		)
		SELECT 
			cte.*,
			COALESCE(animals_number - LAG(animals_number) OVER (ORDER BY entry_date), 0) number_difference,
			COALESCE(((total_milk / LAG(total_milk) OVER (ORDER BY entry_date)) - 1)*100, 0) total_rate,
			COALESCE(((avg_milk / LAG(avg_milk) OVER (ORDER BY entry_date)) - 1)*100, 0) avg_rate
		FROM cte
    `
	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, "entry_date", order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression, cursorExpression)
	query += whereExpression + " ORDER BY cte.entry_date " + order + fmt.Sprintf(" LIMIT %d", limit)
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetList[LactationGroup](r.DB, query, args...)
}

func (r *MilkRepository) GetLastMilkEntries(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
			SELECT
				l.entry_date AS date,
				SUM(l.quantity) AS value
			FROM milk_entries l
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *MilkRepository) GetYearMilkEntries(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
			SELECT
				DATE_TRUNC('year', l.entry_date) AS date,
				SUM(l.quantity) AS date
			FROM milk_entries l
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 30
		)
		SELECT * FROM cte ORDER BY date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *MilkRepository) GetMilkProduction(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
			SELECT
				DATE_TRUNC('month', l.entry_date) AS date,
				SUM(l.quantity) AS value
			FROM milk_entries l
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 60
		)
		SELECT * FROM cte ORDER BY date
    `
	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *MilkRepository) GetLastAverageMilkEntries(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
			SELECT
				l.entry_date AS date,
				AVG(l.quantity) AS value
			FROM milk_entries l
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 10
		)
		SELECT * FROM cte ORDER BY date
    `

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *MilkRepository) GetYearAverageMilkEntries(userId string) (*[]util.GraphData, error) {
	query := `
        WITH cte AS (
			SELECT
				DATE_TRUNC('year', l.entry_date) AS date,
				AVG(l.quantity) AS value
			FROM milk_entries l
			WHERE l.user_id = $1 AND l.deleted_at IS NULL
			GROUP BY 1
			ORDER BY 1 DESC
			LIMIT 30
		)
		SELECT * FROM cte ORDER BY value
    `

	return util.GetList[util.GraphData](r.DB, query, userId)
}

func (r *MilkRepository) GetLactationEntries(lacId string) (*[]MilkDB, error) {
	query := `
		SELECT
			m.id,
			m.animal_id,
			CONCAT_WS(' - ', a.tag, a.name) AS animal_info,
			COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) AS animal_order,
			COALESCE(p.name, 'Sem Pasto') AS pasture_name,
			m.entry_date,
			m.quantity
		FROM milk_entries m
			JOIN animals a ON a.id = m.animal_id
			JOIN lactations l ON 
				l.id = $1
				AND m.entry_date >= l.start_date
				AND m.entry_date <=	COALESCE(l.end_date, NOW())
				AND m.animal_id = l.animal_id
			LEFT JOIN pasture_entries pe ON
				pe.animal_id = m.animal_id
				AND pe.entry_date <= m.entry_date
				AND m.entry_date <= COALESCE(pe.exit_date, NOW())
			LEFT JOIN pastures p ON p.id = pe.pasture_id
		WHERE m.deleted_at IS NULL
		ORDER BY m.entry_date
    `
	return util.GetList[MilkDB](r.DB, query, lacId)
}

func (r *MilkRepository) GetLactationEntriesFoot(lacId string) (*MilkEntryFoot, error) {
	query := `
		SELECT
			COUNT(*) animals_number,
			(EXTRACT(days FROM MAX(entry_date) - MIN(entry_date)) + 1)*AVG(quantity) total_milk,
			AVG(quantity) avg_milk
		FROM milk_entries m
			JOIN lactations l ON 
				l.id = $1
				AND m.entry_date >= l.start_date
				AND m.entry_date <=	COALESCE(l.end_date, NOW())
				AND m.animal_id = l.animal_id
		WHERE m.deleted_at IS NULL
    `
	return util.GetOne[MilkEntryFoot](r.DB, query, lacId)
}

func (r *MilkRepository) GetLastEntries(userId string) (*[]MilkDB, error) {
	query := `
		WITH max_tbl AS (
			SELECT MAX(entry_date) max_date 
			FROM milk_entries 
			WHERE user_id = $1 AND deleted_at IS NULL
		)

		SELECT 
			m.id,
			m.entry_date,
			m.quantity,
			
			m.animal_id,
			a.tag AS animal_tag,
			a.name AS animal_name,

			p.id AS pasture_id,
			p.name AS pasture_name,
			f.id AS farm_id,
			f.name AS farm_name
		FROM milk_entries m 
			CROSS JOIN max_tbl max
			JOIN animals a ON a.id = m.animal_id
			LEFT JOIN pasture_entries pe ON pe.animal_id = m.animal_id
				AND m.entry_date >= pe.entry_date
				AND m.entry_date < COALESCE(pe.exit_date, NOW())
				AND pe.deleted_at IS NULL
			LEFT JOIN pastures p ON p.id = pe.pasture_id
			LEFT JOIN farms f ON f.id = p.farm_id
		WHERE m.user_id = $1 
			AND m.deleted_at IS NULL 
			AND m.entry_date = max.max_date
		ORDER BY COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0)
	`
	return util.GetList[MilkDB](r.DB, query, userId)
}

func (r *MilkRepository) GetLastGroups(userId string) (*[]LactationGroup, error) {
	query := `
		WITH cte AS (
			SELECT 
				entry_date,
				COUNT(*) animals_number,
				SUM(quantity) total_milk,
				AVG(quantity) avg_milk
			FROM milk_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY 1
		)
		SELECT 
			cte.*,
			COALESCE(animals_number - LAG(animals_number) OVER (ORDER BY entry_date), 0) number_difference,
			COALESCE(((total_milk / LAG(total_milk) OVER (ORDER BY entry_date)) - 1)*100, 0) total_rate,
			COALESCE(((avg_milk / LAG(avg_milk) OVER (ORDER BY entry_date)) - 1)*100, 0) avg_rate
		FROM cte
		ORDER BY entry_date DESC
		LIMIT 5
	`
	return util.GetList[LactationGroup](r.DB, query, userId)
}

func (r *MilkRepository) FindPage(
	filter *MilkEntryFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*[]MilkDB, error) {

	sortMap := map[string]util.SortField{
		"animal_name":  {Field: "a.name", Order: "asc"},
		"animal_order": {Field: "coalesce(regexp_replace(a.tag, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"entry_date":   {Field: "m.entry_date", Order: "desc"},
		"quantity":     {Field: "m.quantity", Order: "asc"},
		"id":           {Field: "m.id", Order: "asc"},
		"created_at":   {Field: "m.created_at", Order: "asc"},
	}

	query := `
		SELECT
			m.id,
			m.entry_date,
			m.quantity,
			m.created_at,

			m.animal_id,
			a.tag AS animal_tag,
			a.name AS animal_name,
			COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0) AS animal_order,

			p.id AS pasture_id,
			p.name AS pasture_name,

			p.farm_id,
			f.name AS farm_name

		FROM milk_entries m
			JOIN animals a ON a.id = m.animal_id
			LEFT JOIN pasture_entries pe ON pe.animal_id = m.animal_id
				AND m.entry_date >= pe.entry_date
				AND m.entry_date < COALESCE(pe.exit_date, NOW())
				AND pe.deleted_at IS NULL
			LEFT JOIN pastures p ON p.id = pe.pasture_id
			LEFT JOIN farms f ON f.id = p.farm_id
    `

	whereExpression := "m.user_id = $1 AND m.deleted_at IS NULL"

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "m", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	whereExpression = util.GetWhereExpression(whereExpression, filterExpression, cursorExpression)
	query += whereExpression

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query += " ORDER BY " + sortExpression + fmt.Sprintf(" LIMIT %d", limit)
	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return util.GetList[MilkDB](r.DB, query, args...)
}

func (r *MilkRepository) GetPageFoot(filter *MilkEntryFilter, userId string) (*MilkEntryFoot, error) {
	query := `
		SELECT
			COUNT(*) animals_number,
			SUM(quantity) total_milk,
			AVG(quantity) avg_milk
		FROM milk_entries m
    `
	whereExpression := "WHERE m.user_id = $1 AND m.deleted_at IS NULL"

	filterExpression, _, err := util.GetFilterExpressions(filter, "m", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " AND " + filterExpression
	}

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query = query + whereExpression
	return util.GetOne[MilkEntryFoot](r.DB, query, args...)
}

func (r *MilkRepository) GetGroupEntries(userId string, entryDate time.Time) (*[]MilkDB, error) {

	query := `
		SELECT
			m.id,
			m.entry_date,
			m.quantity,

			m.animal_id,
			a.tag AS animal_tag,
			a.name AS animal_name,

			p.id AS pasture_id,
			p.name AS pasture_name,
			p.farm_id,
			f.name AS farm_name

		FROM milk_entries m
			JOIN animals a ON a.id = m.animal_id
			LEFT JOIN pasture_entries pe ON pe.animal_id = m.animal_id
				AND pe.entry_date <= m.entry_date
				AND COALESCE(pe.exit_date, NOW()) > m.entry_date
				AND pe.deleted_at IS NULL
			LEFT JOIN pastures p ON p.id = pe.pasture_id
			LEFT JOIN farms f ON f.id = p.farm_id
		WHERE m.entry_date = $2
			AND m.user_id = $1 
			AND m.deleted_at IS NULL 
		ORDER BY COALESCE(REGEXP_REPLACE(a.tag, '[^0-9]', '', 'g')::int, 0)
    `
	return util.GetList[MilkDB](r.DB, query, userId, entryDate)
}

func (r *MilkRepository) GetGroupEntriesFoot(userId string, entryDate time.Time) (*MilkEntryFoot, error) {
	query := `
		SELECT
			COUNT(*) AS animals_number,
			SUM(quantity) AS total_milk,
			AVG(quantity) AS avg_milk
		FROM milk_entries m
		WHERE m.user_id = $1 AND deleted_at IS NULL 
			AND m.entry_date = $2
    `
	return util.GetOne[MilkEntryFoot](r.DB, query, userId, entryDate)
}

func (r *MilkRepository) UpdateGroup(groupEntry *LactationGroupSave) (*LactationGroup, error) {

	query := `
		UPDATE milk_entries
		SET entry_date = :entry_date
		WHERE entry_date = :old_entry
			AND user_id = :user_id
			AND deleted_at IS NULL
	`
	err := util.NamedExec(r.DB, query, groupEntry)
	if err != nil {
		return nil, err
	}

	returnQuery := `
		WITH milk_stats AS (
			SELECT 
				entry_date,
				COUNT(*) animals_number,
				SUM(quantity) total_milk,
				AVG(quantity) avg_milk
			FROM milk_entries
			WHERE user_id = :user_id AND deleted_at IS NULL
			GROUP BY 1
		),
		cte AS (
			SELECT 
				s.*,
				COALESCE(animals_number - LAG(animals_number) OVER (ORDER BY entry_date), 0) AS number_difference,
				COALESCE(((total_milk / LAG(total_milk) OVER (ORDER BY entry_date)) - 1) * 100, 0) AS total_rate,
				COALESCE(((avg_milk / LAG(avg_milk) OVER (ORDER BY entry_date)) - 1) * 100, 0) AS avg_rate
			FROM milk_stats s
		)
		SELECT * 
		FROM cte 
		WHERE entry_date = :entry_date 
			AND user_id = :user_id 
			AND deleted_at IS NULL
	`
	response, err := util.NamedGet(r.DB, returnQuery, LactationGroup{}, groupEntry)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *MilkRepository) DeleteGroup(entryDate time.Time, userId string) error {
	deleteQuery := `
		UPDATE milk_entries
		SET deleted_at = NOW()
		WHERE entry_date = $1 AND user_id = $2
	`

	err := util.Exec(r.DB, deleteQuery, entryDate, userId)
	if err != nil {
		return err
	}

	return nil
}

func (r *MilkRepository) Add(entry *MilkEntrySave) error {

	tx, err := r.DB.Beginx()
	if err != nil {
		return err
	}

	defer tx.Rollback()

	if entry.TransferPasture {
		updateQuery := `
			UPDATE pasture_entries
			SET exit_date = :entry_date
			WHERE animal_id = :animal_id
				AND exit_date IS NULL
				AND user_id = :user_id
				AND deleted_at IS NULL
		`
		err = util.NamedExecTx(tx, updateQuery, entry)
		if err != nil {
			return err
		}

		insertQuery := `
			INSERT INTO pasture_entries (animal_id, entry_date, user_id)
			VALUES (:animal_id, :entry_date, :user_id)
		`
		err = util.NamedExecTx(tx, insertQuery, entry)
		if err != nil {
			return err
		}
	}

	if entry.Overwrite {
		query := `
			UPDATE milk_entries
			SET quantity = :quantity
			WHERE animal_id = :animal_id 
				AND entry_date = :entry_date
				AND user_id = :user_id
		`
		err = util.NamedExecTx(tx, query, entry)
		if err != nil {
			return err
		}
	} else {
		insertQuery := `
			INSERT INTO milk_entries (animal_id, entry_date, quantity, user_id) 
			VALUES (:animal_id, :entry_date, :quantity, :user_id)
		`
		err = util.NamedExecTx(tx, insertQuery, entry)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil

}

func (r *MilkRepository) Update(entry *MilkEntrySave) (*MilkDB, error) {

	query := `
		UPDATE milk_entries 
		SET entry_date = :entry_date,
			quantity = :quantity
		WHERE id = :id AND user_id = :user_id
	`
	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, err
	}

	returnQuery := `
		SELECT 
			m.id,
			m.animal_id,
			CONCAT_WS(' - ', a.tag, a.name) AS animal_info,
			COALESCE(p.name, 'Sem Pasto') AS pasture_name,
			m.entry_date,
			m.quantity
		FROM milk_entries m 
			JOIN animals a ON a.id = m.animal_id
			LEFT JOIN pasture_entries pe ON 
				pe.animal_id = m.animal_id
				AND pe.entry_date <= m.entry_date
				AND m.entry_date <= COALESCE(pe.exit_date, NOW())
			LEFT JOIN pastures p ON p.id = pe.pasture_id
			WHERE m.id = :id
				AND user_id = :user_id
				AND deleted_at IS NULL
	`

	response, err := util.NamedGet(r.DB, returnQuery, MilkDB{}, entry)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (r *MilkRepository) Delete(id string) error {

	deleteQuery := `
		UPDATE milk_entries
		SET deleted_at = NOW()
		WHERE id = $1
	`

	err := util.Exec(r.DB, deleteQuery, id)
	if err != nil {
		return err
	}

	return nil
}
