package entity

import (
	"time"
)

type InseminationEntry struct {
	Id          string      `json:"id"`
	Animal      AnimalShort `json:"animal"`
	GroupId     string      `json:"GroupId"`
	Observation string      `json:"observation"`
	Status      string      `json:"status"`
	Loss        LossShort   `json:"loss"`
	Calf        CalfShort   `json:"calf"`
	CreatedAt   time.Time   `json:"created_at"`
	DeletedAt   *time.Time  `json:"deleted_at"`
}
