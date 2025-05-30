package slaughterhouses

import "time"

type Slaughterhouse struct {
	Id             string     `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	TaxNumber      string     `json:"taxNumber" db:"tax_number"`
	City           string     `json:"city" db:"city"`
	State          string     `json:"state" db:"state"`
	WeightDiscount float64    `json:"weightDiscount" db:"weight_discount"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt      *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId         string     `json:"userId" db:"user_id"`
}
