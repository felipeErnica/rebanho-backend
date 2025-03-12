package entity

import (
	"time"

	"github.com/google/uuid"
)

type Animal struct {
    Id                   *string         `json:"id"`
    Name                 *string         `json:"name"`
    IdentificationNumber *string         `json:"identification_number"`
    RingOrder            int             `json:"ring_order"`  
    Sex                  string          `json:"sex"`
    FatherId             *string         `json:"father_id"`
    MotherId             *string         `json:"mother_id"`
    BirthDate            *time.Time      `json:"birth_date"`
    DeathDate            *time.Time      `json:"death_date"`
    PastureId            *string         `json:"pasture_id"`
    Status               *string         `json:"status"`
    Isr                  float32         `json:"isr"`
    AvarageProd          float32         `json:"avarage_prod"`
    AvarageBirthInterval float32         `json:"avarage_birth_interval"`
    MaxPeak              float32         `json:"max_peak"`
    ChildrenQuantity     int             `json:"children_quantity"`
    CreatedAt            time.Time      `json:"created_at"`
    DeletedAt            *time.Time      `json:"deleted_at"`
}

func NewAnimal(create *CreateAnimal) *Animal {
    id:=uuid.New().String()
    return &Animal{
        Id: &id,
        Name: &create.Name,
        IdentificationNumber: &create.IdentificationNumber,
        FatherId: &create.FatherId,
        MotherId: &create.MotherId,
        BirthDate: &create.BirthDate,
        DeathDate: &create.DeathDate,
        PastureId: &create.PastureId,
        Status: &create.Status,
        Isr: create.Isr,
        AvarageProd: create.AvarageProd,
        AvarageBirthInterval: create.AvarageBirthInterval,
        MaxPeak: create.MaxPeak,
        ChildrenQuantity: create.ChildrenQuantity,
        CreatedAt: time.Now(),
    }
}

type CreateAnimal struct {
    Name                 string         `json:"name"`
    IdentificationNumber string         `json:"identification_number"`
    RingOrder            int            `json:"ring_order"`  
    Sex                  string         `json:"sex"`
    FatherId             string         `json:"father_id"`
    MotherId             string         `json:"mother_id"`
    BirthDate            time.Time      `json:"birth_date"`
    DeathDate            time.Time      `json:"death_date"`
    PastureId            string         `json:"pasture_id"`
    Status               string         `json:"status"`
    Isr                  float32        `json:"isr"`
    AvarageProd          float32        `json:"avarage_prod"`
    AvarageBirthInterval float32        `json:"avarage_birth_interval"`
    MaxPeak              float32        `json:"max_peak"`
    ChildrenQuantity     int            `json:"children_quantity"`
    CreatedAt            time.Time      `json:"created_at"`
    DeletedAt            time.Time      `json:"deleted_at"`
}

type AnimalComplete struct {
	Id                   string          `json:"id"`
	Name                 *string         `json:"name"`
	IdentificationNumber *string         `json:"identification_number"`
    RingOrder            int             `json:"ring_order"`  
    Sex                  string          `json:"sex"`
	Father               AnimalShort     `json:"father"`
	Mother               AnimalShort     `json:"mother"`
	BirthDate            *time.Time      `json:"birth_date"`
	DeathDate            *time.Time      `json:"death_date"`
	Pasture              PastureShort    `json:"pasture"`
	Status               *string         `json:"status"`
	Isr                  float32         `json:"isr"`
	AvarageProd          float32         `json:"avarage_prod"`
	AvarageBirthInterval float32         `json:"avarage_birth_interval"`
	MaxPeak              float32         `json:"max_peak"`
	ChildrenQuantity     int             `json:"children_quantity"`
    CreatedAt            *time.Time      `json:"created_at"`
}

type AnimalShort struct {
    Id                   string
    IdentificationNumber string
    Name                 string
}

type PageAnimalComplete struct {
    NextCursor      string
    HasNextPage     bool
    List            *[]AnimalComplete
}

func (p *PageAnimalComplete) GetPage() PageInterface {
    return p
}

type PageAnimal struct {
    NextCursor      string
    HasNextPage     bool
    List            *[]Animal
}

