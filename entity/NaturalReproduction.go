package entity

import (
	"time"
)

type NaturalReproduction struct {
	Id           string     `json:"id"`
	AnimalId     string     `json:"animalId"`
	AnimalNumber string     `json:"animalNumber"`
    AnimalOrder  string     `json:"animalOrder"`
	AnimalName   string     `json:"animalName"`
	BullName     string     `json:"bullName"`
	Observation  string     `json:"observation"`
	Status       string     `json:"status"`
	LossId       string     `json:"lossId"`
	CalfId       string     `json:"calfId"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
	UserId       string     `json:"user_id"`
}

type NaturalReproductionSave struct {
	Id          string     `json:"id"`
	AnimalId    string     `json:"animal"`
	Observation string     `json:"observation"`
	Status      string     `json:"status"`
	LossId      string     `json:"loss"`
	CalfId      string     `json:"calf"`
	CreatedAt   time.Time  `json:"created_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	UserId      string     `json:"user_id"`
}
