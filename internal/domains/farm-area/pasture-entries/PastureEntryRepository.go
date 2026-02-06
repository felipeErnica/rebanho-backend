package pastureEntries

import (
	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type PastureEntryRepository struct {
	SelectQuery string
	TableName   string
	DB          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PastureEntryRepository {
	selectQuery := `
        SELECT pastures_entries.*, 
            animals.name AS animal_name, animals.number AS animal_number,
            pastures.name AS pasture_name
        FROM pastures_entries
            LEFT JOIN animals ON animals.id = pastures_entries.animal_id
            LEFT JOIN pastures ON pastures.id = pastures_entries.pasture_id
    `
	return &PastureEntryRepository{selectQuery, "pastures_entries", db}
}

func (r *PastureEntryRepository) SearchPastureAnimals(pastureId string, userId string) (*[]entity.SearchEntity, error) {
	query := `
        SELECT
            animals.id,
            CONCAT_WS(' - ', animals.ring_number, animals.name, TO_CHAR(animals.birth_date, 'dd/mm/yyyy')) AS label
        FROM animals
            LEFT JOIN pastures ON pastures.id = animals.pasture_id
        WHERE 
			animals.deleted_at IS NULL 
            AND animals.pasture_id <> $1
            AND animals.user_id = $2
            AND animals.animal_type <> 'OUTSIDE_ANIMAL'
        ORDER BY 
			COALESCE(REGEXP_REPLACE(animals.ring_number, '[^0-9]', '', 'g')::int, 0), 
			COALESCE(animals.name, ''),
			COALESCE(animals.birth_date, '-infinity')
        LIMIT 20
    `
	return util.GetList[entity.SearchEntity](r.DB, query, pastureId, userId)
}

func (r *PastureEntryRepository) FindByPasture(
	pastureId string,
	userId string,
	filter *PastureEntryFilter,
	cursor string,
	sort string,
	order string,
) (*entity.Page[PastureEntry], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"animal_order":      {Field: "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":       {Field: "coalesce(animals.name, '')", Order: "asc"},
		"animal_birth_date": {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"entry_date":        {Field: "coalesce(pasture_entries.entry_date, '-infinity')", Order: "asc"},
		"id":                {Field: "pasture_entries.id", Order: "asc"},
		"created_at":        {Field: "pasture_entries.created_at", Order: "asc"},
	}

	sortExpression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT
            pasture_entries.id,
            pasture_entries.entry_date,
            COALESCE(REGEXP_REPLACE(animals.ring_number, '[^0-9]', '', 'g')::int, 0) AS animal_order,
            pasture_entries.animal_id,
            pasture_entries.created_at,
            animals.ring_number AS animal_ring_number, 
            animals.name AS animal_name, 
            animals.birth_date AS animal_birth_date,
            CONCAT_WS(' - ', mother.ring_number, mother.name) AS animal_mother,
            CONCAT_WS(' - ', father.ring_number, father.name) AS animal_father
        FROM pasture_entries
            LEFT JOIN animals ON animals.id = pasture_entries.animal_id
            LEFT JOIN animals AS mother ON mother.id = animals.mother_id
            LEFT JOIN animals AS father ON father.id = animals.father_id
    `
	whereExpression := `
        WHERE pasture_entries.deleted_at IS NULL 
            AND pasture_entries.pasture_id = $1
            AND pasture_entries.user_id = $2
    `

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := util.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression = whereExpression + " AND " + filterExpression
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if cursorExpression != "" {
		whereExpression = whereExpression + " AND " + cursorExpression
	}

	orderByExpression := " ORDER BY " + sortExpression
	query = query + whereExpression + orderByExpression

	args := []any{pastureId, userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	return util.GetPage[PastureEntry](r.DB, query, sort, 200, args...)
}

func (r *PastureEntryRepository) FindByPastureTotal(
	pastureId string,
	userId string,
	filter *PastureEntryFilter,
) (*PastureTotal, error) {

	query := "SELECT COUNT(pasture_entries.id) AS total FROM pasture_entries"
	whereExpression := `
        WHERE pasture_entries.deleted_at IS NULL 
            AND pasture_entries.pasture_id = $1
            AND pasture_entries.user_id = $2
    `

	filterExpression, _, err := util.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression = whereExpression + " AND " + filterExpression
	}

	query = query + whereExpression

	args := []any{pastureId, userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetOne[PastureTotal](r.DB, query, args...)
}

func (r *PastureEntryRepository) FindByAnimalId(animalId string) (*[]PastureEntry, error) {
	query := r.SelectQuery + " WHERE pastures_entries.deleted_at IS NULL AND pastures_entries.animal_id = $1"
	return util.GetList[PastureEntry](r.DB, query, animalId)
}

func (r *PastureEntryRepository) Delete(id string) error {
	return util.Delete(r.DB, r.TableName, id)
}

func (r *PastureEntryRepository) AddEntry(entry *PastureEntry) *log.APIError {

	validateErr := validateAddEntry(r.DB, *entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		INSERT INTO pasture_entries (animal_id, pasture_id, entry_date, user_id)
		VALUES(:animal_id, :pasture_id, :entry_date, :user_id)
	`

	err := util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *PastureEntryRepository) TransferEntry(entry *PastureEntry) *log.APIError {

	validateErr := validateTransferEntry(r.DB, *entry)
	if validateErr != nil {
		return validateErr
	}

	updateQuery := `
		UPDATE pasture_entries
		SET exit_date = :entry_date
		WHERE animal_id = :animal_id
			AND exit_date IS NULL
			AND user_id = :user_id
			AND deleted_at IS NULL
	`

	err := util.Exec(r.DB, updateQuery, entry.EntryDate, entry.AnimalId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	query := `
		INSERT INTO pasture_entries (animal_id, pasture_id, entry_date, user_id)
		VALUES(:animal_id, :pasture_id, :entry_date, :user_id)
	`

	err = util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (r *PastureEntryRepository) TransferCalfEntry(entry *PastureEntry) *log.APIError {

	cancelTransfer, validateErr := cancelChangeCalf(r.DB, *entry)
	if validateErr != nil {
		return validateErr
	}

	if cancelTransfer {
		return nil
	}

	validateErr = validateTransferCalf(r.DB, *entry)
	if validateErr != nil {
		return validateErr
	}

	updateQuery := `
		UPDATE pasture_entries
		SET exit_date = $1
		WHERE animal_id = $2
			AND exit_date IS NULL
			AND deleted_at IS NULL
	`

	err := util.Exec(r.DB, updateQuery, entry.EntryDate, entry.AnimalId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	query := `
		INSERT INTO pasture_entries (animal_id, pasture_id, entry_date, user_id)
		VALUES(:animal_id, :pasture_id, :entry_date, :user_id)
	`

	err = util.NamedExec(r.DB, query, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
