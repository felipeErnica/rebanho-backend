package entity

import "time"

type BirthEntry struct {
	Id          string      `json:"id"`
	Animal      AnimalShort `json:"animal"`
	Calf        Calf        `json:"calf"`
	Observation string      `json:"observation"`
	CreatedAt   time.Time   `json:"created_at"`
	DeletedAt   *time.Time  `json:"deleted_at"`
    UserId      string      `json:"user_id"`
}
