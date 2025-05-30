package slaughterGroup

import "time"

type SlaughterGroup struct {
	Id                 string     `json:"id" db:"id"`
	SlaughterhouseId   string     `json:"slaughterhouseId" db:"slaughterhouse_id"`
	SlaughterhouseName string     `json:"slaughterhouseName" db:"slaughterhouse_name"`
	WeightDecrease     float32    `json:"weightDecrease" db:"weight_decrease"`
	SlaughterDate      time.Time  `json:"slaughterDate" db:"slaughter_date"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt          *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId             string     `json:"userId" db:"user_id"`
}

type SlaughterGroupSave struct {
	Id                 string     `json:"id" db:"id"`
	SlaughterhouseId   string     `json:"slaughterhouseId" db:"slaughterhouse_id"`
	WeightDecrease     float32    `json:"weightDecrease" db:"weight_decrease"`
	SlaughterDate      time.Time  `json:"slaughterDate" db:"slaughter_date"`
	CreatedAt          time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt          *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId             string     `json:"userId" db:"user_id"`
}
