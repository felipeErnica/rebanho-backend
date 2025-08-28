package pastureEntries

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type PastureEntryRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PastureEntryRepository {
	selectQuery := `
        SELECT pastures_entries.*, 
            animals.name as animal_name, animals.number as animal_number,
            pastures.name as pasture_name
        FROM pastures_entries
            LEFT JOIN animals ON animals.id = pastures_entries.animal_id
            LEFT JOIN pastures ON pastures.id = pastures_entries.pasture_id
    `
	return &PastureEntryRepository{selectQuery, "pastures_entries", db}
}

func (r *PastureEntryRepository) SearchPastureAnimals(pastureId string, userId string) (*[]entity.SearchEntity, error) {
	query := `
        select
            animals.id,
            concat_ws(' - ', animals.ring_number, animals.name, to_char(animals.birth_date, 'dd/mm/yyyy')) as label
        from animals
            left join pastures on pastures.id = animals.pasture_id
        where 
			animals.deleted_at is null 
            and animals.pasture_id <> $1
            and animals.user_id = $2
            and animals.animal_type <> 'OUTSIDE_ANIMAL'
        order by 
			coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0), 
			coalesce(animals.name, ''),
			coalesce(animals.birth_date, '-infinity')
        limit 20
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.Db, query, pastureId, userId)
}

func (r *PastureEntryRepository) FindByPasture(
	pastureId string,
	userId string,
	filter PastureEntryFilter,
	cursor string,
	sort string,
	order string,
) (*entity.Page[PastureEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_order":      {Field: "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":       {Field: "coalesce(animals.name, '')", Order: "asc"},
		"animal_birth_date": {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"entry_date":        {Field: "coalesce(pasture_entries.entry_date, '-infinity')", Order: "asc"},
		"id":                {Field: "pasture_entries.id", Order: "asc"},
		"created_at":        {Field: "pasture_entries.created_at", Order: "asc"},
	}

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query := `
        select
            pasture_entries.id,
            pasture_entries.entry_date,
            coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
            pasture_entries.animal_id,
            pasture_entries.created_at,
            animals.ring_number as animal_ring_number, 
            animals.name as animal_name, 
            animals.birth_date as animal_birth_date,
            concat_ws(' - ', mother.ring_number, mother.name) as animal_mother,
            concat_ws(' - ', father.ring_number, father.name) as animal_father
        from pasture_entries
            left join animals on animals.id = pasture_entries.animal_id
            left join animals as mother on mother.id = animals.mother_id
            left join animals as father on father.id = animals.father_id
    `
	whereExpression := `
        where pasture_entries.deleted_at is null 
            and pasture_entries.pasture_id = $1
            and pasture_entries.user_id = $2
    `

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression = whereExpression + " and " + filterExpression
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if cursorExpression != "" {
		whereExpression = whereExpression + " and " + cursorExpression
	}

	orderByExpression := " order by " + sortExpression
	query = query + whereExpression + orderByExpression

	args := []any{pastureId, userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	return repositoriesUtil.GetPage[PastureEntry](r.Db, query, sort, 200, args...)
}

func (r *PastureEntryRepository) FindByPastureTotal(
	pastureId string,
	userId string,
	filter PastureEntryFilter,
) (*PastureTotal, error) {

	query := "select count(pasture_entries.id) as total from pasture_entries"
	whereExpression := `
        where pasture_entries.deleted_at is null 
            and pasture_entries.pasture_id = $1
            and pasture_entries.user_id = $2
    `

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression = whereExpression + " and " + filterExpression
	}

	query = query + whereExpression

	args := []any{pastureId, userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[PastureTotal](r.Db, query, args...)
}

func (r *PastureEntryRepository) FindByAnimalId(animalId string) (*[]PastureEntry, error) {
	query := r.SelectQuery + " WHERE pastures_entries.deleted_at is null and pastures_entries.animal_id = $1"
	return repositoriesUtil.GetList[PastureEntry](r.Db, query, animalId)
}

func (r *PastureEntryRepository) Add(newEntry *PastureEntrySave) (*PastureEntrySave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newEntry)
}

func (r *PastureEntryRepository) Update(entry *PastureEntrySave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, entry)
}

func (r *PastureEntryRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
