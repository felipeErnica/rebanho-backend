package entity

import (
	"time"
)

type Lactation struct {
    Id                string        `json:"id"`
    Cow               Cow           `json:"cow"`
    Calf              Calf          `json:"calf"`
    StartDate         time.Time     `json:"start_date"`
    EndDate           *time.Time    `json:"end_date"`
    ProductionPeriod  uint          `json:"production_period"`
    ProductionTotal   float32       `json:"production_total"`
    AverageProduction float32       `json:"average_production"`
    PeakProduction    float32       `json:"peak_production"`
    Isr               float32       `json:"isr"`
    Observation       *string       `json:"observation"`
    CreatedAt         time.Time     `json:"created_at"`
    DeletedAt         *time.Time    `json:"deleted_at"`
}

type LactationShort struct {
    Id                string        `json:"id"`
    AnimalId          string        `json:"animal_id"`
    CalfId            string        `json:"calf_id"`
    StartDate         time.Time     `json:"start_date"`
    EndDate           time.Time     `json:"end_date"`
    ProductionPeriod  uint          `json:"production_period"`
    ProductionTotal   float32       `json:"production_total"`
    AverageProduction float32       `json:"average_production"`
    PeakProduction    float32       `json:"peak_production"`
    Isr               float32       `json:"isr"`
    Observation       string        `json:"observation"`
}
