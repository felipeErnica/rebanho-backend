package pastureEntries

import (
	"time"

	"github.com/google/uuid"
)

type PastureEntry struct {
	Id               uuid.UUID  `json:"id" db:"id"`
	AnimalId         uuid.UUID  `json:"animalId" db:"animal_id"`
	AnimalRingNumber *string    `json:"animalRingNumber" db:"animal_ring_number"`
	AnimalOrder      int        `db:"animal_order"`
	AnimalName       *string    `json:"animalName" db:"animal_name"`
	AnimalBirthDate  *time.Time `json:"animalBirthDate" db:"animal_birth_date"`
	EntryDate        time.Time  `json:"entryDate" db:"entry_date"`
	CreatedAt        time.Time  `db:"created_at"`
	DeletedAt        *time.Time `db:"deleted_at"`
	UserId           uuid.UUID  `db:"user_id"`
}

type PastureEntryFilter struct {
	IsFiltered         bool      `json:"isFiltered" db:"is_filtered"`
	AnimalRingNumber   string    `json:"animalRingNumber" db:"animal_ring_number"`
	AnimalName         string    `json:"animalName" db:"animal_name"`
	MinAnimalBirthDate time.Time `json:"minAnimalBirthDate" db:"animal_birth_date"`
	MaxAnimalBirthDate time.Time `json:"maxAnimalBirthDate" db:"animal_birth_date"`
	MaxEntryDate       time.Time `json:"maxEntryDate" db:"entry_date"`
	MinEntryDate       time.Time `json:"minEntryDate" db:"entry_date"`
}

type PastureEntrySave struct {
	Id        string     `json:"id" db:"id"`
	AnimalId  string     `json:"animalId" db:"animal_id"`
	PastureId string     `json:"pastureId" db:"pasture_id"`
	EntryDate time.Time  `json:"entryDate" db:"entry_date"`
	ExitDate  *time.Time `json:"exitDate" db:"exit_date"`
	UserId    string     `json:"userId" db:"user_id"`
}
