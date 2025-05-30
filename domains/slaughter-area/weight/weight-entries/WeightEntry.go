package weightEntries

import "time"

type WeightEntry struct {
	Id           string    `json:"id" db:"id"`
	AnimalId     string    `json:"animalId" db:"animal_id"`
	AnimalOrder  string    `json:"animalOrder" db:"animal_order"`
	AnimalNumber string    `json:"animalNumber" db:"animal_number"`
	AnimalName   string    `json:"animalName" db:"animal_name"`
	AnimalSex    string    `json:"animalSex" db:"animal_sex"`
	GroupId      string    `json:"groupId" db:"group_id"`
	GroupDate    time.Time `json:"groupDate" db:"group_date"`
	Weight       float32   `json:"weight" db:"weight"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	DeletedAt    time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string    `json:"userId" db:"user_id"`
}

type WeightEntrySave struct {
	Id           string    `json:"id" db:"id"`
	AnimalId     string    `json:"animalId" db:"animal_id"`
	GroupId      string    `json:"groupId" db:"group_id"`
	Weight       float32   `json:"weight" db:"weight"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	DeletedAt    time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string    `json:"userId" db:"user_id"`
}
