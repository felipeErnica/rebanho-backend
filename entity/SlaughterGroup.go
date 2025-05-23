package entity

import "time"

type SlaughterGroup struct {
	Id                 string     `json:"id"`
	SlaughterhouseId   string     `json:"slaughterhouseId"`
	SlaughterhouseName string     `json:"slaughterhouseName"`
	WeightDecrease     float32    `json:"weightDecrease"`
	SlaughterDate      time.Time  `json:"slaughterDate"`
	CreatedAt          time.Time  `json:"createdAt"`
	DeletedAt          *time.Time `json:"deletedAt"`
	UserId             string     `json:"userId"`
}
