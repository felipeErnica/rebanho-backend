package slaughter

import "time"

type SlaughterEntry struct {
	Id               string     `json:"id" db:"id"`
	AnimalId         *string    `json:"animalId" db:"animal_id"`
	AnimalName       *string    `json:"-" db:"animal_name"`
	AnimalInfo       *string    `json:"animalInfo" db:"animal_info"`
	FatherName       *string    `json:"fatherName" db:"father_name"`
	BirthDate        *time.Time `json:"-" db:"birth_date"`
	MotherName       *string    `json:"motherName" db:"mother_name"`
	AnimalOrder      *string    `json:"-" db:"animal_order"`
	EntryDate        *time.Time `json:"entryDate" db:"entry_date"`
	DiscountRate     *float64   `json:"discountRate" db:"discount_rate"`
	SlaughterhouseId string     `json:"slaughterhouseId" db:"slaughterhouse_id"`
	Slaughterhouse   string     `json:"slaughterhouse" db:"slaughterhouse"`
	Weight           float64    `json:"weight" db:"weight"`
	DiscountWeight   float64    `json:"discountWeight" db:"discount_weight"`
	DeadWeight       float64    `json:"deadWeight" db:"dead_weight"`
	PerformanceRate  float64    `json:"performanceRate" db:"performance_rate"`
	CreatedAt        time.Time  `json:"-" db:"created_at"`
	DeletedAt        *time.Time `json:"-" db:"deleted_at"`
	UserId           string     `json:"-" db:"user_id"`
}

type SlaughterEntryFilter struct {
	IsFiltered      bool       `json:"IsFiltered" db:"is_filtered"`
	Animals         *[]string  `json:"animals" db:"animal_id"`
	Fathers         *[]string  `json:"fathers" db:"father_id" table:"a"`
	Mothers         *[]string  `json:"mothers" db:"mother_id" table:"a"`
	Slaughterhouses *[]string  `json:"slaughterhouses" db:"slaughterhouse_id"`
	MinAnimalBirth  *time.Time `json:"minAnimalBirth" db:"birth_date" table:"a"`
	MaxAnimalBirth  *time.Time `json:"maxAnimalBirth" db:"birth_date" table:"a"`
	MinEntryDate    *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate    *time.Time `json:"maxEntryDate" db:"entry_date"`
	MinWeight       *float64   `json:"minWeight" db:"weight"`
	MaxWeight       *float64   `json:"maxWeight" db:"weight"`
	MinDeadWeight   *float64   `json:"minDeadWeight" db:"dead_weight"`
	MaxDeadWeight   *float64   `json:"maxDeadWeight" db:"dead_weight"`
}

type SlaughterFoot struct {
	AnimalsNumber     int     `json:"animalsNumber" db:"animals_number"`
	AverageWeight     float64 `json:"averageWeight" db:"avg_weight"`
	AverageDeadWeight float64 `json:"averageDeadWeight" db:"avg_dead_weight"`
	AverageRate       float64 `json:"averageRate" db:"avg_rate"`
}

type SlaughterGroup struct {
	EntryDate           time.Time `json:"entryDate" db:"entry_date"`
	Slaughterhouse      string    `json:"slaughterhouse" db:"slaughterhouse"`
	AnimalsNumber       int       `json:"animalsNumber" db:"animals_number"`
	AverageWeight       float64   `json:"averageWeight" db:"avg_weight"`
	AverageDeadWeight   float64   `json:"averageDeadWeight" db:"avg_dead_weight"`
	WeightVariation     float64   `json:"weightVariation" db:"weight_variation"`
	DeadWeightVariation float64   `json:"deadWeightVariation" db:"dead_weight_variation"`
	AverageRate         float64   `json:"averageRate" db:"avg_rate"`
	RateVariation       float64   `json:"rateVariation" db:"rate_variation"`
}

type PerformanceRateHist struct {
	EntryDate       time.Time `json:"entryDate" db:"entry_date"`
	PerformanceRate float64   `json:"performanceRate" db:"performance_rate"`
}

type PerformanceRateCard struct {
	Current float64               `json:"current"`
	Trend   float64               `json:"trend"`
	Hist    []PerformanceRateHist `json:"hist"`
}

type AverageWeightHist struct {
	EntryDate     time.Time `json:"entryDate" db:"entry_date"`
	AverageWeight float64   `json:"averageWeight" db:"avg_weight"`
}

type AverageWeightCard struct {
	Current float64             `json:"current"`
	Trend   float64             `json:"trend"`
	Hist    []AverageWeightHist `json:"hist"`
}

type RateHist struct {
	EntryDate   time.Time `json:"entryDate" db:"entry_date"`
	AverageRate float64   `json:"averageRate" db:"avg_rate"`
}

type WeightHist struct {
	EntryDate  time.Time `json:"entryDate" db:"entry_date"`
	Weight     float64   `json:"weight" db:"avg_weight"`
	DeadWeight float64   `json:"deadWeight" db:"dead_weight"`
}

type TableRatings struct {
	Name             string  `json:"name" db:"name"`
	AnimalNumber     int     `json:"animalsNumber" db:"animals_number"`
	AverageWeight    float64 `json:"averageWeight" db:"avg_weight"`
	WeightComparison float64 `json:"weightComparison" db:"weight_comparison"`
	AverageRate      float64 `json:"performanceRate" db:"avg_rate"`
	RateComparison   float64 `json:"rateComparison" db:"rate_comparison"`
}
