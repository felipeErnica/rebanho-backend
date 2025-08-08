package birth

import (
	"time"

	"github.com/google/uuid"
)

type BirthEntry struct {
	Id            uuid.UUID  `json:"id" db:"id"`
	MotherId      uuid.UUID  `json:"motherId" db:"mother_id"`
	MotherName    string     `json:"motherName" db:"mother_name" table:"m"`
	MotherOrder   int        `json:"motherOrder" db:"mother_order" table:"m"`
	CalfId        string     `json:"calfId" db:"calf_id"`
	CalfName      string     `json:"calfName" db:"calf_name"`
	CalfBirthDate time.Time  `json:"calfBirthDate" db:"calf_birth_date"`
	CalfSex       string     `json:"calfSex" db:"calf_sex"`
	CalfFatherId  *uuid.UUID `json:"calfFatherId" db:"calf_father_id"`
	CalfFather    *string    `json:"calfFather" db:"calf_father"`
	BirthInterval *int       `json:"birthInterval" db:"birth_interval"`
	Observation   *string    `json:"observation" db:"observation"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        uuid.UUID  `json:"userId" db:"user_id"`
}

type BirthEntryFilter struct {
	IsFiltered       bool       `json:"isFiltered"`
	Mothers          *[]string  `json:"mothers" db:"mother_id"`
	MinBirthDate     *time.Time `json:"minBirthDate" db:"birth_date" table:"c" `
	MaxBirthDate     *time.Time `json:"maxBirthDate" db:"birth_date" table:"c"`
	Sex              *string    `json:"sex" db:"sex" table:"c"`
	Fathers          *[]string  `json:"calfFatherId" db:"father_id" table:"c"`
	MinBirthInterval *int       `json:"maxBirthInterval" db:"birth_interval"`
	MaxBirthInterval *int       `json:"minBirthInterval" db:"birth_interval"`
}

type BirthEntrySave struct {
	Id            string     `json:"id" db:"id"`
	MotherId      string     `json:"motherId" db:"mother_id"`
	CalfId        string     `json:"calfId" db:"calf_id"`
	BirthInterval *int       `json:"birthInterval" db:"birth_interval"`
	Observation   *string    `json:"observation" db:"observation"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        string     `json:"userId" db:"user_id"`
}

type BirthFooter struct {
	Total           int     `json:"total"`
	IntervalAverage float64 `json:"intervalAverage"`
}

type TotalBirthsBySex struct {
	BirthMonth time.Time `json:"birthMonth" db:"birth_month"`
	Males      int       `json:"males" db:"males"`
	Females    int       `json:"females" db:"females"`
}

type BirthsByDate struct {
	Date       time.Time `json:"date" db:"date"`
	BirthTotal int       `json:"birthTotal" db:"birth_total"`
	DeathTotal int       `json:"deathTotal" db:"death_total"`
}

type BirthIntervalHist struct {
	Month           time.Time `json:"month" db:"month"`
	IntervalAverage float64   `json:"intervalAverage" db:"interval_average"`
}

type DeathIndexHist struct {
	Month      time.Time `json:"month" db:"date_month"`
	DeathIndex float64   `json:"deathIndex" db:"death_index"`
}

type LossHist struct {
	Month  time.Time `json:"month" db:"month"`
	Losses int       `json:"losses" db:"losses"`
}

type IntervalAnimal struct {
	AnimalName      string  `json:"animalName" db:"animal_name"`
	BirthNumbers    int     `json:"birthNumbers" db:"birth_numbers"`
	IntervalAverage float64 `json:"intervalAverage" db:"interval_average"`
	AverageRate     float64 `json:"averageRate" db:"average_rate"`
}

type BirthStats struct {
	CurrentInterval     float64             `json:"currentInterval"`
	IntervalTrend       float64             `json:"intervalTrend"`
	IntervalHist        []BirthIntervalHist `json:"intervalHist"`
	DeathIndex          float64             `json:"deathIndex"`
	DeathTrend          float64             `json:"deathTrend"`
	DeathIndexHist      []DeathIndexHist    `json:"deathIndexHist"`
	CurrentBirthNumbers int                 `json:"currentBirthNumbers"`
	BirthNumbersTrend   int                 `json:"birthNumbersTrend"`
	PregnantsNumber     int                 `json:"pregnantsNumber"`
	Losses              int                 `json:"losses"`
	LossTrend           int                 `json:"lossTrend"`
	LossHist            []LossHist          `json:"lossHist"`
	BirthHistory        []BirthsByDate      `json:"birthHistory" db:"birth_history"`
}

type IntervalStats struct {
	CurrentInterval   float64             `json:"currentInterval"`
	IntervalTrend     float64             `json:"intervalTrend"`
	BirthIntervalHist []BirthIntervalHist `json:"intervalHist"`
}

type DeathStats struct {
	DeathIndex      float64          `json:"currentDeathIndex"`
	DeathIndexTrend float64          `json:"deathTrend"`
	DeathIndexHist  []DeathIndexHist `json:"deathIndexHist"`
}

type LossStats struct {
	LossNumbers int        `json:"lossNumbers"`
	LossTrend   int        `json:"lossTrend"`
	LossHist    []LossHist `json:"lossHist"`
}

type CurrentStats struct {
	CurrentBirthNumbers int            `json:"currentBirthNumbers"`
	BirthNumbersTrend   int            `json:"birthNumbersTrend"`
	BirthHistory        []BirthsByDate `json:"birthHistory" db:"birth_history"`
}
