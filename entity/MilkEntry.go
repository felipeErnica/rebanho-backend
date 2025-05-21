package entity

import (
	"time"
)

type MilkEntry struct {
	Id           string       `json:"id"`
	Animal       Animal       `json:"animal"`
	Pasture      PastureShort `json:"pasturte"`
	LactationId  string       `json:"lactation_id"`
	EntryDate    time.Time    `json:"entry_date"`
	MilkQuantity float32      `json:"milk_quantity"`
	CreatedAt    time.Time    `json:"created_at"`
	DeletedAt    *time.Time   `json:"deleted_at"`
	UserId       string       `json:"user_id"`
}
