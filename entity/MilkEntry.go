package entity

import "time"

type MilkEntry struct {
	Id           string     `json:"id"`
	AnimalId     string     `json:"animalId"`
	AnimalOrder  int        `json:"animalOrder"`
	AnimalNumber string     `json:"animalNumber"`
	AnimalName   string     `json:"animalName"`
	PastureId    string     `json:"pastureId"`
	PastureName  string     `json:"pastureName"`
	LactationId  string     `json:"lactationId"`
	EntryDate    time.Time  `json:"entryDate"`
	MilkQuantity float32    `json:"milkQuantity"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	UserId       string     `json:"userId"`
}
