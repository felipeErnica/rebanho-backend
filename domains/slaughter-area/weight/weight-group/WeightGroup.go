package weightGroup

import "time"

type WeightGroup struct {
    Id         string     `json:"id" db:"string"`
    WeightDate time.Time  `json:"weightDate" db:"weight_date"`
    CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
    DeletedAt  *time.Time `json:"deletedAt" db:"deleted_at"`
    UserId     string     `json:"userId" db:"user_id"`
}
