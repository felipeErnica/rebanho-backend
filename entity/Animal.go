package entity

import "time"

type Animal struct {
	Id                   string       `json:"id"`
	Name                 *string      `json:"name"`
	IdentificationNumber *string      `json:"identification_number"`
	AnimalOrder          int          `json:"animal_order"`
	Sex                  *string      `json:"sex"`
	WeaningDate          *time.Time   `json:"weaning_date"`
	Father               AnimalShort  `json:"father"`
	Mother               AnimalShort  `json:"mother"`
	BirthDate            *time.Time   `json:"birth_date"`
	DeathDate            *time.Time   `json:"death_date"`
	Pasture              PastureShort `json:"pasture"`
	Status               string       `json:"status"`
	Isr                  float32      `json:"isr"`
	AverageProd          float32      `json:"average_prod"`
	AverageBirthInterval float32      `json:"average_birth_interval"`
	AveragePeak          float32      `json:"average_peak"`
	ChildrenQuantity     int          `json:"children_quantity"`
	Observation          *string      `json:"observation"`
	CreatedAt            time.Time    `json:"created_at"`
	DeletedAt            *time.Time   `json:"deleted_at"`
	UserId               string       `json:"user_id"`
}

type AnimalFilter struct {
	IsFiltered              bool       `json:"is_filtered"`
	Name                    *string    `json:"name"`
	IdentificationNumber    *string    `json:"identification_number"`
	Sex                     *string    `json:"sex"`
	MinWeaningDate          *time.Time `json:"min_weaning_date"`
	MaxWeaningDate          *time.Time `json:"max_weaning_date"`
	Fathers                 *[]string  `json:"fathers"`
	Mothers                 *[]string  `json:"mothers"`
	MinBirthDate            *time.Time `json:"min_birth_date"`
	MaxBirthDate            *time.Time `json:"max_birth_date"`
	MinDeathDate            *time.Time `json:"min_death_date"`
	MaxDeathDate            *time.Time `json:"max_death_date"`
	Pastures                *[]string  `json:"pastures"`
	Status                  *[]string  `json:"status"`
	MinIsr                  *float32   `json:"min_isr"`
	MaxIsr                  *float32   `json:"max_isr"`
	MinAverageProd          *float32   `json:"min_average_prod"`
	MaxAverageProd          *float32   `json:"max_average_prod"`
	MinAverageBirthInterval *float32   `json:"min_average_birth_interval"`
	MaxAverageBirthInterval *float32   `json:"max_average_birth_interval"`
	MinAveragePeak          *float32   `json:"min_average_peak"`
	MaxAveragePeak          *float32   `json:"max_average_peak"`
	MinChildrenQuantity     *int       `json:"min_children_quantity"`
	MaxChildrenQuantity     *int       `json:"max_children_quantity"`
}

type AnimalMaxValues struct {
	MaxWeaningDate          *time.Time `json:"max_weaning_date"`
	MaxBirthDate            *time.Time `json:"max_birth_date"`
	MaxDeathDate            *time.Time `json:"max_death_date"`
	MaxIsr                  *float32   `json:"max_isr"`
	MaxAverageProd          *float32   `json:"max_average_prod"`
	MaxAverageBirthInterval *float32   `json:"max_average_birth_interval"`
	MaxAveragePeak          *float32   `json:"max_average_peak"`
	MaxChildrenQuantity     *int       `json:"max_children_quantity"`
}

type AnimalMinValues struct {
	MinWeaningDate          *time.Time `json:"min_weaning_date"`
	MinBirthDate            *time.Time `json:"min_birth_date"`
	MinDeathDate            *time.Time `json:"min_death_date"`
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
	Id        *string    `json:"id"`
	Sex       *string    `json:"sex"`
	BirthDate *time.Time `json:"birth_date"`
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
