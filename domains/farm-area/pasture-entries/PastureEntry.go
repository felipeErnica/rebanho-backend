package pastureEntries

import "time"

type PastureEntry struct {
	Id           string     `json:"id" db:"id"`
	AnimalId     string     `json:"animalId" db:"animal_id"`
	AnimalName   string     `json:"animalName" db:"animal_name"`
	AnimalNumber string     `json:"animalNumber" db:"animal_number"`
	PastureId    string     `json:"pastureId" db:"pasture_id"`
	PastureName  string     `json:"pastureName" db:"pasture_name"`
	EntryDate    time.Time  `json:"entryDate" db:"entry_date"`
	ExitDate     time.Time  `json:"exitDate" db:"exit_date"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string     `json:"userId" db:"user_id"`
}

type PastureEntrySave struct {
	Id           string     `json:"id" db:"id"`
	AnimalId     string     `json:"animalId" db:"animal_id"`
	PastureId    string     `json:"pastureId" db:"pasture_id"`
	EntryDate    time.Time  `json:"entryDate" db:"entry_date"`
	ExitDate     time.Time  `json:"exitDate" db:"exit_date"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string     `json:"userId" db:"user_id"`
}
