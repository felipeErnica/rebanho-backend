package entity

import "github.com/google/uuid"

type Pasture struct {
	Id     string
	BullId string
	Name   string
}

func (p *Pasture) NewPasture(create *CreatePasture) *Pasture {
	p = &Pasture{
		Id: uuid.NewString(),
        BullId: create.BullId,
        Name: create.BullId,
	}
    return p
}

type PastureShort struct {
	Id     *string
	Name   *string
}

type CreatePasture struct {
	BullId string
	Name   string
}
