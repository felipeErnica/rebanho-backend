package entity

import "github.com/google/uuid"

type Animal struct {
    Id string `json "id"`
    Name string `json "name"`
    Number string `json "number"`
}

func NewAnimal(create *CreateAnimal) *Animal {
    return &Animal{
        Id: uuid.New().String(),
        Name: create.Name,
        Number: create.Number,
    }
}

type CreateAnimal struct {
    Name string
    Number string
}
