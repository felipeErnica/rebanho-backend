package entity

import "time"

type BirthEntry struct {
	Id             string     `json:"id"`
	AnimalId       string     `json:"animalId"`
	AnimalName     string     `json:"animalName"`
	AnimalNumber   string     `json:"animalNumber"`
	CalfId         string     `json:"calfId"`
	CalfBirthDate  string     `json:"calfBirthDate"`
	CalfSex        string     `json:"calfSex"`
	CalfFatherName string     `json:"calfFatherName"`
	Observation    string     `json:"observation"`
	CreatedAt      time.Time  `json:"createdAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	UserId         string     `json:"userId"`
}
