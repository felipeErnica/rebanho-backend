package entity

import (
	"time"

	"github.com/google/uuid"
)

type MilkEntry struct {
	Id           string
	AnimalId     string
	PastureId    string
	LactationId  string
	EntryDate     time.Time
	MilkQuantity float32
}

func (m *MilkEntry) New(create *CreateMilkMark) *MilkEntry {
    m = &MilkEntry{
        Id: uuid.NewString(),
        AnimalId: create.AnimalId,
        PastureId:  create.PastureId,
        LactationId: create.LactationId,
        EntryDate: create.MarkDate,
        MilkQuantity: create.MilkQuantity,
    }
    return m
}

type CreateMilkMark struct {
	AnimalId     string
	PastureId    string
	LactationId  string
	MarkDate     time.Time
	MilkQuantity float32
}
