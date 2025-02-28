package repositories

import (
	"database/sql"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type AnimalRepository struct {}

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

func (r *AnimalRepository) createPage(arr []entity.Animal) *entity.PageAnimal {
    lastEntry:=arr[len(arr) - 1]
    return &entity.PageAnimal{
        List: &arr,
        HasNextPage: len(arr) != PAGE_LIMIT,
        NextCursor: encodeCursor(*lastEntry.CreatedAt, *lastEntry.Id),
    }
}

func (r *AnimalRepository) GetAll() (*[]entity.Animal, error) {
    query:= "SELECT * FROM animals"
    sqlStatement, err:= selectQueryList(query)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    return &animals, err
}

func (r *AnimalRepository) GetFirstPage() (*entity.PageAnimal, error) {
    query:= fmt.Sprintf(`SELECT * 
        FROM animals 
        ORDER BY name, created_at DESC, id DESC 
        LIMIT %d`, PAGE_LIMIT)
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

    query:= fmt.Sprintf(`SELECT * 
        FROM animals 
        WHERE (created_at,id) < (%s,%s)
        ORDER BY name, created_at DESC, id DESC 
        LIMIT %d`, createdAt, id, PAGE_LIMIT)
    sqlStatement, err:= selectQueryList(query)
    if err != nil {
        return nil, err
    }

    animals, err:= r.animalArray(sqlStatement)
    page:= r.createPage(animals)
    return page, err
}

func (r *AnimalRepository) GetById(id string) (*entity.Animal, error) {
    query:= "SELECT * FROM animals WHERE id = $1"
    sqlStatement:= selectQueryOne(query, id)

    var animal entity.Animal
    err:= sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber)
    if err != nil {
        return nil, err
    }

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
