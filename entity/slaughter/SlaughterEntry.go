package slaughter

import "github.com/google/uuid"

type SlaughterEntry struct {
	Id         string
	GroupId    string
	Weight     string
	DeadWeight string
}

func (s *SlaughterEntry) New(c *CreateSlaughterEntry) *SlaughterEntry {
	s = &SlaughterEntry{
		Id: uuid.NewString(),
        GroupId: c.GroupId,
        Weight: c.Weight,
        DeadWeight: c.DeadWeight,
	}
    return s
}

type CreateSlaughterEntry struct {
	GroupId    string
	Weight     string
	DeadWeight string
}
