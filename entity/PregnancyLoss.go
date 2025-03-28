package entity

import "time"

type PregnancyLoss struct {
	Id        string      `json:"id"`
	Animal    AnimalShort `json:"animal"`
	LossType  string      `json:"loss_type"`
	LossDate  time.Time   `json:"loss_date"`
	CreatedAt time.Time   `json:"created_at"`
	DeletedAt *time.Time  `json:"deleted_at"`
    UserId    string      `json:"user_id"`
}

type LossShort struct {
	Id       string    `json:"id"`
	LossType string    `json:"loss_type"`
	LossDate time.Time `json:"loss_date"`
}
