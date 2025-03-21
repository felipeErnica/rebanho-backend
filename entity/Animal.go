package entity

import "time"

type Animal struct {
	Id                   string       `json:"id"`
	Name                 *string      `json:"name"`
	IdentificationNumber *string      `json:"identification_number"`
	AnimalOrder          int          `json:"animal_order"`
	Sex                  string       `json:"sex"`
	WeaningDate          time.Time    `json:"weaning_date"`
	Father               AnimalShort  `json:"father"`
	Mother               AnimalShort  `json:"mother"`
	BirthDate            *time.Time   `json:"birth_date"`
	DeathDate            *time.Time   `json:"death_date"`
	Pasture              PastureShort `json:"pasture"`
	Status               string       `json:"status"`
	Isr                  float32      `json:"isr"`
	AverageProd          float32      `json:"average_prod"`
	AverageBirthInterval float32      `json:"average_bitrh_interval"`
	MaxPeak              float32      `json:"max_peak"`
	ChildrenQuantity     int          `json:"children_quantity"`
	CreatedAt            time.Time    `json:"created_at"`
	DeletedAt            *time.Time   `json:"deleted_at"`
	UserId               string       `json:"user_id"`
}

type AnimalShort struct {
	Id                   *string `json:"id"`
	Name                 *string `json:"name"`
	IdentificationNumber *string `json:"identification_number"`
	AnimalOrder          *int    `json:"animal_order"`
}

type Cow struct {
	Id                   *string      `json:"id"`
	Name                 *string      `json:"name"`
	IdentificationNumber *string      `json:"identification_number"`
	Pasture              PastureShort `json:"pasture"`
	Status               *string      `json:"status"`
}

type Calf struct {
	Id                   *string     `json:"id"`
	Name                 *string     `json:"name"`
	IdentificationNumber *string     `json:"identification_number"`
	Sex                  *string     `json:"sex"`
	Father               AnimalShort `json:"father"`
	BirthDate            *time.Time  `json:"birth_date"`
}

type CalfShort struct {
	Id                   *string     `json:"id"`
	Sex                  *string     `json:"sex"`
	BirthDate            *time.Time  `json:"birth_date"`
}

type AnimalName struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}

type BullInsemintation struct {
	Id     *string    `json:"id"`
	Name   *string    `json:"name"`
	Mother AnimalName `json:"mother"`
	Father AnimalName `json:"father"`
}
