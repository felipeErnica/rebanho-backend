package lactation

import "time"

type Lactation struct {
	Id                string     `json:"id" db:"id"`
	CowId             string     `json:"cowId" db:"cow_id"`
	CowName           string     `json:"cowName" db:"cow_name"`
	CowNumber         string     `json:"cowNumber" db:"cow_number"`
	CowPasture        string     `json:"cowPasture" db:"cow_pasture"`
	CowOrder          int        `json:"cowOrder" db:"cow_order"`
	CalfId            string     `json:"calfId" db:"calf_id"`
	CalfBirthDate     time.Time  `json:"calfBirthDate" db:"calf_birth_date"`
	CalfSex           string     `json:"calfSex" db:"calf_sex"`
	CalfFather        string     `json:"calfFather" db:"calf_father"`
	StartDate         time.Time  `json:"startDate" db:"start_date"`
	EndDate           *time.Time `json:"endDate" db:"end_date"`
	ProductionPeriod  uint       `json:"productionPeriod" db:"production_period"`
	ProductionTotal   float32    `json:"productionTotal" db:"production_total"`
	AverageProduction float32    `json:"averageProduction" db:"average_production"`
	PeakProduction    float32    `json:"peakProduction" db:"peak_production"`
	Isr               float32    `json:"isr" db:"isr"`
	Observation       *string    `json:"observation" db:"observation"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt         *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId            string     `json:"userId" db:"user_id"`
}

type LactationSave struct {
	Id                string     `json:"id" db:"id"`
	CowId             string     `json:"cowId" db:"cow_id"`
	CalfId            string     `json:"calfId" db:"calf_id"`
	StartDate         time.Time  `json:"startDate" db:"start_date"`
	EndDate           *time.Time `json:"endDate" db:"end_date"`
	ProductionPeriod  uint       `json:"productionPeriod" db:"production_period"`
	ProductionTotal   float32    `json:"productionTotal" db:"production_total"`
	AverageProduction float32    `json:"averageProduction" db:"average_production"`
	PeakProduction    float32    `json:"peakProduction" db:"peak_production"`
	Isr               float32    `json:"isr" db:"isr"`
	Observation       *string    `json:"observation" db:"observation"`
	CreatedAt         time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt         *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId            string     `json:"userId" db:"user_id"`
}

type LactationFilter struct {
	CowId                *[]string  `json:"cowId" db:"cow_id"`
	CowPasture           *[]string  `json:"cowPasture" db:"cow_pasture"`
	MinCalfBirthDate     *time.Time `json:"minCalfBirthDate" db:"calf_birth_date"`
	MaxCalfBirthDate     *time.Time `json:"maxCalfBirthDate" db:"calf_birth_date"`
	CalfSex              *string    `json:"calfSex" db:"calf_sex"`
	CalfFather           *[]string  `json:"calfFather" db:"calf_father"`
	MinStartDate         *time.Time `json:"minStartDate" db:"start_date"`
	MaxStartDate         *time.Time `json:"maxStartDate" db:"start_date"`
	MinEndDate           *time.Time `json:"minEndDate" db:"end_date"`
	MaxEndDate           *time.Time `json:"maxEndDate" db:"end_date"`
	MinProductionPeriod  uint       `json:"minProductionPeriod" db:"production_period"`
	MaxProductionPeriod  uint       `json:"maxProductionPeriod" db:"production_period"`
	MinProductionTotal   float32    `json:"minProductionTotal" db:"production_total"`
	MaxProductionTotal   float32    `json:"maxProductionTotal" db:"production_total"`
	MinAverageProduction float32    `json:"minAverageProduction" db:"average_production"`
	MaxAverageProduction float32    `json:"maxAverageProduction" db:"average_production"`
	MinPeakProduction    float32    `json:"minPeakProduction" db:"peak_production"`
	MaxPeakProduction    float32    `json:"maxPeakProduction" db:"peak_production"`
	MinIsr               float32    `json:"minIsr" db:"isr"`
	MaxIsr               float32    `json:"maxIsr" db:"isr"`
}
