package slaughtergroups

import "time"

type DTO struct {
	Id                string    `json:"id"`
	EntryDate         time.Time `json:"entryDate"`
	Discount          float64   `json:"discount"`
	Butcher           Butcher   `json:"butcher"`
	AnimalsNumber     *int      `json:"animalsNumber"`
	AverageWeight     *float64  `json:"averageWeight"`
	AverageDeadWeight *float64  `json:"averageDeadWeight"`
	AverageRate       *float64  `json:"averageRate"`
}

type Butcher struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Discount *float64 `json:"discount"`
}
