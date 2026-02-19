package tests

import "time"

type TestDB struct {
	Id              string     `db:"id"`
	TestDate        time.Time  `db:"test_date"`
	AnimalId        string     `db:"animal_id"`
	AnimalTag       *string    `db:"animal_tag"`
	AnimalName      *string    `db:"animal_name"`
	AnimalOrder     int        `db:"animal_order"`
	BirthForecast   *time.Time `db:"birth_forecast"`
	BirthStatus     string     `db:"birth_status"`
	PregnancyStatus string     `db:"pregnancy_status"`
	Observation     *string    `db:"observation"`
	CalfId          *string    `db:"calf_id"`
	CalfTag         *string    `db:"calf_tag"`
	CalfName        *string    `db:"calf_name"`
	CalfSex         *string    `db:"calf_sex"`
	CalfBirthDate   *time.Time `db:"calf_birth_date"`
	CalfDeathDate   *time.Time `db:"calf_death_date"`
	CreatedAt       time.Time  `db:"created_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	UserId          string     `db:"user_id"`
}

type TestFilter struct {
	MinTestDate      *time.Time `json:"minTestDate" db:"test_date"`
	MaxTestDate      *time.Time `json:"maxTestDate" db:"test_date"`
	Animals          *[]string  `json:"animals" db:"animal_id"`
	MinBirthForecast *time.Time `json:"minBirthForecast" db:"birth_forecast"`
	MaxBirthForecast *time.Time `json:"maxBirthForecast" db:"birth_forecast"`
	BirthStatus      *string    `json:"birthStatus" db:"birth_status"`
	PregnancyStatus  *string    `json:"pregnancyStatus" db:"pregnancy_status"`
}

type TestFoot struct {
	Totals        int     `json:"totals" db:"totals"`
	PregnancyRate float64 `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthRate     float64 `json:"birthRate" db:"birth_rate"`
}

type TestSave struct {
	Id              string    `json:"id" db:"id"`
	TestDate        time.Time `json:"testDate" db:"test_date"`
	AnimalId        string    `json:"animalId" db:"animal_id"`
	PregnancyStatus string    `json:"pregnancyStatus" db:"pregnancy_status"`
	PregnancyTime   *int      `json:"pregnancyTime" db:"pregnancy_time"`
	Observation     *string   `json:"observation" db:"observation"`
	Overwrite       bool      `json:"overwrite"`
	UserId          string    `json:"-" db:"user_id"`
}

type TestGroups struct {
	TestDate            time.Time `json:"testDate" db:"test_date"`
	OldTestDate         time.Time `json:"oldTestDate" db:"old_test_date"`
	AnimalsNumber       int       `json:"animalsNumber" db:"animals_number"`
	PregnancyRate       float64   `json:"pregnancyRate" db:"pregnancy_rate"`
	PregnancyComparison float64   `json:"pregnancyComparison" db:"pregnancy_comparison"`
	BirthRate           float64   `json:"birthRate" db:"birth_rate"`
	BirthComparison     float64   `json:"birthComparison" db:"birth_comparison"`
	UserId              string    `json:"-" db:"user_id"`
}

type TestGroupSave struct {
	TestDate            time.Time `json:"testDate" db:"test_date"`
	OldTestDate         time.Time `json:"oldTestDate" db:"old_test_date"`
	UserId              string    `json:"-" db:"user_id"`
}

type TestAnimal struct {
	AnimalName          string  `json:"animalName" db:"animal_name"`
	Totals              int     `json:"totals" db:"totals"`
	PregnancyRate       float64 `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthRate           float64 `json:"birthRate" db:"birth_rate"`
	PregnancyComparison float64 `json:"pregnancyComparison" db:"pregnancy_comparison"`
	BirthComparison     float64 `json:"birthComparison" db:"birth_comparison"`
}

type PregnancyTestHist struct {
	TestDate    time.Time `json:"testDate" db:"test_date"`
	Totals      int       `json:"totals" db:"totals"`
	Pregnancies int       `json:"pregnancies" db:"pregnancies"`
	Births      int       `json:"births" db:"births"`
}
