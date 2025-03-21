package entity

import "time"

type BirthEntry struct {
	Id          string
	CalfId      string
	Observation string
	CreatedAt   time.Time
}
