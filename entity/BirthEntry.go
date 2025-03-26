package entity

import "time"

type BirthEntry struct {
	Id          string
    Animal      AnimalShort
	Calf        Calf
	Observation string
	CreatedAt   time.Time
}
