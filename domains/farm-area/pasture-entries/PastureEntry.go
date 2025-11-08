package pastureEntries

import (
	"time"
)

type PastureEntry struct {
	Id              string     `json:"id" db:"id"`
	AnimalId        string     `json:"animalId" db:"animal_id"`
	PastureId       string     `json:"pastureId" db:"pasture_id"`
	AnimalName      *string    `json:"animalName" db:"animal_name"`
	AnimalOrder     int        `db:"animal_order"`
	AnimalBirthDate *time.Time `json:"animalBirthDate" db:"animal_birth_date"`
	AnimalMother    *string    `json:"animalMother" db:"animal_mother"`
	AnimalFather    *string    `json:"animalFather" db:"animal_father"`
	EntryDate       time.Time  `json:"entryDate" db:"entry_date"`
	ExitDate        *time.Time `json:"exitDate" db:"exit_date"`
	CreatedAt       time.Time  `db:"created_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	UserId          string     `db:"user_id"`
}

type PastureTotal struct {
	Total int `json:"total" db:"total"`
}

type PastureEntryFilter struct {
	IsFiltered         bool       `json:"isFiltered" db:"is_filtered"`
	AnimalRingNumber   *string    `json:"animalRingNumber" db:"ring_number"`
	AnimalName         *string    `json:"animalName" db:"name"`
	Fathers            *[]string  `json:"fathers" db:"father_id"`
	Mothers            *[]string  `json:"mothers" db:"mother_id"`
	MinAnimalBirthDate *time.Time `json:"minAnimalBirthDate" db:"birth_date"`
	MaxAnimalBirthDate *time.Time `json:"maxAnimalBirthDate" db:"birth_date"`
	MaxEntryDate       *time.Time `json:"maxEntryDate" db:"entry_date" table:"pasture_entries"`
	MinEntryDate       *time.Time `json:"minEntryDate" db:"entry_date" table:"pasture_entries"`
}

type PastureEntrySave struct {
	Id        string     `json:"id" db:"id"`
	AnimalId  string     `json:"animalId" db:"animal_id"`
	PastureId string     `json:"pastureId" db:"pasture_id"`
	EntryDate time.Time  `json:"entryDate" db:"entry_date"`
	ExitDate  *time.Time `json:"exitDate" db:"exit_date"`
	UserId    string     `json:"userId" db:"user_id"`
}
