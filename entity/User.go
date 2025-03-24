package entity

import "time"

type User struct {
	Id        string     `json:"id"`
	Username  string     `json:"username"`
    Password  string     `json:"password"`
    CreatedAt time.Time  `json:"created_at"`
    DeletedAt *time.Time `json:"deleted_at"`
}
