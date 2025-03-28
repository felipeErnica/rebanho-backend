package entity

import "time"

type SlaughterEntry struct {
	Id         string              `json:"id"`
	Group      SlaughterGroupShort `json:"group_id"`
	Weight     string              `json:"weight"`
	DeadWeight string              `json:"dead_weight"`
	CreatedAt  time.Time           `json:"created_at"`
	DeletedAt  *time.Time          `json:"deleted_at"`
	UserId     string              `json:"user_id"`
}
