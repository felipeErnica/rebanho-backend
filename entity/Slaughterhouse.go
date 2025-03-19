package entity

import "time"

type Slaughterhouse struct {
    Id         string     `json:"id"`
    Name       string     `json:"name"`
    TaxNumber  string     `json:"tax_number"`
    CreatedAt  time.Time  `json:"created_at"`
    UserId     string     `json:"user_id"`
}

type SlaughterhouseShort struct {
    Id    string `json:"id"`
	Name  string `json:"name"`
}
