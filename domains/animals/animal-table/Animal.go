package animalTable

import "time"

type Animal struct {
	Id                   string     `json:"id" db:"id"`
	OldId                int        `json:"oldId" db:"old_id"`
	Name                 *string    `json:"name" db:"name"`
	RingNumber           *string    `json:"ringNumber" db:"ring_number"`
	AnimalOrder          int        `json:"animalOrder" db:"animal_order"`
	WeightBirth          float64    `json:"weightBirth" db:"weight_birth"`
	Sex                  string     `json:"sex" db:"sex"`
	WeaningDate          *time.Time `json:"weaningDate" db:"weaning_date"`
	FatherId             *string    `json:"fatherId" db:"father_id"`
	FatherName           *string    `json:"fatherName" db:"father_name"`
	MotherId             *string    `json:"motherId" db:"mother_id"`
	MotherName           *string    `json:"motherName" db:"mother_name"`
	BirthDate            *time.Time `json:"birthDate" db:"birth_date"`
	DeathDate            *time.Time `json:"deathDate" db:"death_date"`
	PastureId            *string    `json:"pastureId" db:"pasture_id"`
	PastureName          *string    `json:"pastureName" db:"pasture_name"`
	FarmId               *string    `json:"farmId" db:"farm_id"`
	FarmName             *string    `json:"farmName" db:"farm_name"`
	AnimalType           string     `json:"animalType" db:"animal_type"`
	Isr                  *float64   `json:"isr" db:"isr"`
	AverageProd          *float64   `json:"averageProd" db:"average_prod"`
	AverageProdInterval  *float64   `json:"averageProdInterval" db:"average_prod_interval"`
	AverageBirthInterval *float64   `json:"averageBirthInterval" db:"average_birth_interval"`
	AveragePeak          *float64   `json:"averagePeak" db:"average_peak"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt            *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId               string     `json:"userId" db:"user_id"`
}

type AnimalSave struct {
	Id                   *string    `json:"id" db:"id"`
	Name                 *string    `json:"name" db:"name"`
	WeightBirth          *float64   `json:"weightBirth" db:"weight_birth"`
	RingNumber           *string    `json:"ringNumber" db:"ring_number"`
	Sex                  *string    `json:"sex" db:"sex"`
	WeaningDate          *time.Time `json:"weaningDate" db:"weaning_date"`
	FatherId             *string    `json:"fatherId" db:"father_id"`
	MotherId             *string    `json:"motherId" db:"mother_id"`
	BirthDate            *time.Time `json:"birthDate" db:"birth_date"`
	DeathDate            *time.Time `json:"deathDate" db:"death_date"`
	PastureId            *string    `json:"pasture" db:"pasture_id"`
	AnimalType           string     `json:"animalType" db:"animal_type"`
	Isr                  float64    `json:"isr" db:"isr"`
	AverageProd          float64    `json:"averageProd" db:"average_prod"`
	AverageProdInterval  float64    `json:"averageProdInterval" db:"average_prod_interval"`
	AverageBirthInterval float64    `json:"averageBirthInterval" db:"average_birth_interval"`
	AveragePeak          *float64   `json:"averagePeak" db:"average_peak"`
	DeletedAt            *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId               string     `json:"userId" db:"user_id"`
}

type AnimalFilter struct {
	IsFiltered              bool       `json:"isFiltered"`
	Name                    *string    `json:"name" db:"name"`
	Number                  *string    `json:"ringNumber" db:"ring_number"`
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
	Farms                   *[]string  `json:"farms" db:"farm_id" table:"pastures"`
	Types                   *[]string  `json:"types" db:"type"`
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
