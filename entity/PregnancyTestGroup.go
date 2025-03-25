package entity

import "time"

type PregancyTestGroup struct {
	Id              string    `json:"id"`
	TestDate        string    `json:"test_date"`
	NumberEntries   int       `json:"number_entries"`
	NumberPregnants int       `json:"number_pregnant"`
	UserId          string    `json:"user_id"`
	CreatedAt       time.Time `json:"created_at"`
}
