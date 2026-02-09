package lactation

import "time"

type AnimalDTO struct {
	Id        string        `json:"id"`
	Tag       *string       `json:"tag"`
	Name      *string        `json:"name"`
	BirthDate *time.Time    `json:"birthDate"`
	Lactation *LactationDTO `json:"lactation"`
}

type LactationDTO struct {
	Id                string     `json:"id"`
	Cow               Animal     `json:"cow"`
	Calf              *Calf      `json:"calf"`
	StartDate         time.Time  `json:"startDate"`
	EndDate           *time.Time `json:"endDate"`
	LacPeriod         *float64   `json:"lacPeriod"`
	AverageProduction *float64   `json:"averageProduction"`
	TotalProduction   *float64   `json:"totalProduction"`
	LacInterval       *int       `json:"lacInterval"`
	Peak              *float64   `json:"peak"`
	Observation       *string    `json:"observation"`
}

type Animal struct {
	Id   string `json:"id"`
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

type Calf struct {
	Id        string     `json:"id"`
	Tag       *string    `json:"tag"`
	Sex       string     `json:"sex"`
	Name      *string    `json:"name"`
	BirthDate *time.Time `json:"birthDate"`
	DeathDate *time.Time `json:"deathDate"`
}
