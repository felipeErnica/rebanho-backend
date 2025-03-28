package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PastureEntryRepository struct {
	Base        PageRepositoryImpl[entity.PastureEntry]
	SelectQuery util.QueryConstructor
}

func (r *PastureEntryRepository) Init() {

	dateFields := []string{
		"entry_date",
		"exit_date",
		"created_at",
	}

	selectQuery := new(util.QueryConstructor).Select("entries", "id", "entry_date", "exit_date")
	selectQuery.AndSelect("animals", "id", "identification_number", "animal_order", "name")
	selectQuery.AndSelect("pastures", "id", "identification_number", "animal_order", "name")
	selectQuery.From("pasture_entries", "entries")
	selectQuery.LeftJoin("pastures", "").On("pastures.id", "entries.pasture_id")
	selectQuery.LeftJoin("animals", "").On("animals.id", "entries.animal_id")
	r.SelectQuery = *selectQuery

    insertQuery := new(util.QueryConstructor).Insert("pasture_entries", "id", "entry_date", "exit_date", "animal_id", "pasture_id", 
        "created_at", "user_id")
    updateQuery := new(util.QueryConstructor).Update("pasture_entries", "id", "entry_date", "exit_date", "animal_id", "pasture_id", 
        "created_at", "user_id")

	mainRepository := RepositoryImpl[entity.PastureEntry]{
		TableName:       "entries",
		SelectQueryBody: r.SelectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		Repository:      r,
	}

	r.Base = PageRepositoryImpl[entity.PastureEntry]{
		PageRepository: r,
		Base:           &mainRepository,
		DateFields:     dateFields,
	}
}

func (r *PastureEntryRepository) getFields(sort string) (firstField string, secondField string) {
	switch sort {
	case "name":
		return "animal.name", "entries.id"
	case "identification_number":
		return "animal.animal_order", "entries.id"
	case "entry_date":
		return "entries.entry_date", "entries.id"
	case "exit_date":
		return "entries.exit_date", "entries.id"
	default:
		return "entries.created_at", "entries.id"
	}
}

func (r *PastureEntryRepository) createKey(sort string, lastEntry *entity.PastureEntry) (key string) {
	switch sort {
	case "name":
		return fmt.Sprintf("%s,%s", *lastEntry.Animal.Name, lastEntry.Id)
	case "identification_number":
		return fmt.Sprintf("%d,%s", lastEntry.Animal.AnimalOrder, lastEntry.Id)
	case "entry_date":
		return fmt.Sprintf("%s,%s", lastEntry.EntryDate, lastEntry.Id)
	case "exit_date":
		return fmt.Sprintf("%s,%s", lastEntry.ExitDate, lastEntry.Id)
	default:
		return fmt.Sprintf("%s,%s", lastEntry.CreatedAt, lastEntry.Id)
	}
}

func (r *PastureEntryRepository) setNewEntity(model *entity.PastureEntry, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
    model.UserId = GetUserId()
}

func (r *PastureEntryRepository) buildEntity(row *sql.Row) (*entity.PastureEntry, error) {
	var entry entity.PastureEntry
	err := row.Scan(&entry.Id, &entry.EntryDate, &entry.ExitDate,
		&entry.Animal.Id, &entry.Animal.IdentificationNumber, &entry.Animal.AnimalOrder, &entry.Animal.Name,
		&entry.Pasture.Id, &entry.Pasture.Name)
	return &entry, err
}

func (r *PastureEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.PastureEntry, err error) {
	var entries []entity.PastureEntry
	for rows.Next() {
		var entry entity.PastureEntry
		err = rows.Scan(&entry.Id, &entry.EntryDate, &entry.ExitDate,
			&entry.Animal.Id, &entry.Animal.IdentificationNumber, &entry.Animal.AnimalOrder, &entry.Animal.Name,
			&entry.Pasture.Id, &entry.Pasture.Name)
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}
	return &entries, err
}

func (r *PastureEntryRepository) saveOrUpdateScan(query string, model *entity.PastureEntry) error {
	return execQuery(query, model.Id, model.EntryDate, model.ExitDate, model.Animal.Id,
		model.Pasture.Id, model.CreatedAt, model.UserId)
}

func (r *PastureEntryRepository) FindByPastureId(sort string, direction string,
	cursor string, pastureId string) (*entity.Page[entity.PastureEntry], error) {
    query := r.SelectQuery.Where("entries.pasture_id = $1")
	return r.Base.FindRandomQueryPage(query, sort, direction, cursor, pastureId)
}

func (r *PastureEntryRepository) FindByAnimalId(animalId string) (*[]entity.PastureEntry, error) {
    query := r.SelectQuery.Where("entries.deleted_at is null").And("entries.pasture_id = $1")
	return r.Base.FindListByQuery(query, animalId)
}

func (r *PastureEntryRepository) FindByDeletedPasturePage(sort string, direction string,
	cursor string, pastureId string) (*entity.Page[entity.PastureEntry], error) {
    query := r.SelectQuery.Where("entries.deleted_at is not null").And("entries.pasture_id = $1")
	return r.Base.FindRandomQueryPage(query, sort, direction, cursor, pastureId)
}

func (r *PastureEntryRepository) Add(newEntry *entity.PastureEntry) (*entity.PastureEntry, error) {
	return r.Base.Add(newEntry)
}

func (r *PastureEntryRepository) Save(entry *entity.PastureEntry) error {
	return r.Base.Save(entry)
}

func (r *PastureEntryRepository) Delete(id string) error {
	return r.Base.Delete(id)
}
