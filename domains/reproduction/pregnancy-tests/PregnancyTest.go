package pregnancyTests

import "time"

type TestEntry struct {
	Id               string     `json:"id,omitempty" db:"id"`
	TestDate         time.Time  `json:"testDate" db:"test_date"`
	AnimalId         string     `json:"animalId,omitempty" db:"animal_id"`
	AnimalName       string     `json:"animalName" db:"animal_name"`
	AnimalOrder      int        `json:"animalOrder,omitempty" db:"animal_order"`
	BirthForecast    *time.Time `json:"birthForecast,omitempty" db:"birth_forecast"`
	BirthStatus      string     `json:"birthStatus" db:"birth_status"`
	PregnancyStatus  string     `json:"pregnancyStatus" db:"pregnancy_status"`
	Observation      *string    `json:"observation" db:"observation"`
	ChildInformation *string    `json:"childInformation" db:"child_information"`
	LossId           *string    `json:"lossId" db:"loss_id"`
	CalfId           *string    `json:"calfId" db:"calf_id"`
	CreatedAt        time.Time  `json:"-" db:"created_at"`
	DeletedAt        *time.Time `json:"-" db:"deleted_at"`
	UserId           string     `json:"-" db:"user_id"`
}

type TestEntryFilter struct {
	IsFiltered       bool       `json:"isFiltered"`
	MinTestDate      *time.Time `json:"minTestDate" db:"test_date"`
	MaxTestDate      *time.Time `json:"maxTestDate" db:"test_date"`
	Animals          *[]string  `json:"animals" db:"animal_id"`
	MinBirthForecast *time.Time `json:"minBirthForecast" db:"birth_forecast"`
	MaxBirthForecast *time.Time `json:"maxBirthForecast" db:"birth_forecast"`
	BirthStatus      *string    `json:"birthStatus" db:"birth_status"`
	PregnancyStatus  *string    `json:"pregnancyStatus" db:"pregnancy_status"`
}

type TestEntryFoot struct {
	Totals        int     `json:"totals" db:"totals"`
	PregnancyRate float64 `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthRate     float64 `json:"birthRate" db:"birth_rate"`
}

type TestEntrySave struct {
	Id            string     `json:"id" db:"id"`
	GroupId       string     `json:"groupId" db:"group_id"`
	AnimalId      string     `json:"animalId" db:"animal_id"`
	IsPregnant    bool       `json:"isPregnant" db:"is_pregnant"`
	BirthForecast *time.Time `json:"birthForecast" db:"birth_forecast"`
	Observation   *string    `json:"observation" db:"observation"`
	Status        string     `json:"status" db:"status"`
	LossId        *string    `json:"lossId" db:"loss_id"`
	CalfId        *string    `json:"calfId" db:"calf_id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        string     `json:"userId" db:"user_id"`
}

type TestGroups struct {
	TestDate            time.Time `json:"testDate" db:"test_date"`
	AnimalsNumber       int       `json:"animalsNumber" db:"animals_number"`
	PregnancyRate       float64   `json:"pregnancyRate" db:"pregnancy_rate"`
	PregnancyComparison float64   `json:"pregnancyComparison" db:"pregnancy_comparison"`
	BirthRate           float64   `json:"birthRate" db:"birth_rate"`
	BirthComparison     float64   `json:"birthComparison" db:"birth_comparison"`
}

type LastEntries struct {
	TestDate time.Time   `json:"testDate"`
	Entries  []TestEntry `json:"entries"`
}

type NextBirths struct {
	BirthForecast time.Time `json:"birthForecast" db:"birth_forecast"`
	BirthNumbers  int       `json:"birthNumbers" db:"birth_numbers"`
}

type TestAnimal struct {
	AnimalName          string  `json:"animalName" db:"animal_name"`
	Totals              int     `json:"totals" db:"totals"`
	PregnancyRate       float64 `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthRate           float64 `json:"birthRate" db:"birth_rate"`
	PregnancyComparison float64 `json:"pregnancyComparison" db:"pregnancy_comparison"`
	BirthComparison     float64 `json:"birthComparison" db:"birth_comparison"`
}

type PregnancyHist struct {
	TestDate      time.Time `json:"testDate" db:"test_date"`
	PregnancyRate float64   `json:"pregnancyRate" db:"pregnancy_rate"`
}

type AnimalsNumberHist struct {
	TestDate      time.Time `json:"testDate" db:"test_date"`
	AnimalsNumber float64   `json:"animalsNumber" db:"totals"`
}

type CardStats struct {
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
	Hist    any     `json:"hist"`
}

type BirthHist struct {
	TestDate  time.Time `json:"testDate" db:"test_date"`
	BirthRate float64   `json:"birthRate" db:"birth_rate"`
}

type BirthStats struct {
	Current float64     `json:"current"`
	Trend   float64     `json:"trend"`
	Hist    []BirthHist `json:"hist"`
}

type PregnancyTestHist struct {
	TestDate    time.Time `json:"testDate" db:"test_date"`
	Totals      int       `json:"totals" db:"totals"`
	Pregnancies int       `json:"pregnancies" db:"pregnancies"`
	Births      int       `json:"births" db:"births"`
}
