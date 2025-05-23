package entity

import "time"

type Slaughterhouse struct {
	Id             string     `json:"id"`
	Name           string     `json:"name"`
	TaxNumber      string     `json:"taxNumber"`
	City           string     `json:"city"`
	State          string     `json:"state"`
	WeightDiscount float64    `json:"weightDiscount"`
	CreatedAt      time.Time  `json:"createdAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	UserId         string     `json:"userId"`
}
