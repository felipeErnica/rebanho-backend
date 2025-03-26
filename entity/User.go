package entity

import "time"

type User struct {
	Id           string     `json:"id"`
	Name         string     `json:"name"`
	EmailAddress string     `json:"email_address"`
	PhoneNumber  string     `json:"phone_number"`
	Password     string     `json:"password"`
	CreatedAt    time.Time  `json:"created_at"`
	DeletedAt    *time.Time `json:"deleted_at"`
}
