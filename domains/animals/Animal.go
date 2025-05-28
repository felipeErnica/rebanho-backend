package animals

import "time"

type Animal struct {
	Id                   string     `json:"id" db:"id"`
	ChipId               *string    `json:"chipId" db:"chip_id"`
	Name                 *string    `json:"name" db:"name"`
	Number               *string    `json:"number" db:"number"`
	Color                *string    `json:"color" db:"color"`
	WeightBirth          float64    `json:"weightBirth" db:"weight_birth"`
	AnimalOrder          int        `json:"animalOrder" db:"animal_order"`
	Sex                  string     `json:"sex" db:"sex"`
	WeaningDate          *time.Time `json:"weaningDate" db:"weaning_date"`
	FatherId             *string    `json:"fatherId" db:"father_id"`
	FatherName           *string    `json:"fatherName" db:"father_name"`
	FatherNumber         *string    `json:"fatherNumber" db:"father_number"`
	MotherId             *string    `json:"motherId" db:"mother_id"`
	MotherName           *string    `json:"motherName" db:"mother_name"`
	MotherNumber         *string    `json:"motherNumber" db:"mother_number"`
	BirthDate            *time.Time `json:"birthDate" db:"birth_date"`
	DeathDate            *time.Time `json:"deathDate" db:"death_date"`
	PastureId            *string    `json:"pastureId" db:"pasture_id"`
	PastureName          *string    `json:"pastureName" db:"pasture_name"`
	Status               string     `json:"status" db:"status"`
	Isr                  float64    `json:"isr" db:"isr"`
	AverageProd          float64    `json:"averageProd" db:"average_prod"`
	AverageBirthInterval float64    `json:"averageBirthInterval" db:"average_birth_interval"`
	AveragePeak          float64    `json:"averagePeak" db:"average_peak"`
	ChildrenQuantity     int        `json:"childrenQuantity" db:"children_quantity"`
	Observation          *string    `json:"observation" db:"observation"`
	IsDna                bool       `json:"isDna" db:"is_dna"`
	IsGenotipagem        bool       `json:"isGenotipagem" db:"is_genotipagem"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt            *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId               string     `json:"userId" db:"user_id"`
}

type AnimalSave struct {
    Id                   string     `json:"id" db:"id"`
    ChipId               string     `json:"chipId" db:"chip_id"`
    Name                 *string    `json:"name" db:"name"`
    Color                *string    `json:"color" db:"color"`
    WeightBirth          *float64   `json:"weightBirth" db:"weight_birth"`
    IdentificationNumber *string    `json:"number" db:"number"`
    AnimalOrder          int        `json:"animalOrder" db:"animal_order"`
    Sex                  *string    `json:"sex" db:"sex"`
    WeaningDate          *time.Time `json:"weaningDate" db:"weaning_date"`
    Father               string     `json:"fatherId" db:"father_id"`
    Mother               string     `json:"motherId" db:"mother_id"`
    BirthDate            *time.Time `json:"birthDate" db:"birth_date"`
    DeathDate            *time.Time `json:"deathDate" db:"death_date"`
    PastureId            string     `json:"pasture" db:"pasture_id"`
    Status               string     `json:"status" db:"status"`
    Isr                  float64    `json:"isr" db:"isr"`
    AverageProd          float64    `json:"averageProd" db:"average_prod"`
    AverageBirthInterval float64    `json:"averageBirthInterval" db:"average_birth_interval"`
    AveragePeak          *float64   `json:"averagePeak" db:"average_peak"`
    ChildrenQuantity     int        `json:"childrenQuantity" db:"children_quantity"`
    Observation          *string    `json:"observation" db:"observation"`
    IsDna                bool       `json:"isDna" db:"is_dna"`
    IsGenotipagem        bool       `json:"isGenotipagem" db:"is_genotipagem"`
    CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
    DeletedAt            *time.Time `json:"deletedAt" db:"deleted_at"`
    UserId               string     `json:"userId" db:"user_id"`
}

type AnimalFilter struct {
	IsFiltered              bool       `json:"isFiltered"`
	Name                    *string    `json:"name" db:"name"`
	Number                  *string    `json:"number" db:"number"`
	Sex                     *string    `json:"sex" db:"sex"`
	MinWeaningDate          *time.Time `json:"minWeaningDate" db:"weaning_date"`
	MaxWeaningDate          *time.Time `json:"maxWeaningDate" db:"weaning_date"`
	Fathers                 *[]string  `json:"fathers" db:"father_id"`
	Mothers                 *[]string  `json:"mothers" db:"mother_id"`
	MinBirthDate            *time.Time `json:"minBirthDate" db:"birth_date"`
	MaxBirthDate            *time.Time `json:"maxBirthDate" db:"birth_date"`
	MinDeathDate            *time.Time `json:"minDeathDate" db:"death_date"`
	MaxDeathDate            *time.Time `json:"maxDeathDate" db:"death_date"`
	Pastures                *[]string  `json:"pastures" db:"pasture_id"`
	Status                  *[]string  `json:"status" db:"status"`
	MinIsr                  *float64   `json:"minIsr" db:"isr"`
	MaxIsr                  *float64   `json:"maxIsr" db:"isr"`
	MinAverageProd          *float64   `json:"minAverageProd" db:"average_prod"`
	MaxAverageProd          *float64   `json:"maxAverageProd" db:"average_prod"`
	MinAverageBirthInterval *float64   `json:"minAverageBirthInterval" db:"average_birth_interval"`
	MaxAverageBirthInterval *float64   `json:"maxAverageBirthInterval" db:"average_birth_interval"`
	MinAveragePeak          *float64   `json:"minAveragePeak" db:"average_peak"`
	MaxAveragePeak          *float64   `json:"maxAveragePeak" db:"average_peak"`
	MinChildrenQuantity     *int       `json:"minChildrenQuantity" db:"children_quantity"`
	MaxChildrenQuantity     *int       `json:"maxChildrenQuantity" db:"children_quantity"`
}
