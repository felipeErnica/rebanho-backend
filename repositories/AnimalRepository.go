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
    Base BaseRepository
}

func (r *AnimalRepository) Init() {

    dateFields:= []string{
        "birth_date",
        "death_date",
        "created_at",
    }

    DeletedAt:= "animal.deleted_at"
    SelectPageBody:= `
        SELECT animal.id, animal.name, animal.identification_number, animal.birth_date, animal.death_date, 
        animal.status, animal.avarage_prod, animal.avarage_birth_interval, animal.max_peak,
        animal.children_quantity, animal.created_at, animal.isr,
        mother.id AS mother_id, mother.name, mother.identification_number,
        father.id AS father_id, father.name, father.identification_number,
        pasture.id AS pasture_id, pasture.name
        FROM animals as animal
        LEFT JOIN animals AS father ON father.id = animal.father_id
        LEFT JOIN animals AS mother ON mother.id = animal.mother_id
        LEFT JOIN pastures AS pasture ON pasture.id = animal.pasture_id

        `
    r.Base = BaseRepository{
        Repository: r,
        SelectPageBody: SelectPageBody,
        DeletedAtField: DeletedAt,
        DateFields: dateFields,
    }
}

func (r *AnimalRepository) returnSimpleQuery(criteria string) string {
    return fmt.Sprintf(
    `SELECT id, name, identification_number, father_id, mother_id, birth_date, death_date, 
        pasture_id, status, avarage_prod, avarage_birth_interval, max_peak,
        children_quantity, created_at, isr
    FROM animals
    %s`, 
    criteria)
}

func (r *AnimalRepository) animalUnique(sqlStatement *sql.Row) (animal entity.Animal, err error) {

    err = sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.FatherId, &animal.MotherId, &animal.BirthDate,
        &animal.DeathDate, &animal.PastureId, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
        &animal.ChildrenQuantity, &animal.CreatedAt, &animal.Isr)
    if err != nil {
        return 
    }

    return animal, nil
}

func (r *AnimalRepository) saveOrUpdateScan(query string, animal *entity.Animal) error {
    return execQuery(query, animal.Id, animal.Name, animal.IdentificationNumber, animal.FatherId, animal.MotherId,
        animal.BirthDate, animal.DeathDate, animal.PastureId, animal.Status, animal.AvarageProd, animal.AvarageBirthInterval, animal.MaxPeak,
        animal.ChildrenQuantity, animal.CreatedAt, animal.DeletedAt, animal.Isr)
}

func (r *AnimalRepository) buildPage(query string, sort string, args... any) (page entity.PageInterface, err error) {

    sqlStatement, err:= selectQueryList(query, args...)
    if err != nil {
        return
    }

    var animals []entity.AnimalComplete

    for sqlStatement.Next() {
        var animal entity.AnimalComplete

        err = sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.BirthDate,
            &animal.DeathDate, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
            &animal.ChildrenQuantity, &animal.CreatedAt, &animal.Isr,
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
    
    pageAnimal:= entity.PageAnimalComplete{
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
    default:
        return "animal.created_at", "animal.id"
    }

}
func (r *AnimalRepository) createNextCursor(sort string, arr []entity.AnimalComplete) (cursor string, err error) {

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
        key = fmt.Sprintf("%s,%s", lastEntry.BirthDate.Format(time.RFC3339Nano), lastEntry.Id) 
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
    default:
        key = fmt.Sprintf("%s,%s", lastEntry.CreatedAt, lastEntry.Id) 
    }

    cursor = base64.StdEncoding.EncodeToString([]byte(key))
    return cursor, err
}

func (r *AnimalRepository) GetPage(sort string, direction string, cursor string) (page entity.PageInterface, err error) {
    return r.Base.GetPage(sort, direction, cursor)
}

func (r *AnimalRepository) GetPageDelete(sort string, direction string, cursor string) (page entity.PageInterface, err error) {
    return r.Base.GetPageDelete(sort, direction, cursor)
}

func (r *AnimalRepository) GetById(id string) (*entity.Animal, error) {
    query:= r.returnSimpleQuery("WHERE animal.id = $1")
    sqlStatement:= selectQueryOne(query, id)
    
    animal, err:=r.animalUnique(sqlStatement)

    return &animal, err
}

func (r *AnimalRepository) GetByMotherId(motherId string) (*[]entity.Animal, error) {
    query:= r.returnSimpleQuery("WHERE mother_id = $1 ORDER BY birth_date ASC")
    sqlStatement, err:= selectQueryList(query, motherId)
    defer sqlStatement.Close()
    if err != nil {
        return nil, err
    }
    
    var animals []entity.Animal
    for sqlStatement.Next() {
        var animal entity.Animal
        err:= sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.FatherId, &animal.MotherId, &animal.BirthDate,
        &animal.DeathDate, &animal.PastureId, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
        &animal.ChildrenQuantity, &animal.CreatedAt, &animal.Isr)
        if err != nil {
            return nil, err
        }
        animals = append(animals, animal)
    }
    
    return &animals, err
}

func (r *AnimalRepository) GetByFatherId(fatherId string) (*[]entity.Animal, error) {
    query:= r.returnSimpleQuery("WHERE father_id = $1 ORDER BY birth_date ASC")
    sqlStatement, err:= selectQueryList(query, fatherId)
    defer sqlStatement.Close()
    if err != nil {
        return nil, err
    }
    
    var animals []entity.Animal
    for sqlStatement.Next() {
        var animal entity.Animal
        err:= sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.FatherId, &animal.MotherId, &animal.BirthDate,
        &animal.DeathDate, &animal.PastureId, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
        &animal.ChildrenQuantity, &animal.CreatedAt, &animal.Isr)
        if err != nil {
            return nil, err
        }
        animals = append(animals, animal)
    }
    
    return &animals, err
}

func (r *AnimalRepository) GetByPastureId(pastureId string) (page *entity.PageAnimal, err error) {
    query:= fmt.Sprintf(
        `SELECT id, name, identification_number, birth_date, status
        FROM animals
        WHERE pasture_id = $1 AND (identification_number, id) > ($2, $3)
        ORDER BY identification_number ASC, id ASC
        LIMIT %d`, PAGE_LIMIT)
    sqlStatement, err:= selectQueryList(query, pastureId)

    var animals []entity.Animal
    for sqlStatement.Next() {
        var animal entity.Animal
        err = sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.BirthDate, &animal.Status)
        if err != nil {
            return
        }
        animals = append(animals, animal)
    }

    if len(animals) == 0 {
        err = serverErrors.EmptyList()
        return
    }

    lastEntry:= animals[len(animals) - 1]
    key:= fmt.Sprintf("%s, %s", *lastEntry.IdentificationNumber, *lastEntry.Id)
    nextCursor:= base64.RawStdEncoding.EncodeToString([]byte(key))
    page = &entity.PageAnimal{
        HasNextPage: PAGE_LIMIT != len(animals),
        NextCursor: nextCursor,
        List: &animals,
    }

    return page, err
}

func (r *AnimalRepository) Add(animal *entity.CreateAnimal) (*entity.Animal, error) {
    query:= 
    `INSERT INTO animals (id, name, identification_number, father_id, mother_id, 
        birth_date, death_date, pasture_id, status, avarage_prod, avarage_birth_interval, max_peak,
        children_quantity, created_at, deleted_at, isr) 
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`
    newAnimal := entity.NewAnimal(animal)
    err:= r.saveOrUpdateScan(query, newAnimal)
    return newAnimal, err
}

func (r *AnimalRepository) Save(animal *entity.Animal) error {
    query:= 
    `UPDATE animals 
    SET name = $1, identification_number = $2, father_id = $3, mother_id = $4, birth_date = $5, death_date = $6, 
        pasture_id = $7, status = $8, avarage_prod = $9, avarage_birth_interval = $10, max_peak = $11, children_quantity = $12,
        created_at = $13, deleted_at = $14, isr = $16
    WHERE id = $1`
    err:= r.saveOrUpdateScan(query, animal)
    return err
}
