package weight

import (
	"time"

	"github.com/google/uuid"
)

type WeightGroup struct {
	Id         string
	WeightDate time.Time
	UserId     string
}

func (w *WeightGroup) New(create *CreateWeightGroup) *WeightGroup {
    w = &WeightGroup{
        Id: uuid.NewString(),
        WeightDate: create.WeightDate,
        UserId: create.UserId,
    }
    return w
}

type CreateWeightGroup struct {
	WeightDate time.Time
	UserId     string
}
