package lactation

import "time"

type MilkEntry struct {
	Id          string     `json:"id" db:"id"`
	AnimalId    string     `json:"animalId" db:"animal_id"`
	AnimalOrder int        `json:"-" db:"animal_order"`
	AnimalName  string     `json:"animalName" db:"animal_name"`
	EntryDate   time.Time  `json:"entryDate" db:"entry_date"`
	Quantity    float64    `json:"quantity" db:"quantity"`
	CreatedAt   time.Time  `json:"-" db:"created_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
	UserId      string     `json:"-" db:"user_id"`
}

type MilkEntryFilter struct {
	IsFiltered   bool       `json:"isFiltered" db:"is_filtered"`
	Animals      *[]string  `json:"animals" db:"animal_id"`
	MinEntryDate *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate *time.Time `json:"maxEntryDate" db:"entry_date"`
	MinQuantity  *float64   `json:"minQuantity" db:"quantity"`
	MaxQuantity  *float64   `json:"maxQuantity" db:"quantity"`
}

type MilkEntryFoot struct {
	AnimalsNumber int     `json:"animalsNumber" db:"animals_number"`
	TotalMilk     float64 `json:"totalMilk" db:"total_milk"`
	AverageMilk   float64 `json:"averageMilk" db:"avg_milk"`
}

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

type LactationGroup struct {
	EntryDate        time.Time `json:"entryDate" db:"entry_date"`
	AnimalsNumber    int       `json:"animalsNumber" db:"animals_number"`
	TotalMilk        float64   `json:"totalMilk" db:"total_milk"`
	AverageMilk      float64   `json:"averageMilk" db:"avg_milk"`
	NumberDifference int       `json:"numberDifference" db:"number_difference"`
	AverageRate      float64   `json:"averageRate" db:"avg_rate"`
	TotalRate        float64   `json:"totalRate" db:"total_rate"`
}

type LactationGroupFilter struct {
	IsFiltered   bool       `json:"isFiltered" db:"is_filtered"`
	MinEntryDate *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate *time.Time `json:"maxEntryDate" db:"entry_date"`
}

type AnimalsRating struct {
	AnimalName   string  `json:"animalName" db:"animal_name"`
	AvgTotal     float64 `json:"avgTotal" db:"avg_total"`
	AvgPeriod    float64 `json:"avgPeriod" db:"avg_period"`
	AvgProd      float64 `json:"avgProd" db:"avg_prod"`
	AvgInterval  float64 `json:"avgInterval" db:"avg_interval"`
	LacNum       int     `json:"lacNum" db:"lac_num"`
	PeriodRate   float64 `json:"periodRate" db:"period_rate"`
	TotalRate    float64 `json:"totalRate" db:"total_rate"`
	ProdRate     float64 `json:"prodRate" db:"prod_rate"`
	IntervalRate float64 `json:"intervalRate" db:"interval_rate"`
}
