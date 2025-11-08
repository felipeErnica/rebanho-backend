package lactation

import "time"

type MilkEntry struct {
	Id          string     `json:"id" db:"id"`
	AnimalId    string     `json:"animalId" db:"animal_id"`
	AnimalOrder int        `json:"-" db:"animal_order"`
	AnimalName  string     `json:"-" db:"animal_name"`
	AnimalInfo  string     `json:"animalInfo" db:"animal_info"`
	PastureName string     `json:"pastureName" db:"pasture_name"`
	EntryDate   time.Time  `json:"entryDate" db:"entry_date"`
	Quantity    float64    `json:"quantity" db:"quantity"`
	CreatedAt   time.Time  `json:"-" db:"created_at"`
	DeletedAt   *time.Time `json:"-" db:"deleted_at"`
	UserId      string     `json:"-" db:"user_id"`
}

type MilkEntryFilter struct {
	IsFiltered   bool       `json:"isFiltered" db:"is_filtered"`
	Animals      *[]string  `json:"animals" db:"animal_id"`
	Pastures     *[]string  `json:"pastures" db:"pasture_id" table:"pe"`
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

type LactationHist struct {
	Id                string     `json:"id" db:"id"`
	AnimalId          string     `json:"animalId" db:"animal_id"`
	Name              string     `json:"-" db:"name"`
	AnimalName        string     `json:"animalName" db:"animal_name"`
	AnimalOrder       int        `json:"-" db:"animal_order"`
	CalfId            *string    `json:"calfId" db:"calf_id"`
	CalfBirthDate     *time.Time `json:"-" db:"calf_birth_date"`
	CalfInfo          *string    `json:"calfInfo" db:"calf_info"`
	StartDate         time.Time  `json:"startDate" db:"start_date"`
	EndDate           *time.Time `json:"endDate" db:"end_date"`
	LacPeriod         float64    `json:"lacPeriod" db:"lac_period"`
	AverageProduction float64    `json:"averageProduction" db:"avg_production"`
	TotalProduction   float64    `json:"totalProduction" db:"total_production"`
	LacInterval       *int       `json:"lacInterval" db:"lac_interval"`
	Peak              float64    `json:"peak" db:"peak"`
	Observation       *string    `json:"observation" db:"observation"`
	CreatedAt         time.Time  `json:"-" db:"created_at"`
	DeletedAt         *time.Time `json:"-" db:"deleted_at"`
	UserId            string     `json:"-" db:"user_id"`
}

type LactationHistFilter struct {
	IsFiltered           bool       `json:"isFiltered"`
	Animals              *[]string  `json:"animals" db:"animal_id"`
	MinCalfBirthDate     *time.Time `json:"minCalfBirthDate" db:"calf_birth_date"`
	MaxCalfBirthDate     *time.Time `json:"maxCalfBirthDate" db:"calf_birth_date"`
	MinStartDate         *time.Time `json:"minStartDate" db:"start_date"`
	MaxStartDate         *time.Time `json:"maxStartDate" db:"start_date"`
	MinEndDate           *time.Time `json:"minEndDate" db:"end_date"`
	MaxEndDate           *time.Time `json:"maxEndDate" db:"end_date"`
	MinLacPeriod         *float64   `json:"minLacPeriod" db:"lac_period"`
	MaxLacPeriod         *float64   `json:"maxLacPeriod" db:"lac_period"`
	MinAverageProduction *float64   `json:"minAverageProduction" db:"avg_production"`
	MaxAverageProduction *float64   `json:"maxAverageProduction" db:"avg_production"`
	MinTotalProduction   *float64   `json:"minTotalProduction" db:"total_production"`
	MaxTotalProduction   *float64   `json:"maxTotalProduction" db:"total_production"`
	MinLacInterval       *int       `json:"minLacInterval" db:"lac_interval"`
	MaxLacInterval       *int       `json:"maxLacInterval" db:"lac_interval"`
	MinPeak              *float64   `json:"minPeak" db:"peak"`
	MaxPeak              *float64   `json:"maxPeak" db:"peak"`
}

type LactationHistFoot struct {
	TotalLacs         int     `json:"totalLacs" db:"total_lacs"`
	AveragePeriod     float64 `json:"averagePeriod" db:"avg_lac_period"`
	AverageProduction float64 `json:"averageProduction" db:"avg_production"`
	AverageTotal      float64 `json:"averageTotal" db:"avg_total_production"`
	AverageInterval   float64 `json:"averageInterval" db:"avg_lac_interval"`
	AveragePeak       float64 `json:"averagePeak" db:"avg_peak"`
}

type AddLactationStruct struct {
	AnimalId    string     `json:"animalId"`
	CalfId      *string    `json:"calfId"`
	PastureId   *string    `json:"pastureId"`
	StartDate   time.Time  `json:"startDate"`
	EndDate     *time.Time `json:"endDate"`
	Observation *string    `json:"observation"`
	UserId      string     `json:"-"`
}

type AverageMilkEntry struct {
	EntryDate   time.Time `json:"entryDate" db:"entry_date"`
	AverageMilk float64   `json:"averageMilk" db:"avg_milk"`
}

type TotalMilkEntry struct {
	EntryDate time.Time `json:"entryDate" db:"entry_date"`
	TotalMilk float64   `json:"totalMilk" db:"total_milk"`
}

type AnimalsAverageHist struct {
	EntryDate     time.Time `json:"entryDate" db:"entry_date"`
	AnimalsNumber float64   `json:"animalsNumber" db:"animals_number"`
}

type MilkProductionEntry struct {
	EntryDate     time.Time `json:"entryDate" db:"entry_date"`
	TotalMilk     float64   `json:"totalMilk" db:"total_milk"`
	AnimalsNumber float64   `json:"animalsNumber" db:"animals_number"`
}

type CardContainer struct {
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
	Hist    any     `json:"hist"`
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
	AnimalName    string  `json:"animalName" db:"animal_name"`
	AvgTotal      float64 `json:"avgTotal" db:"avg_total"`
	AvgPeriod     float64 `json:"avgPeriod" db:"avg_period"`
	AvgProd       float64 `json:"avgProd" db:"avg_prod"`
	AvgInterval   float64 `json:"avgInterval" db:"avg_interval"`
	LacNum        int     `json:"lacNum" db:"lac_num"`
	PeriodRate    float64 `json:"periodRate" db:"period_rate"`
	TotalRate     float64 `json:"totalRate" db:"total_rate"`
	ProdRate      float64 `json:"prodRate" db:"prod_rate"`
	IntervalRate  float64 `json:"intervalRate" db:"interval_rate"`
	TotalScore    float64 `json:"-" db:"total_score"`
	IntervalScore float64 `json:"-" db:"interval_score"`
}

type ParentsRating struct {
	ParentName     string  `json:"parentName" db:"parent_name"`
	ChildrenNumber float64 `json:"childrenNumber" db:"children_number"`
	AvgTotal       float64 `json:"avgTotal" db:"avg_total"`
	AvgPeriod      float64 `json:"avgPeriod" db:"avg_period"`
	AvgProd        float64 `json:"avgProd" db:"avg_prod"`
	AvgInterval    float64 `json:"avgInterval" db:"avg_interval"`
	PeriodRate     float64 `json:"periodRate" db:"period_rate"`
	TotalRate      float64 `json:"totalRate" db:"total_rate"`
	ProdRate       float64 `json:"prodRate" db:"prod_rate"`
	IntervalRate   float64 `json:"intervalRate" db:"interval_rate"`
	TotalScore     float64 `json:"-" db:"total_score"`
	IntervalScore  float64 `json:"-" db:"interval_score"`
}
