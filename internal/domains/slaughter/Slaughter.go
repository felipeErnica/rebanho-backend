package slaughter

import "time"

type SlaughterDB struct {
	Id string `db:"id"`

	GroupId       string    `db:"group_id"`
	GroupDate     time.Time `db:"group_date"`
	GroupDiscount float64   `db:"group_discount"`

	ButcherId       string   `db:"butcher_id"`
	ButcherName     string   `db:"butcher_name"`
	ButcherDiscount *float64 `db:"butcher_discount"`

	AnimalId    *string    `db:"animal_id"`
	AnimalTag   *string    `db:"animal_tag"`
	AnimalName  *string    `db:"animal_name"`
	AnimalSex   *string    `db:"animal_sex"`
	AnimalBirth *time.Time `db:"animal_birth"`
	AnimalOrder *string    `db:"animal_order"`

	FatherId   *string `db:"father_id"`
	FatherTag  *string `db:"father_tag"`
	FatherName *string `db:"father_name"`

	MotherId   *string `db:"mother_id"`
	MotherTag  *string `db:"mother_tag"`
	MotherName *string `db:"mother_name"`

	Weight          float64   `db:"weight"`
	DiscountWeight  *float64  `db:"discount_weight"`
	DeadWeight      *float64  `db:"dead_weight"`
	PerformanceRate *float64  `db:"performance_rate"`
	CreatedAt       time.Time `db:"created_at"`
	UserId          string    `db:"user_id"`
}

type SlaughterSave struct {
	Id           *string    `json:"id" db:"id"`
	AnimalId     *string    `json:"animalId" db:"animal_id"`
	GroupId      string     `json:"groupId" db:"group_id"`
	Weight       float64    `json:"weight" db:"weight"`
	DeadWeight   *float64   `json:"deadWeight" db:"dead_weight"`
	UserId       string     `json:"-" db:"user_id"`
	Overwrite    bool       `json:"overwrite"`
	IgnoreDeath  bool       `json:"ignoreDeath"`
}

type SlaughterFilter struct {
	Animals         *[]string  `json:"animals" db:"animal_id"`
	Fathers         *[]string  `json:"fathers" db:"father_id"`
	Mothers         *[]string  `json:"mothers" db:"mother_id"`
	Slaughterhouses *[]string  `json:"slaughterhouses" db:"butcher_id"`
	MinAnimalBirth  *time.Time `json:"minAnimalBirth" db:"birth_date"`
	MaxAnimalBirth  *time.Time `json:"maxAnimalBirth" db:"birth_date"`
	MinEntryDate    *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate    *time.Time `json:"maxEntryDate" db:"entry_date"`
	MinWeight       *float64   `json:"minWeight" db:"weight"`
	MaxWeight       *float64   `json:"maxWeight" db:"weight"`
	MinDeadWeight   *float64   `json:"minDeadWeight" db:"dead_weight"`
	MaxDeadWeight   *float64   `json:"maxDeadWeight" db:"dead_weight"`
}

type SlaughterFoot struct {
	AnimalsNumber     int      `json:"animalsNumber" db:"animals_number"`
	AverageWeight     *float64 `json:"averageWeight" db:"avg_weight"`
	AverageDeadWeight *float64 `json:"averageDeadWeight" db:"avg_dead_weight"`
	AverageRate       *float64 `json:"averageRate" db:"avg_rate"`
}

type SlaughterGroupDB struct {
	EntryDate           time.Time `db:"entry_date"`
	ButcherId           string    `db:"butcher_id"`
	ButcherName         string    `db:"butcher_name"`
	ButcherDiscount     *float64  `db:"butcher_discount"`
	AnimalsNumber       int       `db:"animals_number"`
	AverageWeight       *float64  `db:"avg_weight"`
	AverageDeadWeight   *float64  `db:"avg_dead_weight"`
	WeightVariation     *float64  `db:"weight_variation"`
	DeadWeightVariation *float64  `db:"dead_weight_variation"`
	AverageRate         *float64  `db:"avg_rate"`
	RateVariation       *float64  `db:"rate_variation"`
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
