package repositories

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
)

type AnimalRepository struct {
    Base      BaseRepository[entity.Animal]
}

func (r *AnimalRepository) Init() {

    dateFields:= []string{
        "birth_date",
        "death_date",
        "created_at",
        "deleted_at",
    }

    deletedAt:= "animal.deleted_at"
    selectPageBody:= `
        SELECT animal.id, animal.name, animal.identification_number, animal.birth_date, animal.death_date, 
            animal.status, animal.avarage_prod, animal.avarage_birth_interval, animal.max_peak,
            animal.children_quantity, animal.created_at, animal.deleted_at, animal.isr,
            mother.id AS mother_id, mother.name, mother.identification_number,
            father.id AS father_id, father.name, father.identification_number,
            pasture.id AS pasture_id, pasture.name
        FROM animals as animal
            LEFT JOIN animals AS father ON father.id = animal.father_id
            LEFT JOIN animals AS mother ON mother.id = animal.mother_id
            LEFT JOIN pastures AS pasture ON pasture.id = animal.pasture_id
    `
    simpleQueryBody:=`
        SELECT id, name, identification_number, father_id, mother_id, birth_date, death_date, 
            pasture_id, status, avarage_prod, avarage_birth_interval, max_peak,
            children_quantity, created_at, deleted_at, isr
        FROM animals
    `
    insertQuery:=`
        INSERT INTO animals (id, name, identification_number, father_id, mother_id, 
            birth_date, death_date, pasture_id, status, avarage_prod, avarage_birth_interval, max_peak,
            children_quantity, created_at, deleted_at, isr) 
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
    `

    updateQuery:=`
        UPDATE animals 
        SET name = $2, identification_number = $3, father_id = $4, mother_id = $5, birth_date = $6, death_date = $7, 
            pasture_id = $8, status = $9, avarage_prod = $10, avarage_birth_interval = $11, max_peak = $12, children_quantity = $13,
            created_at = $14, deleted_at = $15, isr = $16
        WHERE id = $1
    `

    deleteQuery:="UPDATE animals SET deleted_at = $1 WHERE id = $2"

    r.Base = BaseRepository[entity.Animal]{
        Repository: r,
        DeleteQuery: deleteQuery,
        SimpleQueryBody: simpleQueryBody,
        SelectPageBody: selectPageBody,
        InsertQuery: insertQuery,
        UpdateQuery: updateQuery,
        DeletedAtField: deletedAt,
        DateFields: dateFields,
    }
}

func (r *AnimalRepository) setNewId(model *entity.Animal, id string) {
    model.Id = id
}

func (r *AnimalRepository) buildListEntity(sqlRows *sql.Rows) (list *[]entity.Animal, err error) {
    var animals []entity.Animal
    for sqlRows.Next() {
        var animal entity.Animal
        err:= sqlRows.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.Father.Id, &animal.Mother.Id, &animal.BirthDate,
        &animal.DeathDate, &animal.Pasture.Id, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
        &animal.ChildrenQuantity, &animal.CreatedAt, &animal.DeletedAt, &animal.Isr)
        if err != nil {
            return nil, err
        }
        animals = append(animals, animal)
    }
    sqlRows.Close()
    return &animals, err
}

func (r *AnimalRepository) buildEntity(sqlStatement *sql.Row) (model *entity.Animal, err error) {
    var animal entity.Animal
    err = sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.Father.Id, &animal.Mother.Id, &animal.BirthDate,
        &animal.DeathDate, &animal.Pasture.Id, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
        &animal.ChildrenQuantity, &animal.CreatedAt, &animal.DeletedAt, &animal.Isr)
    if err != nil {
        return 
    }
    return &animal, nil
}

func (r *AnimalRepository) saveOrUpdateScan(query string, animal *entity.Animal) error {
    return execQuery(query, animal.Id, animal.Name, animal.IdentificationNumber, animal.Father.Id, animal.Mother.Id,
        animal.BirthDate, animal.DeathDate, animal.Pasture.Id, animal.Status, animal.AvarageProd, animal.AvarageBirthInterval, animal.MaxPeak,
        animal.ChildrenQuantity, animal.CreatedAt, animal.DeletedAt, animal.Isr)
}

func (r *AnimalRepository) buildPage(query string, sort string, args... any) (page *entity.Page[entity.Animal], err error) {
    sqlStatement, err:= selectQueryList(query, args...)
    if err != nil {
        return
    }

    var animals []entity.Animal
    for sqlStatement.Next() {
        var animal entity.Animal

        err = sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.BirthDate,
            &animal.DeathDate, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
            &animal.ChildrenQuantity, &animal.CreatedAt, &animal.DeletedAt, &animal.Isr,
            &animal.Mother.Id, &animal.Mother.Name, &animal.Mother.IdentificationNumber, 
            &animal.Father.Id, &animal.Father.Name, &animal.Father.IdentificationNumber,
            &animal.Pasture.Id, &animal.Pasture.Name)
        if err != nil {
            return 
        }

        animals = append(animals, animal)
    }
    sqlStatement.Close()

    nextCursor, err:= r.createNextCursor(sort, animals)
    if err !=  nil {
        return
    }
    
    pageAnimal:= entity.Page[entity.Animal]{
        HasNextPage: len(animals) == PAGE_LIMIT,
        NextCursor: nextCursor,
        List: &animals,
    }

    return &pageAnimal, err
}

func (r *AnimalRepository) getFields(sort string) (firstField string, secondField string) {

    switch (sort) {
    case "name": 
        return "animal.name", "animal.id"
    case "identification_number": 
        return "animal.animal_order", "animal.id"
    case "birth_date":
        return "animal.birth_date", "animal.id"
    case "death_date":
        return "animal.death_date", "animal.id"
    case "avarage_prod":
        return "animal.avarage_prod", "animal.id"
    case "avarage_birth_interval":
        return "animal.avarage_birth_interval", "animal.id"
    case "max_peak":
        return "animal.max_peak", "animal.id"
    case "children_quantity":
        return "animal.children_quantity", "animal.id"
    case "isr":
        return "animal.isr", "animal.id"
    case "deleted_at":
        return "animal.deleted_at", "animal.id"
    default:
        return "animal.created_at", "animal.id"
    }

}
func (r *AnimalRepository) createNextCursor(sort string, arr []entity.Animal) (cursor string, err error) {

    if (len(arr) == 0) {
        err = serverErrors.EmptyList()
        return 
    }

    lastEntry:= arr[len(arr) - 1]

    var key string;
    switch (sort) {
    case "name": 
        key = fmt.Sprintf("%s,%s", *lastEntry.Name, lastEntry.Id) 
    case "identification_number": 
        key = fmt.Sprintf("%s,%s",*lastEntry.IdentificationNumber, lastEntry.Id) 
    case "birth_date":
        key = fmt.Sprintf("%s,%s", "null", lastEntry.Id) 
        if lastEntry.BirthDate != nil {
            key = fmt.Sprintf("%s,%s", lastEntry.BirthDate.Format(time.RFC3339Nano), lastEntry.Id) 
        }
    case "death_date":
        key = fmt.Sprintf("%s,%s", "null", lastEntry.Id) 
        if lastEntry.DeathDate != nil {
            key = fmt.Sprintf("%s,%s", lastEntry.DeathDate.Format(time.RFC3339Nano), lastEntry.Id) 
        }
    case "avarage_prod":
        key = fmt.Sprintf("%f,%s", lastEntry.AvarageProd, lastEntry.Id) 
    case "avarage_birth_interval":
        key = fmt.Sprintf("%f,%s", lastEntry.AvarageBirthInterval, lastEntry.Id) 
    case "max_peak":
        key = fmt.Sprintf("%f,%s", lastEntry.MaxPeak, lastEntry.Id) 
    case "children_quantity":
        key = fmt.Sprintf("%d,%s", lastEntry.ChildrenQuantity, lastEntry.Id) 
    case "isr":
        key = fmt.Sprintf("%f,%s", lastEntry.Isr, lastEntry.Id) 
    case "deleted_at":
        key = fmt.Sprintf("%s,%s", "null", lastEntry.Id) 
        if lastEntry.DeletedAt != nil {
            key = fmt.Sprintf("%s,%s", lastEntry.DeletedAt.Format(time.RFC3339Nano), lastEntry.Id) 
        }
    default:
        key = fmt.Sprintf("%s,%s", lastEntry.CreatedAt.Format(time.RFC3339Nano), lastEntry.Id) 
    }

    cursor = base64.StdEncoding.EncodeToString([]byte(key))
    return cursor, err
}

func (r *AnimalRepository) FindPage(sort string, direction string, cursor string) (page *entity.Page[entity.Animal], err error) {
    return r.Base.FindPage(sort, direction, cursor)
}

func (r *AnimalRepository) FindDeletedPage(sort string, direction string, 
    cursor string) (page *entity.Page[entity.Animal], err error) {
    return r.Base.FindDeletedPage(sort, direction, cursor)
}

func (r *AnimalRepository) FindById(id string) (*entity.Animal, error) {
    return r.Base.FindById(id)
}

func (r *AnimalRepository) FindByMotherId(motherId string) (*[]entity.Animal, error) {
    return r.Base.FindListByQuery("WHERE mother_id = $1 ORDER BY birth_date ASC", motherId)
}

func (r *AnimalRepository) FindByFatherId(fatherId string) (*[]entity.Animal, error) {
    return r.Base.FindListByQuery("WHERE father_id = $1 ORDER BY birth_date ASC", fatherId)
}

func (r *AnimalRepository) FindByPastureId(sort string, direction string, 
    cursor string, pastureId string) (page *entity.Page[entity.Animal], err error) {
    return r.Base.FindPageCondional(sort, direction, cursor, "animal.pasture_id", pastureId)
}

func (r *AnimalRepository) Add(create entity.Animal) (*entity.Animal, error) {
    return r.Base.Add(create)
}

func (r *AnimalRepository) Save(animal *entity.Animal) error {
    return r.Base.Save(animal)
}

func (r *AnimalRepository) Delete(id string) error {
    return r.Base.Delete(id)
}

