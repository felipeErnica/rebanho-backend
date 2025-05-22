package entity

import "time"

type NaturalMating struct {
	Id           string     `json:"id"`
	AnimalId     string     `json:"animalId"`
	AnimalNumber string     `json:"animalNumber"`
	AnimalName   string     `json:"animalName"`
	MatingDate   time.Time  `json:"matingDate"`
	BullId       string     `json:"bullId"`
	BullName     string     `json:"bullName"`
	Observation  string     `json:"observation"`
	Status       string     `json:"status"`
	LossId       string     `json:"lossId"`
	CalfId       string     `json:"calfId"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	UserId       string     `json:"userId"`
}

type NaturalMatingSave struct {
	Id           string     `json:"id"`
	AnimalId     string     `json:"animalId"`
	MatingDate   time.Time  `json:"matingDate"`
	BullId       string     `json:"bullId"`
	Observation  string     `json:"observation"`
	Status       string     `json:"status"`
	LossId       string     `json:"lossId"`
	CalfId       string     `json:"calfId"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	UserId       string     `json:"userId"`
}
