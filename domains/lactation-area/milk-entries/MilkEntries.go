package milkEntries

import "time"

type MilkEntry struct {
	Id           string     `json:"id" db:"id"`
	AnimalId     string     `json:"animalId" db:"animal_id"`
	AnimalOrder  int        `json:"animalOrder" db:"animal_order"`
	AnimalNumber string     `json:"animalNumber" db:"animal_number"`
	AnimalName   string     `json:"animalName" db:"animal_name"`
	PastureId    string     `json:"pastureId" db:"pasture_id"`
	PastureName  string     `json:"pastureName" db:"pasture_name"`
	LactationId  string     `json:"lactationId" db:"lactation_id"`
	EntryDate    time.Time  `json:"entryDate" db:"entry_date"`
	MilkQuantity float32    `json:"milkQuantity" db:"milk_quantity"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string     `json:"userId" db:"user_id"`
}

type MilkEntrySave struct {
	Id           string     `json:"id" db:"id"`
	AnimalId     string     `json:"animalId" db:"animal_id"`
	PastureId    string     `json:"pastureId" db:"pasture_id"`
	LactationId  string     `json:"lactationId" db:"lactation_id"`
	EntryDate    time.Time  `json:"entryDate" db:"entry_date"`
	MilkQuantity float32    `json:"milkQuantity" db:"milk_quantity"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string     `json:"userId" db:"user_id"`
}

type MilkEntryFilter struct {
	AnimalId        *[]string  `json:"animalId" db:"animal_id"`
	PastureId       *[]string  `json:"pastureId" db:"pasture_id"`
	MinEntryDate    *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate    *time.Time `json:"maxEntryDate" db:"entry_date"`
	MinMilkQuantity *float32   `json:"minMilkQuantity" db:"milk_quantity"`
	MaxMilkQuantity *float32   `json:"maxMilkQuantity" db:"milk_quantity"`
}
