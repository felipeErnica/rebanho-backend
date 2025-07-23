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

func (r *PastureEntryRepository) FindByPasture(
	pastureId string,
	userId string,
    filter PastureEntryFilter,
	cursor string,
	sort string,
	order string,
) (*entity.Page[PastureEntry], error) {

	sortMap := map[string]string{
		"animal_order":      "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)",
		"animal_name":       "coalesce(animals.name, '')",
		"animal_birth_date": "coalesce(animals.birth_date, '-infinity')",
		"entry_date":        "coalesce(pasture_entries.entry_date, '-infinity')",
	}

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	countQuery := `
        select count(pasture_entries.id) from pasture_entries 
        left join animals on animals.id = pasture_entries.animal_id
    `

	query := `
        select
            pasture_entries.id,
            pasture_entries.entry_date,
            coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
            pasture_entries.animal_id,
            pasture_entries.created_at,
            animals.ring_number as animal_ring_number, 
            animals.name as animal_name, 
            animals.birth_date as animal_birth_date
        from pasture_entries
            left join animals on animals.id = pasture_entries.animal_id

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

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, "pasture_entries", cursorArgs, 3)
	if err != nil {
		return nil, err
	}
	if cursorExpression != "" {
		cursorExpression = " and " + cursorExpression
	}

	countQuery = countQuery + whereExpression
	whereExpression = whereExpression + cursorExpression
	orderByExpression := " order by " + sortExpression
	query = query + whereExpression + orderByExpression

	return repositoriesUtil.GetPage[PastureEntry](r.Db, query, countQuery, sort, 200, cursorArgs, pastureId, userId)
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
