package insemination

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type InseminationEntry struct {
	Id          string             `json:"id"`
	Animal      entity.AnimalShort `json:"animal"`
	GroupId     string             `json:"GroupId"`
	Observation string             `json:"observation"`
	Status      string             `json:"status"`
	Loss        entity.LossShort   `json:"loss"`
	Calf        entity.CalfShort   `json:"calf"`
	CreatedAt   time.Time          `json:"created_at"`
	DeletedAt   *time.Time         `json:"deleted_at"`
}
