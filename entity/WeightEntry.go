package entity

import "time"

type WeightEntry struct {
	Id           string    `json:"id"`
	AnimalId     string    `json:"animalId"`
	AnimalOrder  string    `json:"animalOrder"`
	AnimalNumber string    `json:"animalNumber"`
	AnimalName   string    `json:"animalName"`
	AnimalSex    string    `json:"animalSex"`
	GroupId      string    `json:"groupId"`
	GroupDate    time.Time `json:"groupDate"`
	Weight       float32   `json:"weight"`
	CreatedAt    time.Time `json:"createdAt"`
	DeletedAt    time.Time `json:"deletedAt"`
	UserId       string    `json:"userId"`
}
