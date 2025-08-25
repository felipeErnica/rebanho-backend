package lactation

import "time"

type YearProductionHist struct {
	EntryDate time.Time `json:"entryDate" db:"entry_date"`
	TotalMilk float64   `json:"totalMilk" db:"total_milk"`
}

type MonthMilkHist struct {
	EntryDate time.Time `json:"entryDate" db:"entry_date"`
	TotalMilk float64   `json:"totalMilk" db:"total_milk"`
}

type AnimalsAverageHist struct {
	EntryDate     time.Time `json:"entryDate" db:"entry_date"`
	AnimalsNumber float64   `json:"animalsNumber" db:"animals_number"`
}

type CardContainer struct {
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
	Hist    any     `json:"hist"`
}

type MilkProductionHist struct {
	EntryDate     time.Time `json:"entryDate" db:"entry_date"`
	AnimalsNumber float64   `json:"animalsNumber" db:"animals_number"`
	TotalMilk     float64   `json:"totalMilk" db:"total_milk"`
}

type AnimalsRating struct {
	AnimalName string  `json:"animalName" db:"animal_name"`
	AvgTotal   float64 `json:"avgTotal" db:"avg_total"`
	AvgPeriod  float64 `json:"avgPeriod" db:"avg_period"`
	AvgProd    float64 `json:"avgProd" db:"avg_prod"`
	LacNum     int     `json:"lacNum" db:"lac_num"`
	PeriodRate float64 `json:"periodRate" db:"period_rate"`
	TotalRate  float64 `json:"totalRate" db:"total_rate"`
	ProdRate   float64 `json:"prodRate" db:"prod_rate"`
}
