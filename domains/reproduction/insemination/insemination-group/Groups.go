package inseminationGroup

import "time"

type Group struct {
    Id               string     `json:"id" db:"id"`
    InseminationDate time.Time  `json:"inseminationDate" db:"insemination_date"`
    BullId           string     `json:"bullId" db:"bull_id"`
    BullName         string     `json:"bullName" db:"bull_name"`
    UserId           string     `json:"userId" db:"user_id"`
    CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
    DeletedAt        *time.Time `json:"deletedAt" db:"deleted_at"`
}

type GroupSave struct {
    Id               string     `json:"id" db:"id"`
    InseminationDate time.Time  `json:"inseminationDate" db:"insemination_date"`
    BullId           string     `json:"bullId" db:"bull_id"`
    UserId           string     `json:"userId" db:"user_id"`
    CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
    DeletedAt        *time.Time `json:"deletedAt" db:"deleted_at"`
}
