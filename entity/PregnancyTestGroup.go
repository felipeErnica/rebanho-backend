package entity

import "time"

type PregancyTestGroup struct {
	Id              string     `json:"id"`
	TestDate        string     `json:"test_date"`
	NumberEntries   int        `json:"number_entries"`
	NumberPregnants int        `json:"number_pregnant"`
	CreatedAt       time.Time  `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	UserId          string     `json:"user_id"`
}
