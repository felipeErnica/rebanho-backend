package loss

import (
	"time"
)

type PregnancyLoss struct {
	Id           string     `json:"id" db:"id"`
	AnimalId     string     `json:"animalId" db:"animal_id"`
	AnimalOrder  int        `json:"animalOrder" db:"animal_order"`
	AnimalNumber string     `json:"animalNumber" db:"animal_number"`
	AnimalName   string     `json:"animalName" db:"animal_name"`
	LossDate     time.Time  `json:"lossDate" db:"loss_date"`
	Observation  *string    `json:"observation" db:"observation"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string     `json:"userId" db:"user_id"`
}

type PregnancyLossSave struct {
	Id          string     `json:"id" db:"id"`
	AnimalId    string     `json:"animalId" db:"animal_id"`
	LossDate    time.Time  `json:"lossDate" db:"loss_date"`
	Observation *string    `json:"observation" db:"observation"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt   *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId      string     `json:"userId" db:"user_id"`
}

type LossFilter struct {
	IsFiltered  bool       `json:"isFiltered"`
	Animals     *[]string  `json:"animals" db:"animal_id"`
	MinLossDate *time.Time `json:"minLossDate" db:"loss_date"`
	MaxLossDate *time.Time `json:"maxLossDate" db:"loss_date"`
}

type LossFooter struct {
	Totals int `json:"totals"`
}

type LossRate struct {
	Trend   float64        `json:"trend"`
	Current float64        `json:"current"`
	Hist    []LossRateHist `json:"hist"`
}

type LossRateHist struct {
	LossDate time.Time `json:"lossDate" db:"loss_date"`
	LossRate float64   `json:"lossRate" db:"loss_rate"`
}

type LossNumbersHist struct {
	LossDate    time.Time `json:"lossDate" db:"loss_date"`
	LossNumbers int       `json:"lossNumbers" db:"loss_numbers"`
}

type MostLossesAnimals struct {
	AnimalName     string  `json:"animalName" db:"animal_name"`
	Losses         int     `json:"losses" db:"losses"`
	LossRate       float64 `json:"lossRate" db:"loss_rate"`
	RateComparison float64 `json:"rateComparison" db:"rate_comparison"`
}
