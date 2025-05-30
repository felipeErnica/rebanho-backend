package farm

import "time"

type Farm struct {
	Id        string     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	State     string     `json:"state" db:"state"`
	City      string     `json:"city" db:"city"`
	TaxNumber string     `json:"taxNumber" db:"tax_number"`
	Status    string     `json:"status" db:"status"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId    string     `json:"userId" db:"user_id"`
}

type FarmFilter struct {
	Name      string `json:"name" db:"name"`
	State     string `json:"state" db:"state"`
	City      string `json:"city" db:"city"`
	TaxNumber string `json:"taxNumber" db:"tax_number"`
	Status    string `json:"status" db:"status"`
}
