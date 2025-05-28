package entity

import "time"

type Pasture struct {
	Id        string     `json:"id"`
	BullId    string     `json:"bullId"`
	BullName  string     `json:"bullName"`
	Name      string     `json:"name"`
	Farm      FarmShort  `json:"farm"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	UserId    string     `json:"user_id"`
}

type PastureShort struct {
	Id   *string `json:"id"`
	Name *string `json:"name"`
}
