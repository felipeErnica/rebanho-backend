package repositories

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type AnimalRepository struct {}

func (r *AnimalRepository) returnListQuery(criteria string) string {
    return fmt.Sprintf(
    `SELECT animal.id, animal.name, animal.identification_number, animal.birth_date, animal.death_date, 
        animal.status, animal.avarage_prod, animal.avarage_birth_interval, animal.max_peak,
        animal.children_quantity, animal.created_at, animal.isr,
        mother.id AS mother_id, mother.name, mother.identification_number,
        father.id AS father_id, father.name, father.identification_number,
        pasture.id AS pasture_id, pasture.name
    FROM animals AS animal
    LEFT JOIN animals AS father ON father.id = animal.father_id
    LEFT JOIN animals AS mother ON mother.id = animal.mother_id
    LEFT JOIN pastures AS pasture ON pasture.id = animal.pasture_id
    %s
    LIMIT %d`, 
    criteria, PAGE_LIMIT)
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

func (r *AnimalRepository) animalArray(sqlStatement *sql.Rows) ([]entity.AnimalResponse, error) {
    var animals []entity.AnimalResponse

    for sqlStatement.Next() {
        var animal entity.AnimalResponse

        err:= sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.BirthDate,
            &animal.DeathDate, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
            &animal.ChildrenQuantity, &animal.CreatedAt, &animal.Isr,
            &animal.Mother.Id, &animal.Mother.Name, &animal.Mother.IdentificationNumber, 
            &animal.Father.Id, &animal.Father.Name, &animal.Father.IdentificationNumber,
            &animal.Pasture.Id, &animal.Pasture.Name)
        if err != nil {
            return nil, err
        }

        animals = append(animals, animal)
    }
    
    return animals, nil
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

func (r *AnimalRepository) GetAll() (*[]entity.AnimalResponse, error) {
    query:= 
    `SELECT animal.id, animal.name, animal.identification_number, animal.birth_date, animal.death_date, 
        animal.status, animal.avarage_prod, animal.avarage_birth_interval, animal.max_peak,
        animal.children_quantity, animal.created_at, animal.isr,
        mother.id AS mother_id, mother.name, mother.identification_number,
        father.id AS father_id, father.name, father.identification_number,
        pasture.id AS pasture_id, pasture.name
    FROM animals AS animal
    LEFT JOIN animals AS father ON father.id = animal.father_id
    LEFT JOIN animals AS mother ON mother.id = animal.mother_id
    LEFT JOIN pastures AS pasture ON pasture.id = animal.pasture_id`

    sqlStatement, err:= selectQueryList(query)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    return &animals, err
}

func createCriteriaFirstPage(sort string, direction string) string {
    
    var criteria string;
    direction = strings.ToUpper(direction)

    switch (sort) {
    case "name": 
        criteria = fmt.Sprintf("ORDER BY animal.name %[1]s, animal.id %[1]s", direction)
    case "identification_number": 
        criteria = fmt.Sprintf("ORDER BY animal.identification_number %[1]s, animal.id %[1]s", direction)
    case "birth_date":
        criteria = fmt.Sprintf("ORDER BY animal.birth_date %[1]s, animal.id %[1]s", direction)
    case "death_date":
        criteria = fmt.Sprintf("ORDER BY animal.death_date %[1]s, animal.id %[1]s", direction)
    case "avarage_prod":
        criteria = fmt.Sprintf("ORDER BY animal.avarage_prod %[1]s, animal.id %[1]s", direction)
    case "avarage_birth_interval":
        criteria = fmt.Sprintf("ORDER BY animal.avarage_birth_interval %[1]s, animal.id %[1]s", direction)
    case "max_peak":
        criteria = fmt.Sprintf("ORDER BY animal.max_peak %[1]s, animal.id %[1]s", direction)
    case "children_quantity":
        criteria = fmt.Sprintf("ORDER BY animal.children_quantity %[1]s, animal.id %[1]s", direction)
    case "isr":
        criteria = fmt.Sprintf("ORDER BY animal.isr %[1]s, animal.id %[1]s", direction)
    default:
        criteria = "ORDER BY animal.created_at, animal.id"
    }

    return criteria;
}

func (r *AnimalRepository) createNextCursor(sort string, arr []entity.AnimalResponse) (cursor string, err error) {
    
    if (len(arr) == 0) {
        err = errors.New("A matriz está vazia")
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
        key = fmt.Sprintf("%s,%s", lastEntry.BirthDate, lastEntry.Id) 
    case "death_date":
        key = fmt.Sprintf("%s,%s", lastEntry.DeathDate, lastEntry.Id) 
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

func (r *AnimalRepository) GetFirstPage(sort string, direction string) (*entity.PageAnimal, error) {
    criteria:=createCriteriaFirstPage(sort, direction)
    query:= r.returnListQuery(criteria)
    sqlStatement, err:= selectQueryList(query)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    if err != nil {
        return nil, err
    }

    cursor, err:= r.createNextCursor(sort, animals)
    if err != nil {
        return nil, err
    }
    
    page:= &entity.PageAnimal{
        HasNextPage: len(animals) == PAGE_LIMIT,
        NextCursor: cursor,
        List: &animals,
    }
    
    return page, err
}

func (r *AnimalRepository) createCriteriaPages(sort string, direction string) string {

    var signal string;
    switch (direction) {
    case "asc": 
        signal = ">"
    case "desc":
        signal = "<"
    }

    direction = strings.ToUpper(direction)
    var criteria string;

    switch (sort) {
    case "name": 
        criteria = fmt.Sprintf(`
            WHERE (animal.name, animal.id) %s ($1, $2)
            ORDER BY animal.name %[2]s, animal.id %[2]s`, signal, direction)
    case "identification_number": 
        criteria = fmt.Sprintf(`
            ORDER BY animal.identification_number %[1]s, animal.id %[1]s
            WHERE (animal.identification_number, animal.id) %s ($1, $2)`, direction, signal)
    case "birth_date":
        criteria = fmt.Sprintf(`
            ORDER BY animal.birth_date %[1]s, animal.id %[1]s
            WHERE (animal.birth_date, animal.id)`, direction, signal)
    case "death_date":
        criteria = fmt.Sprintf(`
            ORDER BY animal.death_date %[1]s, animal.id %[1]s
            WHERE (animal.death_date, animal.id) %s ($1, $2)`, direction, signal)
    case "avarage_prod":
        criteria = fmt.Sprintf(`
            ORDER BY animal.avarage_prod %[1]s, animal.id %[1]s
            WHERE (animal.avarage_prod, animal.id) %s ($1, $2)`, direction, signal)
    case "avarage_birth_interval":
        criteria = fmt.Sprintf(`
            ORDER BY animal.avarage_birth_interval %[1]s, animal.id %[1]s
            WHERE (animal.avarage_birth_interval, animal.id) %s ($1, $2)`, direction, signal)
    case "max_peak":
        criteria = fmt.Sprintf(`
            ORDER BY animal.max_peak %[1]s, animal.id %[1]s
            WHERE (animal.max_peak, animal.id) %s ($1, $2)`, direction, signal)
    case "children_quantity":
        criteria = fmt.Sprintf(`
            ORDER BY animal.children_quantity %[1]s, animal.id %[1]s
            WHERE (animal.children_quantity, animal.id) %s ($1, $2)`, direction, signal)
    case "isr":
        criteria = fmt.Sprintf(`
            ORDER BY animal.isr %[1]s, animal.id %[1]s
            WHERE (animal.isr, animal.id) %s ($1, $2)`, direction, signal)
    default:
        criteria = `
            ORDER BY animal.created_at, animal.id
            WHERE (animal.created_at, animal.id) > ($1, $2)`
    }

    return criteria;
}

func (r *AnimalRepository) GetNextPage(cursor string, sort string, direction string) (*entity.PageAnimal, error) {
    first, second, err:= decodeCursor(cursor)
    if err != nil {
        return nil, err
    }

    criteria:= r.createCriteriaPages(sort, direction)
    query:= r.returnListQuery(criteria)
    sqlStatement, err:= selectQueryList(query, first, second)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    if err != nil {
        return nil, err
    }
    
    nextCursor, err:= r.createNextCursor(sort, animals)
    if err != nil {
        return nil, err
    }

    page:= &entity.PageAnimal{
        HasNextPage: len(animals) == PAGE_LIMIT,
        NextCursor: nextCursor,
        List: &animals,
    }

    return page, err
}

func (r *AnimalRepository) GetById(id string) (*entity.Animal, error) {
    query:= r.returnSimpleQuery("WHERE animal.id = $1")
    sqlStatement:= selectQueryOne(query, id)
    
    animal, err:=r.animalUnique(sqlStatement)

    return &animal, err
}

func (r *AnimalRepository) saveOrUpdateScan(query string, animal *entity.Animal) error {
    return execQuery(query, animal.Id, animal.Name, animal.IdentificationNumber, animal.FatherId, animal.MotherId,
        animal.BirthDate, animal.DeathDate, animal.PastureId, animal.Status, animal.AvarageProd, animal.AvarageBirthInterval, animal.MaxPeak,
        animal.ChildrenQuantity, animal.CreatedAt, animal.DeletedAt, animal.Isr)
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
