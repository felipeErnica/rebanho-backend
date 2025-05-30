package testGroup

import "time"

type TestGroup struct {
	Id              string     `json:"id" db:"id"`
	TestDate        time.Time  `json:"testDate" db:"test_date"`
	NumberEntries   int        `json:"numberEntries" db:"number_entries"`
	NumberPregnants int        `json:"numberPregnants" db:"number_pregnant"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt       *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId          string     `json:"userId" db:"user_id"`
}
