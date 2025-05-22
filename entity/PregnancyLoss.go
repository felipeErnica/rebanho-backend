package entity

import "time"

type PregnancyLoss struct {
	Id           string     `json:"id"`
	AnimalId     string     `json:"animalId"`
	AnimalOrder  int        `json:"animalOrder"`
	AnimalNumber string     `json:"animalNumber"`
	AnimalName   string     `json:"animalName"`
	LossType     string     `json:"lossType"`
	LossDate     time.Time  `json:"lossDate"`
	Observation  string     `json:"observation"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	UserId       string     `json:"userId"`
}
