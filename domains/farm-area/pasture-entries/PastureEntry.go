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

type PastureEntrySave struct {
	Id        string     `json:"id" db:"id"`
	AnimalId  string     `json:"animalId" db:"animal_id"`
	PastureId string     `json:"pastureId" db:"pasture_id"`
	EntryDate time.Time  `json:"entryDate" db:"entry_date"`
	ExitDate  *time.Time `json:"exitDate" db:"exit_date"`
	UserId    string     `json:"userId" db:"user_id"`
}
