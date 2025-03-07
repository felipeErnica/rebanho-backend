package entity

import (
	"time"

	"github.com/google/uuid"
)

type PastureEntry struct {
    Id        string
    AnimalId  string
    PastureId string
    EntryDate time.Time
    ExitDate  time.Time
    Status    string
}

func (e *PastureEntry) New(create *CreatePastureEntry) *PastureEntry {
    e = &PastureEntry{
        Id: uuid.New().String(),
        AnimalId: create.AnimalId,
        PastureId: create.PastureId,
        EntryDate: create.EntryDate,
        ExitDate: create.ExitDate,
        Status: create.Status,
    }
    return e
}

type CreatePastureEntry struct {
    AnimalId  string
    PastureId string
    EntryDate time.Time
    ExitDate  time.Time
    Status    string
}

