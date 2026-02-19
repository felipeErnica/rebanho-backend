package slaughter

import "time"

type SlaughterDTO struct {
	Id              string     `json:"id"`
	Animal          *Animal    `json:"animal"`
	Butcher         Butcher    `json:"butcher"`
	EntryDate       *time.Time `json:"entryDate"`
	DiscountRate    *float64   `json:"discountRate"`
	Weight          float64    `json:"weight"`
	DiscountWeight  *float64   `json:"discountWeight"`
	DeadWeight      *float64   `json:"deadWeight"`
	PerformanceRate *float64   `json:"performanceRate"`
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
