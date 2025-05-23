package entity

import "time"

type Lactation struct {
	Id                string     `json:"id"`
	CowId             string     `json:"cowId"`
	CowName           string     `json:"cowName"`
	CowNumber         string     `json:"cowNumber"`
	CowPasture        string     `json:"cowPasture"`
	CowOrder          int        `json:"cowOrder"`
	CalfId            string     `json:"calfId"`
	CalfBirthDate     time.Time  `json:"calfBirthDate"`
	CalfSex           string     `json:"calfSex"`
	CalfFather        string     `json:"calfFather"`
	StartDate         time.Time  `json:"startDate"`
	EndDate           *time.Time `json:"endDate"`
	ProductionPeriod  uint       `json:"productionPeriod"`
	ProductionTotal   float32    `json:"productionTotal"`
	AverageProduction float32    `json:"averageProduction"`
	PeakProduction    float32    `json:"peakProduction"`
	Isr               float32    `json:"isr"`
	Observation       *string    `json:"observation"`
	CreatedAt         time.Time  `json:"createdAt"`
	DeletedAt         *time.Time `json:"deletedAt"`
	UserId            string     `json:"userId"`
}
