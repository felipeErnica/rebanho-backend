package tests

import "time"

type TestDTO struct {
	Id              string     `json:"id"`
	TestDate        time.Time  `json:"testDate"`
	BirthForecast   *time.Time `json:"birthForecast"`
	BirthStatus     string     `json:"birthStatus"`
	PregnancyStatus string     `json:"pregnancyStatus"`
	Observation     *string    `json:"observation"`
	Cow             Cow        `json:"cow"`
	Calf            *Calf      `json:"calf"`
}

type Cow struct {
	Id   string  `json:"id"`
	Name *string `json:"name"`
	Tag  *string `json:"tag"`
}

type Calf struct {
	Id        string     `json:"id"`
	Name      *string    `json:"name"`
	Tag       *string    `json:"tag"`
	Sex       string     `json:"sex"`
	BirthDate time.Time  `json:"birthDate"`
	DeathDate *time.Time `json:"deathDate"`
}
