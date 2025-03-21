package entity

import "time"

type PregnancyTestEntry struct {
	Id            string      `json:"id"`
	GroupId       string      `json:"group"`
	Animal        AnimalShort `json:"animal"`
	IsPregnant    bool        `json:"is_pregnant"`
	BirthForecast time.Time   `json:"birth_forecast"`
    Observation   string      `json:"observation"`
    Loss          LossShort   `json:"loss"`
    Calf          CalfShort   `json:"calf"`
	CreatedAt     time.Time   `json:"created_at"`
}
