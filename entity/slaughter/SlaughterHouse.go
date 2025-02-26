package slaughter

import "github.com/google/uuid"

type SlaughterHouse struct {
	Id             string
	Name           string
	TaxNumber      string
	WeightDecrease float32
	UserId         string
}

func (s *SlaughterHouse) New(c *CreateSlaughterHouse) *SlaughterHouse {
	s = &SlaughterHouse{
		Id: uuid.NewString(),
        Name: c.Name,
        TaxNumber: c.TaxNumber,
        WeightDecrease: c.WeightDecrease,
        UserId: c.UserId,
	}
    return s
}

type CreateSlaughterHouse struct {
	Name           string
	TaxNumber      string
	WeightDecrease float32
	UserId         string
}
