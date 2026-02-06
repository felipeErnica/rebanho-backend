package birth

import "time"

type BirthDTO struct {
	Mother        Parent  `json:"mother"`
	Father        *Parent `json:"father"`
	Calf          Calf    `json:"calf"`
	BirthInterval *int    `json:"birthInterval"`
}

type Parent struct {
	Id   string  `json:"id"`
	Name *string `json:"name"`
	Tag  *string `json:"tag"`
}

type Calf struct {
	Id          string    `json:"id"`
	Name        *string   `json:"name"`
	Tag         *string   `json:"tag"`
	Sex         string    `json:"sex"`
	BirthDate   time.Time `json:"birthDate"`
	Observation *string   `json:"observation"`
}
