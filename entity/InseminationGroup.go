package entity

import "time"

type InseminationGroup struct {
	Id               string            `json:"id"`
	InseminationDate time.Time         `json:"insemination_date"`
	Bull             BullInsemintation `json:"bull"`
	UserId           string            `json:"user_id"`
	CreatedAt        time.Time         `json:"created_at"`
    DeletedAt        *time.Time        `json:"deleted_at"`
}
