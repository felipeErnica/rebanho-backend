package insemination

import (
	"time"
)

type InseminationEntry struct {
	Id               string     `json:"id,omitempty" db:"id"`
	AnimalId         string     `json:"animalId,omitempty" db:"animal_id"`
	AnimalOrder      string     `json:"-" db:"animal_order"`
    Name             string     `json:"-" db:"name"`
	AnimalName       string     `json:"animalName" db:"animal_name"`
	GroupId          string     `json:"groupId,omitempty" db:"group_id"`
	InseminationDate *time.Time `json:"inseminationDate,omitempty" db:"insemination_date"`
	BullName         string     `json:"bullName,omitempty" db:"bull_name"`
	Observation      *string    `json:"observation" db:"observation"`
	Status           string     `json:"status" db:"status"`
	LossId           *string    `json:"lossId" db:"loss_id"`
	CalfId           *string    `json:"calfId" db:"calf_id"`
	CreatedAt        time.Time  `json:"-" db:"created_at"`
}

type EntryFooter struct {
	Totals    int     `json:"totals" db:"totals"`
	BirthRate float64 `json:"birthRate" db:"birth_rate"`
}

type InseminationEntryFilter struct {
	IsFiltered          bool       `json:"isFiltered"`
	Animals             *[]string  `json:"animals" db:"animal_id"`
	Groups              *[]string  `json:"groups" db:"group_id"`
	MinInseminationDate *time.Time `json:"minInseminationDate" db:"insemination_date" table:"g"`
	MaxInseminationDate *time.Time `json:"maxInseminationDate" db:"insemination_date" table:"g"`
	Bulls               *[]string  `json:"bulls" db:"bull_id" table:"g"`
	Status              *string    `json:"status" db:"status"`
}

type InseminationGroup struct {
	Id               string    `json:"id" db:"id"`
	BullId           string    `json:"bullId,omitempty" db:"bull_id"`
	BullName         string    `json:"bullName" db:"bull_name"`
	InseminationDate time.Time `json:"inseminationDate" db:"insemination_date"`
	CowNumber        float64   `json:"cowNumber" db:"cow_number"`
	BirthRate        float64   `json:"birthRate" db:"birth_rate"`
	ComparisonRate   float64   `json:"comparisonRate" db:"comparison_rate"`
}

type GroupFooter struct {
	Totals           int     `json:"totals" db:"totals"`
	AverageBirthRate float64 `json:"averageBirthRate" db:"average_birth_rate"`
}

type InseminationEntrySave struct {
	Id          string  `json:"id" db:"id"`
	AnimalId    string  `json:"animal" db:"animal_id"`
	GroupId     string  `json:"groupId" db:"group_id"`
	Observation *string `json:"observation" db:"observation"`
	Status      string  `json:"status" db:"status"`
	LossId      *string `json:"loss" db:"loss_id"`
	CalfId      *string `json:"calf" db:"calf_id"`
}

type InseminationBulls struct {
	BullName       string  `json:"bullName" db:"bull_name"`
	Total          int     `json:"total" db:"total"`
	BirthRate      float64 `json:"birthRate" db:"birth_rate"`
	ComparisonRate float64 `json:"comparisonRate" db:"comparison_rate"`
}

type PregnantsNumber struct {
	PregnantNumber int `json:"pregnantsNumber" db:"pregnants_number"`
}

type InseminationHist struct {
	DateMonth time.Time `json:"dateMonth" db:"date_month"`
	Total     int       `json:"total" db:"total"`
	BirthRate float64   `json:"birthRate" db:"birth_rate"`
}

type BirthRateHist struct {
	DateMonth time.Time `json:"dateMonth" db:"date_month"`
	BirthRate float64   `json:"birthRate" db:"birth_rate"`
}

type BirthRateStats struct {
	Hist    []BirthRateHist `json:"hist"`
	Current float64         `json:"current"`
	Trend   float64         `json:"trend"`
}

type PregnantStats struct {
	PregnantNumber int `json:"pregnantNumber" db:"pregnant_number"`
}
