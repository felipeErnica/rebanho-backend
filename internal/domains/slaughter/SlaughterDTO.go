package slaughter

import "time"

type SlaughterDTO struct {
	Id              string   `json:"id"`
	Group           Group    `json:"group"`
	Animal          *Animal  `json:"animal"`
	Weight          float64  `json:"weight"`
	DiscountWeight  *float64 `json:"discountWeight"`
	DeadWeight      *float64 `json:"deadWeight"`
	PerformanceRate *float64 `json:"performanceRate"`
}

type Group struct {
	Id        string    `json:"groupId"`
	EntryDate time.Time `json:"entryDate"`
	Discount  float64   `json:"discount"`
	Butcher   Butcher   `json:"butcher"`
}

type Animal struct {
	Id        string     `json:"id"`
	Tag       *string    `json:"tag"`
	Name      *string    `json:"name"`
	Sex       string     `json:"sex"`
	BirthDate *time.Time `json:"birthDate"`
	Father    *Parent    `json:"father"`
	Mother    *Parent    `json:"mother"`
}

type Parent struct {
	Id   string  `json:"id"`
	Tag  *string `json:"tag"`
	Name *string `json:"name"`
}

type Butcher struct {
	Id       string   `json:"id"`
	Name     string   `json:"name"`
	Discount *float64 `json:"discount"`
}

type SlaughterGroupDTO struct {
	EntryDate           time.Time `json:"entryDate"`
	Butcher             Butcher   `json:"butcher"`
	AnimalsNumber       int       `json:"animalsNumber"`
	AverageWeight       *float64  `json:"averageWeight"`
	AverageDeadWeight   *float64  `json:"averageDeadWeight"`
	WeightVariation     *float64  `json:"weightVariation"`
	DeadWeightVariation *float64  `json:"deadWeightVariation"`
	AverageRate         *float64  `json:"averageRate"`
	RateVariation       *float64  `json:"rateVariation"`
}
