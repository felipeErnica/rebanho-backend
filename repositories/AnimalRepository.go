package repositories

import (
	"database/sql"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type AnimalRepository struct {}

func (r *AnimalRepository) returnSimpleQuery(criteria string) string {
    return fmt.Sprintf(`SELECT id, name, identification_number, father_id, mother_id, 
    birth_date, death_date, pasture_id, status, avarage_prod, avarage_birth_interval, max_peak,
    children_quantity, created_at, deleted_at, isr
    FROM animals
    %s`, criteria)
}

func (r *AnimalRepository) returnFirstPageQuery(criteria string) string {
    query:=fmt.Sprintf(`SELECT * 
        FROM (%s) 
        ORDER BY created_at DESC, id DESC
        LIMIT %d`, r.returnSimpleQuery(criteria), PAGE_LIMIT)
    return query
}

func (r *AnimalRepository) returnPageQuery(criteria string) string {
    query:=fmt.Sprintf(`SELECT * 
        FROM (%s) 
        WHERE (created_at,id) < ($1, $2) 
        ORDER BY created_at DESC, id DESC
        LIMIT %d`, r.returnSimpleQuery(criteria), PAGE_LIMIT)
    return query
}

func (r *AnimalRepository) animalArray(sqlStatement *sql.Rows) ([]entity.Animal, error) {
    var animals []entity.Animal

    for sqlStatement.Next() {
        var animal entity.Animal

        err:= sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber, &animal.FatherId, &animal.MotherId, &animal.BirthDate,
            &animal.DeathDate, &animal.PastureId, &animal.Status, &animal.AvarageProd, &animal.AvarageBirthInterval, &animal.MaxPeak,
            &animal.ChildrenQuantity, &animal.CreatedAt, &animal.DeletedAt, &animal.Isr)
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
            &animal.ChildrenQuantity, &animal.CreatedAt, &animal.DeletedAt, &animal.Isr)
        if err != nil {
            return 
        }

    return animal, nil
}

func (r *AnimalRepository) createPage(arr []entity.Animal) *entity.PageAnimal {
    lastEntry:=arr[len(arr) - 1]
    return &entity.PageAnimal{
        List: &arr,
        HasNextPage: len(arr) == PAGE_LIMIT,
        NextCursor: encodeCursor(*lastEntry.CreatedAt, *lastEntry.Id),
    }
}

func (r *AnimalRepository) GetAll() (*[]entity.Animal, error) {
    query:=r.returnSimpleQuery("")
    sqlStatement, err:= selectQueryList(query)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    return &animals, err
}

func (r *AnimalRepository) GetFirstPage() (*entity.PageAnimal, error) {
    query:= r.returnFirstPageQuery("ORDER BY name")
    sqlStatement, err:= selectQueryList(query)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    page:= r.createPage(animals)
    return page, err
}

func (r *AnimalRepository) GetNextPage(cursor string) (*entity.PageAnimal, error) {
    createdAt, id, err:= decodeCursor(cursor)
    if err != nil {
        return nil, err
    }

    query:= r.returnPageQuery("ORDER BY name")
    sqlStatement, err:= selectQueryList(query, createdAt, id)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    page:= r.createPage(animals)
    return page, err
}

func (r *AnimalRepository) GetById(id string) (*entity.Animal, error) {
    query:= r.returnSimpleQuery("WHERE id = $1")
    sqlStatement:= selectQueryOne(query, id)
    
    animal, err:=r.animalUnique(sqlStatement)

    return &animal, err
}

func (r *AnimalRepository) Add(animal *entity.CreateAnimal) (*entity.Animal, error) {
    query:= "INSERT INTO animals(id, name, identification_number) VALUES($1, $2, $3)"
    newAnimal := entity.NewAnimal(animal)
    err := execQuery(query, newAnimal.Id, newAnimal.Name, newAnimal.IdentificationNumber)
    return newAnimal, err
}

func (r *AnimalRepository) Save(animal *entity.Animal) error {
    query:= "UPDATE animals SET name = $1, identification_number = $2 WHERE id = $3"
    err := execQuery(query, animal.Name, animal.IdentificationNumber, animal.Id)
    return err
}
