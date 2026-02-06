package insemination

import (
	"time"
)

type InseminationEntry struct {
	Id               string     `json:"id,omitempty" db:"id"`
	AnimalId         string     `json:"animalId,omitempty" db:"animal_id"`
	AnimalOrder      string     `json:"-" db:"animal_order"`
	AnimalName       string     `json:"-" db:"animal_name"`
	AnimalInfo       string     `json:"animalInfo" db:"animal_info"`
	InseminationDate *time.Time `json:"inseminationDate,omitempty" db:"insemination_date"`
	BullId           string     `json:"bullId,omitempty" db:"bull_id"`
	BullName         string     `json:"bullName,omitempty" db:"bull_name"`
	Observation      *string    `json:"observation" db:"observation"`
	ChildInformation string     `json:"childInformation" db:"child_information"`
	PregnancyStatus  string     `json:"pregnancyStatus" db:"pregnancy_status"`
	BirthStatus      string     `json:"birthStatus" db:"birth_status"`
	CreatedAt        time.Time  `json:"-" db:"created_at"`
}

type InseminationEntrySave struct {
	Id               *string   `json:"id" db:"id"`
	AnimalId         string    `json:"animalId" db:"animal_id"`
	InseminationDate time.Time `json:"inseminationDate" db:"insemination_date"`
	BullId           string    `json:"bullId" db:"bull_id"`
	Observation      *string   `json:"observation" db:"observation"`
	Overwrite        bool      `json:"overwrite"`
	IgnoreWarnings   bool      `json:"ignoreWarnings"`
	UserId           string    `json:"-" db:"user_id"`
}

type InseminationEntryDelete struct {
	Id              string `json:"id" db:"id"`
	ChangeFather    bool   `json:"changeFather"`
	IgnorePregnancy bool   `json:"ignorePregnancy"`
	UserId          string `json:"-" db:"user_id"`
}

type InseminationEntryFilter struct {
	Animals             *[]string  `json:"animals" db:"animal_id"`
	Groups              *[]string  `json:"groups" db:"group_id"`
	MinInseminationDate *time.Time `json:"minInseminationDate" db:"insemination_date"`
	MaxInseminationDate *time.Time `json:"maxInseminationDate" db:"insemination_date"`
	Bulls               *[]string  `json:"bulls" db:"bull_id"`
	BirthStatus         *string    `json:"birthStatus" db:"birth_status"`
	PregnancyStatus     *string    `json:"pregnancyStatus" db:"pregnancy_status"`
}

type InseminationGroup struct {
	BullId                  string    `json:"bullId,omitempty" db:"bull_id"`
	BullName                string    `json:"bullName" db:"bull_name"`
	InseminationDate        time.Time `json:"inseminationDate" db:"insemination_date"`
	CowNumber               float64   `json:"cowNumber" db:"cow_number"`
	BirthRate               float64   `json:"birthRate" db:"birth_rate"`
	PregnancyRate           float64   `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthComparisonRate     float64   `json:"birthComparisonRate" db:"birth_comparison_rate"`
	PregnancyComparisonRate float64   `json:"pregnancyComparisonRate" db:"pregnancy_comparison_rate"`
	UserId                  string    `json:"-" db:"user_id"`
}

type InseminationGroupSave struct {
	InseminationDate    time.Time `json:"inseminationDate" db:"insemination_date"`
	OldInseminationDate time.Time `json:"oldInseminationDate" db:"old_insemination_date"`
	UserId              string    `json:"-" db:"user_id"`
}

type InseminationGroupDelete struct {
	InseminationDate time.Time `json:"inseminationDate" db:"insemination_date"`
	ChangeFather     bool      `json:"changeFather"`
	UserId           string    `json:"-" db:"user_id"`
}

type InseminationFoot struct {
	Totals               int     `json:"totals" db:"totals"`
	AverageBirthRate     float64 `json:"averageBirthRate" db:"average_birth_rate"`
	AveragePregnancyRate float64 `json:"averagePregnancyRate" db:"average_pregnancy_rate"`
}

type LastEntry struct {
	InseminationDate time.Time           `json:"inseminationDate"`
	Entries          []InseminationEntry `json:"entries"`
}

type InseminationBulls struct {
	BullName                string  `json:"bullName" db:"bull_name"`
	Total                   int     `json:"total" db:"total"`
	BirthRate               float64 `json:"birthRate" db:"birth_rate"`
	PregnancyRate           float64 `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthComparisonRate     float64 `json:"birthComparisonRate" db:"birth_comparison_rate"`
	PregnancyComparisonRate float64 `json:"pregnancyComparisonRate" db:"pregnancy_comparison_rate"`
}

type PregnantsNumber struct {
	PregnantNumber int `json:"pregnantsNumber" db:"pregnants_number"`
}

type InseminationHist struct {
	InseminationDate time.Time `json:"inseminationDate" db:"insemination_date"`
	Total            int       `json:"total" db:"total"`
	BirthNumbers     int       `json:"birthNumbers" db:"birth_numbers"`
	PregnancyNumbers int       `json:"pregnancyNumbers" db:"pregnancy_numbers"`
}

type BirthRateHist struct {
	InseminationDate time.Time `json:"inseminationDate" db:"insemination_date"`
	BirthRate        float64   `json:"birthRate" db:"birth_rate"`
	PregnancyRate    float64   `json:"pregnancyRate" db:"pregnancy_rate"`
}

type PregnancyRateHist struct {
	InseminationDate time.Time `json:"inseminationDate" db:"insemination_date"`
	PregnancyRate    float64   `json:"pregnancyRate" db:"pregnancy_rate"`
}

type FutureBirths struct {
	InseminationDate time.Time `json:"birthForecast" db:"birth_forecast"`
	BirthsNumber     int       `json:"birthsNumber" db:"births_number"`
}

type AnimalsHist struct {
	InseminationDate time.Time `json:"inseminationDate" db:"insemination_date"`
	AnimalsNumber    float64   `json:"animalsNumber" db:"animals_number"`
}

type CardStats struct {
	Hist    any     `json:"hist"`
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
}

type BirthRateStats struct {
	Hist    []BirthRateHist `json:"hist"`
	Current float64         `json:"current"`
	Trend   float64         `json:"trend"`
}
