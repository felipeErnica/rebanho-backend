package entity

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/enums"
	"github.com/google/uuid"
)

type Animal struct {
    Id                   string         `json:"id"`
    Name                 string         `json:"name"`
    IdentificationNumber string         `json:"identification_Number"`
    FatherId             string         `json:"father_id"`
    MotherId             string         `json:"mother_id"`
    BirthDate            time.Time      `json:"birth_date"`
    DeathDate            time.Time      `json:"death_date"`
    PastureId            string         `json:"pasture_id"`
    Status               enums.Status   `json:"status"`
    Isr                  float32        `json:"isr"`
    AvarageProd          float32        `json:"avarage_prod"`
    AvarageBirthInterval float32        `json:"avarage_birth_interval"`
    MaxPeak              float32        `json:"max_peak"`
    ChildrenQuantity     int            `json:"children_quantity"`
    CreatedAt            time.Time      `json:"created_at"`
    DeletedAt            time.Time      `json:"deleted_at"`
}

func NewAnimal(create *CreateAnimal) *Animal {
    return &Animal{
        Id: uuid.New().String(),
        Name: create.Name,
        IdentificationNumber: create.IdentificationNumber,
    }
}

type CreateAnimal struct {
    Name string `json:"name"`
    IdentificationNumber string `json:"identification_number"`
}
