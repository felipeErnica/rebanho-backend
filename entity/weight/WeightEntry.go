package weight

import "github.com/google/uuid"

type WeightEntry struct {
	Id       string
	AnimalId string
	GroupId  string
	Weight   float32
}

func (w *WeightEntry) New(c *CreateWeightEntry) *WeightEntry {
	w = &WeightEntry{
		Id: uuid.NewString(),
        AnimalId: c.AnimalId,
        GroupId: c.GroupId,
        Weight: c.Weight,
	}
    return w
}

type CreateWeightEntry struct {
	AnimalId string
	GroupId  string
	Weight   float32
}
