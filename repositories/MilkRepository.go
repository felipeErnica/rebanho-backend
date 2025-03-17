package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/google/uuid"
)

type MilkRepository struct{
    Base PageRepositoryImpl[entity.MilkEntry]
}

func (m *MilkRepository) Init() {
    
    dateFields:= []string{
        "entry_date",
        "created_at",
    }

    selectQuery:=`
        SELECT milk_entries.id, milk_entries.entry_date, milk_entries.milk_quantity, milk_entries.lactation_id
            animal.id as animal_id, animal.identification_number as animal_number, animal.ring_order as animal_order,
            animal.name as animal_name,
            pasture.id as pasture_id, pasture.name as pasture_name
        FROM milk_actives AS milk_entries
        LEFT JOIN animals AS animal ON animal.id = milk_entries.animal_id
        LEFT JOIN pastures AS pasture ON pasture.id = milk_entries.pasture_id
    `
    
    insertQuery:=`
        INSERT INTO milk_entries (id, entry_date, milk_quantity, animal_id, pasture_id, lactation_id, created_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

    updateQuery:=`
        UPDATE milk_entries 
        SET entry_date = $2, milk_quantity = $3, animal_id = $4, pasture_id = $5, lactation_id = $6, created_at = $7)
        WHERE id = $1
    `

    mainRepository:=&RepositoryImpl[entity.MilkEntry]{
        Repository: m,
        SelectQueryBody: selectQuery,
        UpdateQuery: updateQuery,
        InsertQuery: insertQuery,
        TableName: "milk_entries",
    }

    m.Base = PageRepositoryImpl[entity.MilkEntry] {
        Base: mainRepository,
        PageRepository: m,
        DateFields: dateFields,
    }
}

func (m *MilkRepository) getFields(sort string) (firstField string, secondField string) {
    switch (sort) {
    case "name": 
        return "animal.name", "milk_entries.id"
    case "identification_number": 
        return "animal.animal_order", "milk_entries.id"
    case "entry_date":
        return "milk_entries.entry_date", "milk_entries.id"
    default:
        return "milk_entries.created_at", "milk_entries.id"
    }
}

func (m *MilkRepository) createKey(sort string, lastEntry *entity.MilkEntry) (key string) {
    switch (sort) {
    case "name": 
        return fmt.Sprintf("%s,%s", *lastEntry.Animal.Name, lastEntry.Id)
    case "identification_number": 
        return fmt.Sprintf("%d,%s", *lastEntry.Animal.AnimalOrder, lastEntry.Id)
    case "entry_date":
        return fmt.Sprintf("%s,%s", lastEntry.EntryDate, lastEntry.Id)
    default:
        return fmt.Sprintf("%s,%s", lastEntry.CreatedAt, lastEntry.Id)
    }
}

func (m *MilkRepository) buildPage(query string, sort string, args... any) (page *entity.Page[entity.MilkEntry], err error) {
    rows, err:= selectQueryList(query, args...)
    if err != nil {
        return
    }

    var entries []entity.MilkEntry
    for rows.Next() {
        var milk entity.MilkEntry
        err = rows.Scan(&milk.Id, &milk.EntryDate, &milk.MilkQuantity, &milk.LactationId,
            &milk.Animal.Id, &milk.Animal.IdentificationNumber, &milk.Animal.AnimalOrder, &milk.Animal.Name,
            &milk.Pasture.Id, &milk.Pasture.Name)
        if err != nil {
            return 
        }

        entries = append(entries, milk)
    }
    rows.Close()

    nextCursor, err:= m.Base.CreateNextCursor(sort, entries)
    if err !=  nil {
        return
    }
    
    page = &entity.Page[entity.MilkEntry]{
        HasNextPage: len(entries) == PAGE_LIMIT,
        NextCursor: nextCursor,
        List: &entries,
    }

    return page, err
}

func (m *MilkRepository) setNewEntity(model *entity.MilkEntry) {
    model.Id = uuid.NewString()
    model.CreatedAt = time.Now()
}

func (m *MilkRepository) buildEntity(row *sql.Row) (model *entity.MilkEntry, err error) {
    var milk entity.MilkEntry
    err = row.Scan(&milk.Id, &milk.EntryDate, &milk.MilkQuantity, &milk.LactationId,
        &milk.Animal.Id, &milk.Animal.IdentificationNumber, &milk.Animal.AnimalOrder, &milk.Animal.Name,
        &milk.Pasture.Id, &milk.Pasture.Name)
    return &milk, err
}

func (m *MilkRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.MilkEntry, err error) {
    var entries []entity.MilkEntry
    for rows.Next() {
        var milk entity.MilkEntry
        err = rows.Scan(&milk.Id, &milk.EntryDate, &milk.MilkQuantity, &milk.LactationId,
            &milk.Animal.Id, &milk.Animal.IdentificationNumber, &milk.Animal.AnimalOrder, &milk.Animal.Name,
            &milk.Pasture.Id, &milk.Pasture.Name)
        if err != nil {
            return  
        }
        entries = append(entries, milk)
    }
    return &entries, err
}
	
func (m *MilkRepository) saveOrUpdateScan(query string, model *entity.MilkEntry) error {
    return execQuery(query, model.Id, model.EntryDate, model.MilkQuantity, 
        model.Animal.Id, model.Pasture.Id, model.LactationId, model.CreatedAt)
}

func (m *MilkRepository) FindPage(sort string, direction string, cursor string) (*entity.Page[entity.MilkEntry], error) {
    return m.Base.FindPage(sort, direction, cursor)
}

func (m *MilkRepository) FindByAnimal(animalId string) (*[]entity.MilkEntry, error) {
    query:= "WHERE milk_entries.animal_id = $1"
    return m.Base.FindListByQuery(query, animalId)
}

func (m *MilkRepository) FindByEntryDate(entryDate time.Time) (*[]entity.MilkEntry, error) {
    query:= "WHERE milk_entries.entryDate = $1"
    return m.Base.FindListByQuery(query, entryDate)
}

func (m *MilkRepository) Add(newMilk entity.MilkEntry) (*entity.MilkEntry, error) {
    return m.Base.Add(newMilk)
}

func (m *MilkRepository) Save(milk *entity.MilkEntry) error {
    return m.Base.Save(milk)
}

func (m *MilkRepository) Delete(id string) error {
    return m.Base.Delete(id)
}
