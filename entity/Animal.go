package entity

import "time"

type Animal struct {
	Id                   string       `json:"id"`
	ChipId               string       `json:"chipId"`
	Name                 *string      `json:"name"`
	IdentificationNumber *string      `json:"ringNumber"`
	Color                *string      `json:"color"`
	WeightBirth          *float32     `json:"weightBirth"`
	AnimalOrder          int          `json:"animalOrder"`
	Sex                  *string      `json:"sex"`
	WeaningDate          *time.Time   `json:"weaningDate"`
	Father               AnimalDto    `json:"father"`
	Mother               AnimalDto    `json:"mother"`
	BirthDate            *time.Time   `json:"birthDate"`
	DeathDate            *time.Time   `json:"deathDate"`
	Pasture              PastureShort `json:"pasture"`
	Status               string       `json:"status"`
	Isr                  float32      `json:"isr"`
	AverageProd          float32      `json:"averageProd"`
	AverageBirthInterval float32      `json:"averageBirthInterval"`
	AveragePeak          *float32     `json:"averagePeak"`
	ChildrenQuantity     int          `json:"childrenQuantity"`
	Observation          *string      `json:"observation"`
	IsDna                bool         `json:"isDna"`
	IsGenotipagem        bool         `json:"isGenotipagem"`
	CreatedAt            time.Time    `json:"createdAt"`
	DeletedAt            *time.Time   `json:"deletedAt"`
	UserId               string       `json:"userId"`
}

type AnimalDto struct {
	Id                   string     `json:"id"`
	ChipId               string     `json:"chipId"`
	Name                 *string    `json:"name"`
	Color                *string    `json:"color"`
	WeightBirth          *float32   `json:"weightBirth"`
	IdentificationNumber *string    `json:"ringNumber"`
	AnimalOrder          int        `json:"animalOrder"`
	Sex                  *string    `json:"sex"`
	WeaningDate          *time.Time `json:"weaningDate"`
	Father               string     `json:"fatherId"`
	Mother               string     `json:"motherId"`
	BirthDate            *time.Time `json:"birthDate"`
	DeathDate            *time.Time `json:"deathDate"`
	PastureId            string     `json:"pasture"`
	Status               string     `json:"status"`
	Isr                  float32    `json:"isr"`
	AverageProd          float32    `json:"averageProd"`
	AverageBirthInterval float32    `json:"averageBirthInterval"`
	AveragePeak          *float32   `json:"averagePeak"`
	ChildrenQuantity     int        `json:"childrenQuantity"`
	Observation          *string    `json:"observation"`
	IsDna                bool       `json:"isDna"`
	IsGenotipagem        bool       `json:"isGenotipagem"`
	CreatedAt            time.Time  `json:"createdAt"`
	DeletedAt            *time.Time `json:"deletedAt"`
	UserId               string     `json:"userId"`
}
type AnimalFilter struct {
	IsFiltered              bool       `json:"isFiltered"`
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
