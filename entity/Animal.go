package entity

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/enums"
	"github.com/google/uuid"
)

type Animal struct {
    Id                   string         `json:"Id"`
    Name                 string         `json:"Name"`
    IdentificationNumber string         `json:"IdentificationNumber"`
    FatherId             string         `json:"FatherId"`
    MotherId             string         `json:"MotherId"`
    BirthDate            time.Time      `json:"BirthDate"`
    DeathDate            time.Time      `json:"DeathDate"`
    PastureId            string         `json:"PastureId"`
    Status               enums.Status   `json:"Status"`
    Isr                  float32        `json:"Isr"`
    AvarageProd          float32        `json:"AvarageProd"`
    AvarageBirthInterval float32        `json:"AvarageBirthInterval"`
    MaxPeak              float32        `json:"MaxPeak"`
    ChildrenQuantity     uint8          `json:"ChildrenQuantity"`
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
