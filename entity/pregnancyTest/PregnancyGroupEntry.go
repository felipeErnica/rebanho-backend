package pregnancytest

import "github.com/google/uuid"

type PregnancyTestEntry struct {
	Id        string
	GroupId   string
	AnimalId  string
	IsPregnant bool
}

func (p *PregnancyTestEntry) New(c *CreatePregnancyTestEntry) *PregnancyTestEntry {
	p = &PregnancyTestEntry{
		Id: uuid.NewString(),
        GroupId: c.GroupId,
        AnimalId: c.AnimalId,
        IsPregnant: c.IsPregnant,
	}
    return p
}

type CreatePregnancyTestEntry struct {
	GroupId   string
	AnimalId  string
	IsPregnant bool
}
