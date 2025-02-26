package insemination

import "github.com/google/uuid"

type InseminationEntry struct {
	Id          string
	AnimalId    string
	GroupId     string
	Observation string
}

func (i *InseminationEntry) New(c *CreateInseminationEntry) *InseminationEntry {
	i = &InseminationEntry{
		Id: uuid.NewString(),
        AnimalId: c.AnimalId,
        GroupId: c.GroupId,
        Observation: c.Observation,
	}
    return i
}

type CreateInseminationEntry struct {
	AnimalId    string
	GroupId     string
	Observation string
}
