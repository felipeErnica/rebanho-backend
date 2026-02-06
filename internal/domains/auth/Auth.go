package auth

import "time"

type User struct {
	Id           string     `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	EmailAddress string     `json:"emailAddress" db:"email_address"`
	PhoneNumber  string     `json:"phoneNumber" db:"phone_number"`
	Password     string     `json:"password" db:"password"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
}

type AuthToken struct {
	Token string `json:"token"`
}
