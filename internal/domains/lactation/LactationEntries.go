package lactation

import "time"

type LactationDB struct {
	Id                string     `db:"id"`
	AnimalId          string     `db:"animal_id"`
	AnimalTag         string     `db:"animal_tag"`
	AnimalName        string     `db:"animal_name"`
	AnimalOrder       int        `db:"animal_order"`
	CalfId            *string    `db:"calf_id"`
	CalfBirthDate     *time.Time `db:"calf_birth_date"`
	CalfDeathDate     *time.Time `db:"calf_death_date"`
	CalfTag           *string    `db:"calf_tag"`
	CalfName          *string    `db:"calf_name"`
	CalfSex           *string    `db:"calf_sex"`
	StartDate         time.Time  `db:"start_date"`
	EndDate           *time.Time `db:"end_date"`
	LacPeriod         *float64   `db:"lac_period"`
	AverageProduction *float64   `db:"avg_production"`
	TotalProduction   *float64   `db:"total_production"`
	LacInterval       *int       `db:"lac_interval"`
	Peak              *float64   `db:"peak"`
	Observation       *string    `db:"observation"`
	CreatedAt         time.Time  `db:"created_at"`
	DeletedAt         *time.Time `db:"deleted_at"`
	UserId            string     `db:"user_id"`
}

type LactationHistFilter struct {
	Animals              *[]string  `json:"animals" db:"animal_id"`
	HasEndDate           *bool      `json:"hasEndDate" db:"end_date"`
	HasCalf              *bool      `json:"hasCalf" db:"calf_id"`
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
	Observation          *string    `json:"observation" db:"observation"`
}

type AnimalDB struct {
	Id             string     `db:"id"`
	Tag            *string    `db:"tag"`
	Name           *string     `db:"name"`
	BirthDate      *time.Time `db:"birth_date"`
	TagOrder       int        `db:"tag_order"`
	CalfId         *string    `db:"calf_id"`
	CalfTag        *string    `db:"calf_tag"`
	CalfName       *string    `db:"calf_name"`
	CalfSex        *string    `db:"calf_sex"`
	CalfBirthDate  *time.Time `db:"calf_birth_date"`
	CalfDeathDate  *time.Time `db:"calf_death_date"`
	LacId          *string    `db:"lac_id"`
	LacStart       *time.Time `db:"lac_start"`
	LacEnd         *time.Time `db:"lac_end"`
	LacPeriod      *float64   `db:"lac_period"`
	LacAverage     *float64   `db:"lac_average"`
	LacTotal       *float64   `db:"lac_total"`
	LacInterval    *int       `db:"lac_interval"`
	LacPeak        *float64   `db:"lac_peak"`
	LacObservation *string    `db:"lac_observation"`
	IsLactating    bool       `db:"is_lactating"`
	CreatedAt      time.Time  `db:"created_at"`
}

type AnimalFilter struct {
	Animals              *[]string  `json:"animals" db:"animal_id"`
	IsLactating          *bool      `json:"isLactating" db:"is_lactating"`
	HasLactation         *bool      `json:"hasLactation" db:"lac_start"`
	HasCalf              *bool      `json:"hasCalf" db:"calf_id"`
	MinCalfBirthDate     *time.Time `json:"minCalfBirthDate" db:"calf_birth_date"`
	MaxCalfBirthDate     *time.Time `json:"maxCalfBirthDate" db:"calf_birth_date"`
	MinStartDate         *time.Time `json:"minStartDate" db:"lac_start"`
	MaxStartDate         *time.Time `json:"maxStartDate" db:"lac_start"`
	MinEndDate           *time.Time `json:"minEndDate" db:"lac_end"`
	MaxEndDate           *time.Time `json:"maxEndDate" db:"lac_end"`
	MinLacPeriod         *float64   `json:"minLacPeriod" db:"lac_period"`
	MaxLacPeriod         *float64   `json:"maxLacPeriod" db:"lac_period"`
	MinAverageProduction *float64   `json:"minAverageProduction" db:"lac_average"`
	MaxAverageProduction *float64   `json:"maxAverageProduction" db:"lac_average"`
	MinTotalProduction   *float64   `json:"minTotalProduction" db:"lac_total"`
	MaxTotalProduction   *float64   `json:"maxTotalProduction" db:"lac_total"`
	MinLacInterval       *int       `json:"minLacInterval" db:"lac_interval"`
	MaxLacInterval       *int       `json:"maxLacInterval" db:"lac_interval"`
	MinPeak              *float64   `json:"minPeak" db:"lac_peak"`
	MaxPeak              *float64   `json:"maxPeak" db:"lac_peak"`
	Observation          *string    `json:"observation" db:"lac_observation"`
}

type LactationHistFoot struct {
	TotalLacs         int      `json:"totalLacs" db:"total_lacs"`
	AveragePeriod     *float64 `json:"averagePeriod" db:"avg_lac_period"`
	AverageProduction *float64 `json:"averageProduction" db:"avg_production"`
	AverageTotal      *float64 `json:"averageTotal" db:"avg_total_production"`
	AverageInterval   *float64 `json:"averageInterval" db:"avg_lac_interval"`
	AveragePeak       *float64 `json:"averagePeak" db:"avg_peak"`
}

type LactationHistSave struct {
	Id              *string    `json:"id" db:"id"`
	AnimalId        string     `json:"animalId" db:"animal_id"`
	CalfId          *string    `json:"calfId" db:"calf_id"`
	PastureId       *string    `json:"pastureId" db:"pasture_id"`
	StartDate       time.Time  `json:"startDate" db:"start_date"`
	EndDate         *time.Time `json:"endDate" db:"end_date"`
	Observation     *string    `json:"observation" db:"observation"`
	Overwrite       bool       `json:"overwrite"`
	TransferPasture bool       `json:"transferPasture"`
	UserId          string     `json:"-" db:"user_id"`
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

type DairyTypes struct {
	Lactating int `json:"lactating" db:"lactating"`
	Dry       int `json:"dry" db:"dry"`
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
