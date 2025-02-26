package slaughter

import (
	"time"

	"github.com/google/uuid"
)

type SlaughterGroup struct {
	Id               string
	SlaughterHouseId string
	SlaughterDate    time.Time
}

func (s *SlaughterGroup) New(c *CreateSlaughterGroup) *SlaughterGroup {
    s = &SlaughterGroup{
        Id: uuid.NewString(),
        SlaughterHouseId: c.SlaughterHouseId,
        SlaughterDate: c.SlaughterDate,
    }
    return s
}

type CreateSlaughterGroup struct {
	SlaughterHouseId string
	SlaughterDate    time.Time
}
