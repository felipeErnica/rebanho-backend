package repositories

import (
	"database/sql"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type AnimalRepository struct {
	Db *sql.DB
}

func (r *AnimalRepository) Add(animal *entity.CreateAnimal) (*entity.Animal, error) {
    newAnimal := entity.NewAnimal(animal)
    _, err := r.Db.Exec("INSERT INTO animals(id, name, number) VALUES(?,?,?)", newAnimal.Id, newAnimal.Name, newAnimal.Number)
    return newAnimal, err
}
