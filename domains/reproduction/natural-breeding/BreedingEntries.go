package naturalBreeding

import "time"

type BreedingEntry struct {
	Id               string     `json:"id" db:"id"`
	AnimalId         string     `json:"animalId" db:"animal_id"`
	AnimalOrder      string     `json:"-" db:"animal_number"`
	AnimalName       string     `json:"-" db:"animal_name"`
	AnimalInfo       string     `json:"animalInfo" db:"animal_info"`
	BreedingDate     time.Time  `json:"breedingDate" db:"breeding_date"`
	BullId           string     `json:"bullId" db:"bull_id"`
	BullName         string     `json:"bullName" db:"bull_name"`
	BirthStatus      string     `json:"birthStatus" db:"birth_status"`
	PregnancyStatus  string     `json:"pregnancyStatus" db:"pregnancy_status"`
	Observation      *string    `json:"observation" db:"observation"`
	ChildInformation *string    `json:"childInformation" db:"child_information"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt        *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId           string     `json:"userId" db:"user_id"`
}

type BreedingEntrySave struct {
	Id           string    `json:"id" db:"id"`
	AnimalId     string    `json:"animalId" db:"animal_id"`
	BreedingDate time.Time `json:"breedingDate" db:"breeding_date"`
	BullId       string    `json:"bullId" db:"bull_id"`
	Observation  *string   `json:"observation" db:"observation"`
	UserId       string    `json:"-" db:"user_id"`
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

type CardEntry struct {
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
	Hist    any     `json:"hist"`
}

type BirthRateEntry struct {
	BreedingDate time.Time `json:"breedingDate" db:"breeding_date"`
	BirthRate    float64   `json:"birthRate" db:"birth_rate"`
}

type PregnancyRateEntry struct {
	BreedingDate  time.Time `json:"breedingDate" db:"breeding_date"`
	PregnancyRate float64   `json:"pregnancyRate" db:"pregnancy_rate"`
}

type AnimalsNumberEntry struct {
	BreedingDate  time.Time `json:"breedingDate" db:"breeding_date"`
	AnimalsNumber float64   `json:"animalsNumber" db:"animals_number"`
}

type BreedingHist struct {
	BreedingDate      time.Time `json:"breedingDate" db:"breeding_date"`
	AnimalsNumber     int       `json:"animalsNumber" db:"animals_number"`
	PregnanciesNumber int       `json:"pregnanciesNumber" db:"pregnancies_number"`
	BirthsNumber      int       `json:"birthsNumber" db:"births_number"`
}

type FutureBirths struct {
	BirthForecast time.Time `json:"birthForecast" db:"birth_forecast"`
	BirthsNumber  int       `json:"birthsNumber" db:"births_number"`
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

type LastEntry struct {
	BreedingDate time.Time       `json:"breedingDate"`
	Entries      []BreedingEntry `json:"entries"`
}
