package breeding

import "time"

type BreedingDB struct {
	Id              string     `db:"id"`
	BreedingDate    time.Time  `db:"breeding_date"`
	BirthStatus     string     `db:"birth_status"`
	PregnancyStatus string     `db:"pregnancy_status"`
	Observation     *string    `db:"observation"`

	AnimalId        string     `db:"animal_id"`
	AnimalName      string     `db:"animal_name"`
	AnimalTag       *string    `db:"animal_tag"`
	AnimalOrder     string     `db:"animal_order"`

	BullId          string     `db:"bull_id"`
	BullName        string     `db:"bull_name"`
	BullTag         *string    `db:"bull_tag"`

	ChildId         *string    `db:"child_id"`
	ChildName       *string    `db:"child_name"`
	ChildTag        *string    `db:"child_tag"`
	ChildSex        *string    `db:"child_sex"`
	ChildBirthDate  *time.Time `db:"child_birth_date"`

	CreatedAt       time.Time  `db:"created_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	UserId          string     `db:"user_id"`
}

type BreedingEntrySave struct {
	Id             string    `json:"id" db:"id"`
	AnimalId       string    `json:"animalId" db:"animal_id"`
	BreedingDate   time.Time `json:"breedingDate" db:"breeding_date"`
	BullId         string    `json:"bullId" db:"bull_id"`
	Observation    *string   `json:"observation" db:"observation"`
	Overwrite      bool      `json:"overwrite"`
	SkipValidation bool      `json:"skipValidation"`
	UserId         string    `json:"-" db:"user_id"`
}

type BreedingEntryFilter struct {
	Animals         *[]string  `json:"animals" db:"animal_id"`
	Bulls           *[]string  `json:"bulls" db:"bull_id"`
	MinBreedingDate *time.Time `json:"minBreedingDate" db:"breeding_date"`
	MaxBreedingDate *time.Time `json:"maxBreedingDate" db:"breeding_date"`
	BirthStatus     *string    `json:"birthStatus" db:"birth_status"`
	PregnancyStatus *string    `json:"pregnancyStatus" db:"pregnancy_status"`
}

type BreedingFoot struct {
	Totals               int     `json:"totals" db:"totals"`
	AverageBirthRate     float64 `json:"averageBirthRate" db:"average_birth_rate"`
	AveragePregnancyRate float64 `json:"averagePregnancyRate" db:"average_pregnancy_rate"`
}

type BreedingHist struct {
	BreedingDate      time.Time `json:"breedingDate" db:"breeding_date"`
	AnimalsNumber     int       `json:"animalsNumber" db:"animals_number"`
	PregnanciesNumber int       `json:"pregnanciesNumber" db:"pregnancies_number"`
	BirthsNumber      int       `json:"birthsNumber" db:"births_number"`
}

type BestBulls struct {
	BullName                string  `json:"bullName" db:"bull_name"`
	Total                   int     `json:"total" db:"total"`
	BirthRate               float64 `json:"birthRate" db:"birth_rate"`
	PregnancyRate           float64 `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthComparisonRate     float64 `json:"birthComparisonRate" db:"birth_comparison_rate"`
	PregnancyComparisonRate float64 `json:"pregnancyComparisonRate" db:"pregnancy_comparison_rate"`
}

type BreedingGroup struct {
	BullId                  string    `json:"bullId,omitempty" db:"bull_id"`
	BullName                string    `json:"bullName" db:"bull_name"`
	BreedingDate            time.Time `json:"breedingDate" db:"breeding_date"`
	CowNumber               float64   `json:"cowNumber" db:"cow_number"`
	BirthRate               float64   `json:"birthRate" db:"birth_rate"`
	PregnancyRate           float64   `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthComparisonRate     float64   `json:"birthComparisonRate" db:"birth_comparison_rate"`
	PregnancyComparisonRate float64   `json:"pregnancyComparisonRate" db:"pregnancy_comparison_rate"`
	UserId                  string    `json:"-" db:"user_id"`
}
