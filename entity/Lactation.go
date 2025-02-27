package entity

import (
	"time"

	"github.com/google/uuid"
)

type Lactation struct {
	Id                string  
	AnimalId          string  
    CalfId            string  
	StartDate         time.Time  
	EndDate           time.Time
	ProductionPeriod  uint8
	ProductionTotal   float32
	AvarageProduction float32
	PeakProduction    float32
	Isr               float32
    Observation       string  
    CreatedAt         time.Time  
    DeletedAt         time.Time  
}

func (l *Lactation) New(c *CreateLactation) *Lactation {
    l = &Lactation{
        Id: uuid.NewString(),
        AnimalId: c.AnimalId,
        StartDate: c.StartDate,
        EndDate: c.EndDate,
        ProductionPeriod: c.ProductionPeriod,
        ProductionTotal: c.ProductionTotal,
        AvarageProduction: c.AvarageProduction,
        PeakProduction: c.PeakProduction,
        Isr: c.Isr,
        CreatedAt: time.Now(),
    }
    return l
}

type CreateLactation struct {
	AnimalId          string  
    CalfId            string  
	StartDate         time.Time  
	EndDate           time.Time
	ProductionPeriod  uint8
	ProductionTotal   float32
	AvarageProduction float32
	PeakProduction    float32
	Isr               float32
    Observation       string  
}

type LactationPage struct {
    NextCursor  string
    HasNextPage bool
    List        *[]Lactation
}
