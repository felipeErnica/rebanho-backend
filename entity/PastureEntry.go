package entity

import (
	"time"
)

type PastureEntry struct {
    Id          string          `json:"id"`
    Animal      AnimalShort     `json:"animal"`
    Pasture     PastureShort    `json:"pasture"`
    EntryDate   time.Time       `json:"entry_date"`
    ExitDate    time.Time       `json:"exit_date"`
    CreatedAt   time.Time       `json:"created_at"`
    DeletedAt   *time.Time      `json:"deleted_at"`
}
