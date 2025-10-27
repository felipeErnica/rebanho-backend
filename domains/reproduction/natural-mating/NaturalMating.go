package naturalMating

import "time"

type MatingEntry struct {
	Id               string     `json:"id" db:"id"`
	AnimalId         string     `json:"animalId" db:"animal_id"`
	AnimalNumber     string     `json:"animalNumber" db:"animal_number"`
	AnimalName       string     `json:"animalName" db:"animal_name"`
	MatingDate       time.Time  `json:"matingDate" db:"mating_date"`
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

type MatingEntryFilter struct {
	IsFiltered          bool       `json:"isFiltered"`
	Animals             *[]string  `json:"animals" db:"animal_id"`
	Bulls               *[]string  `json:"bulls" db:"bull_id"`
	MinInseminationDate *time.Time `json:"minMatingDate" db:"mating_date"`
	MaxInseminationDate *time.Time `json:"maxMatingDate" db:"mating_date"`
	BirthStatus         *string    `json:"birthStatus" db:"birth_status"`
	PregnancyStatus     *string    `json:"pregnancyStatus" db:"pregnancy_status"`
}

type MatingFoot struct {
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
	MatingDate time.Time `json:"matingDate" db:"mating_date"`
	BirthRate  float64   `json:"birthRate" db:"birth_rate"`
}

type PregnancyRateEntry struct {
	MatingDate    time.Time `json:"matingDate" db:"mating_date"`
	PregnancyRate float64   `json:"pregnancyRate" db:"pregnancy_rate"`
}

type AnimalsNumberEntry struct {
	MatingDate    time.Time `json:"matingDate" db:"mating_date"`
	AnimalsNumber float64   `json:"animalsNumber" db:"animalsNumber"`
}

type MatingHist struct {
	MatingDate        time.Time `json:"matingDate" db:"mating_date"`
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

type MatingGroup struct {
	BullId                  string    `json:"bullId,omitempty" db:"bull_id"`
	BullName                string    `json:"bullName" db:"bull_name"`
	MatingDate              time.Time `json:"matingDate" db:"mating_date"`
	CowNumber               float64   `json:"cowNumber" db:"cow_number"`
	BirthRate               float64   `json:"birthRate" db:"birth_rate"`
	PregnancyRate           float64   `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthComparisonRate     float64   `json:"birthComparisonRate" db:"birth_comparison_rate"`
	PregnancyComparisonRate float64   `json:"pregnancyComparisonRate" db:"pregnancy_comparison_rate"`
}

type LastEntry struct {
	MatingDate time.Time     `json:"matingDate"`
	Entries    []MatingEntry `json:"entries"`
}
