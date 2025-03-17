package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/google/uuid"
)

type PastureEntryRepository struct {
    Base PageRepositoryImpl[entity.PastureEntry]
}

func (r *PastureEntryRepository) Init() {

    dateFields:= []string {
        "entry_date",
        "exit_date",
        "created_at",
    }

    selectQuery:=`
        SELECT pasture_entries.id, pasture_entries.entry_date, pasture_entries.exit_date,
            animal.id as animal_id, animal.identification_number as animal_number, animal.animal_order as animal_order,
            animal.name,
            pasture.id, pasture.name
        FROM pasture_entries_active as pasture_entries
        LEFT JOIN pastures as pasture ON pasture.id = pasture_entries.pasture_id
        LEFT JOIN animals as animal ON animal.id = pasture_entries.animal_id
    `

    insertQuery:=`
        INSERT INTO pasture_entries (id, entry_date, exit_date, animal_id, pasture_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6)
    `

    updateQuery:=`
        UPDATE pasture_entries 
        SET entry_date = $2, exit_date = $3, animal_id = $4, pasture_id = $5, created_at =$6
        WHERE id = $1
    `


    mainRepository:= RepositoryImpl[entity.PastureEntry] {
        TableName: "pasture_entries",
        SelectQueryBody: selectQuery,
        InsertQuery: insertQuery,
        UpdateQuery: updateQuery,
        Repository: r,
    }

    r.Base = PageRepositoryImpl[entity.PastureEntry]{
        PageRepository: r,
        Base: &mainRepository,
        DateFields: dateFields,
    }
}

func (r *PastureEntryRepository) getFields(sort string) (firstField string, secondField string) {
	switch sort {
	case "name":
		return "animal.name", "pasture_entries.id"
	case "identification_number":
		return "animal.animal_order", "pasture_entries.id"
	case "entry_date":
		return "pasture_entries.entry_date", "pasture_entries.id"
	case "exit_date":
		return "pasture_entries.exit_date", "pasture_entries.id"
	default:
		return "pasture_entries.created_at", "pasture_entries.id"
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

func (r *PastureEntryRepository) buildPage(query string, sort string, args... any) (page *entity.Page[entity.PastureEntry], err error) {
	rows, err := selectQueryList(query, args...)
	if err != nil {
		return
	}

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
	rows.Close()

	nextCursor, err := r.Base.CreateNextCursor(sort, entries)
	if err != nil {
		return
	}

	pageAnimal := entity.Page[entity.PastureEntry]{
		HasNextPage: len(entries) == PAGE_LIMIT,
		NextCursor:  nextCursor,
		List:        &entries,
	}

	return &pageAnimal, err
}

func (r *PastureEntryRepository) setNewEntity(model *entity.PastureEntry) {
    model.Id = uuid.NewString()
    model.CreatedAt = time.Now()
}

func (r *PastureEntryRepository) buildEntity(row *sql.Row) (*entity.PastureEntry, error) {
    var entry entity.PastureEntry
    err:= row.Scan(&entry.Id, &entry.EntryDate, &entry.ExitDate,
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
        model.Pasture.Id, model.CreatedAt)
}

func (r *PastureEntryRepository) FindByPastureId(sort string, direction string, 
    cursor string, pastureId string) (*entity.Page[entity.PastureEntry], error) {
    return r.Base.FindPageCondional(sort, direction, cursor, "pasture_entries.pasture_id", pastureId)
}

func (r *PastureEntryRepository) FindByAnimalId(animalId string) (*[]entity.PastureEntry, error) {
    query:="WHERE pasture_entries.animal_id = $1"
    return r.Base.FindListByQuery(query, animalId)
}

func (r *PastureEntryRepository) Add(newEntry entity.PastureEntry) (*entity.PastureEntry, error) {
    return r.Base.Add(newEntry)
}

func (r *PastureEntryRepository) Save(entry *entity.PastureEntry) error {
    return r.Base.Save(entry)
}

func (r* PastureEntryRepository) Delete(id string) error {
    return r.Base.Delete(id)
}
