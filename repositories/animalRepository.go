package repositories

import (
	"github.com/felipeErnica/rebanho-backend/entity"
)

type AnimalRepository struct {}

func (r *AnimalRepository) GetAll() (*[]entity.Animal, error) {
    query:= "SELECT * FROM animals"
    sqlStatement, err:= selectQueryList(query)
    var animals []entity.Animal

    for sqlStatement.Next() {
        var animal entity.Animal
        
        err:= sqlStatement.Scan(&animal.Id, &animal.Name, &animal.IdentificationNumber)
        if err != nil {
            return nil, err
        }

        animals = append(animals, animal)
    }

    return &animals, err
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
