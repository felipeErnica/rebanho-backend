package entity

import (
	"time"
)

type PastureEntry struct {
	Id        string       `json:"id"`
	Animal    AnimalDto    `json:"animal"`
	Pasture   PastureShort `json:"pasture"`
	Bull      AnimalDto    `json:"bull"`
	EntryDate time.Time    `json:"entry_date"`
	ExitDate  time.Time    `json:"exit_date"`
	CreatedAt time.Time    `json:"created_at"`
	DeletedAt *time.Time   `json:"deleted_at"`
	UserId    string       `json:"user_id"`
}
