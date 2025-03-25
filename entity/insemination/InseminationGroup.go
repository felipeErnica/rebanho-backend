package insemination

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type InseminationGroup struct {
    Id               string                    `json:"id"`
    InseminationDate time.Time                 `json:"insemination_date"`
    Bull             entity.BullInsemintation  `json:"bull"`
    UserId           string                    `json:"user_id"`
    CreatedAt        time.Time                 `json:"created_at"`
}
