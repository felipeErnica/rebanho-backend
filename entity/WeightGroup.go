package entity

import "time"

type WeightGroup struct {
	Id         string     `json:"string"`
	WeightDate time.Time  `json:"weight_date"`
	CreatedAt  time.Time  `json:"created_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
	UserId     string     `json:"user_id"`
}
