package entity

import "time"

type Farm struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	State     string    `json:"state"`
	City      string    `json:"city"`
	TaxNumber string    `json:"tax_number"`
    Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	DeletedAt time.Time `json:"deleted_at"`
	Owner     User      `json:"owner"`
}

type FarmFilter struct {
	Name      *string `json:"name"`
	State     *string `json:"state"`
	City      *string `json:"city"`
	TaxNumber *string `json:"tax_number"`
	OwnerName *string `json:"owner_name"`
}
