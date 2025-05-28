package entity

import "time"

type PastureEntry struct {
	Id           string     `json:"id"`
	AnimalId     string     `json:"animalId"`
	AnimalName   string     `json:"animalName"`
	AnimalNumber string     `json:"animalNumber"`
	PastureId    string     `json:"pastureId"`
	PastureName  string     `json:"pastureName"`
	BullId       string     `json:"bullId"`
	BullName     string     `json:"bullName"`
	EntryDate    time.Time  `json:"entryDate"`
	ExitDate     time.Time  `json:"exitDate"`
	CreatedAt    time.Time  `json:"createdAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	UserId       string     `json:"userId"`
}
