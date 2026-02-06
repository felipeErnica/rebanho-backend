package breeding

import "time"

type BreedingDTO struct {
	Id              string    `json:"id"`
	BreedingDate    time.Time `json:"breedingDate"`
	BirthStatus     string    `json:"birthStatus"`
	PregnancyStatus string    `json:"pregnancyStatus"`
	Observation     *string   `json:"observation"`
	Cow             Parent    `json:"cow"`
	Bull            Parent    `json:"bull"`
	Child           *Child     `json:"child"`
}

type Parent struct {
	Id   string  `json:"id"`
	Name string  `json:"name"`
	Tag  *string `json:"tag"`
}

type Child struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Tag       *string   `json:"tag"`
	Sex       string    `json:"sex"`
	BirthDate time.Time `json:"birthDate"`
}
