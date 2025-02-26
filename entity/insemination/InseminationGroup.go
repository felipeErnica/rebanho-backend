package insemination

import (
	"time"

	"github.com/google/uuid"
)

type InseminationGroup struct {
	Id                string
	InseminationDate  time.Time
	BullId            string
	UserId            string
}

func (i *InseminationGroup) New(c *CreateInseminationGroup) *InseminationGroup {
    i = &InseminationGroup{
        Id: uuid.NewString(),
        InseminationDate: c.InseminationDate,
        BullId: c.BullId,
        UserId: c.UserId,
    }
    return i
}

type CreateInseminationGroup struct {
	InseminationDate  time.Time
	BullId            string
	UserId            string
}
