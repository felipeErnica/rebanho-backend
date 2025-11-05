package weight

import "time"

type WeightEntry struct {
	Id              string     `json:"id" db:"id"`
	AnimalId        string     `json:"animalId" db:"animal_id"`
	AnimalName      string     `json:"-" db:"animal_name"`
	AnimalInfo      string     `json:"animalInfo" db:"animal_info"`
	FatherName      *string    `json:"fatherName" db:"father_name"`
	MotherName      *string    `json:"motherName" db:"mother_name"`
	AnimalOrder     int        `json:"-" db:"animal_order"`
	AnimalBirthDate *time.Time `json:"-" db:"birth_date"`
	EntryDate       time.Time  `json:"entryDate" db:"entry_date"`
	Weight          float64    `json:"weight" db:"weight"`
	WeightVariation float64    `json:"weightVariation" db:"weight_variation"`
	WeightGain      float64    `json:"weightGain" db:"weight_gain"`
	CreatedAt       time.Time  `json:"-" db:"created_at"`
	DeletedAt       *time.Time `json:"-" db:"deleted_at"`
	UserId          string     `json:"-" db:"user_id"`
}

type WeightFilter struct {
	IsFiltered   bool       `json:"isFiltered"`
	Animals      *[]string  `json:"animals" db:"animal_id"`
	Fathers      *[]string  `json:"fathers" db:"father_id" table:"a"`
	Mothers      *[]string  `json:"mothers" db:"mother_id" table:"a"`
	MinEntryDate *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate *time.Time `json:"maxEntryDate" db:"entry_date"`
}

type WeightFoot struct {
	AnimalsNumber int     `json:"animalsNumber" db:"animals_num"`
	AverageWeight float64 `json:"averageWeight" db:"avg_weight"`
	AverageGain   float64 `json:"averageGain" db:"avg_gain"`
}

type WeightGroup struct {
	EntryDate       time.Time `json:"entryDate" db:"entry_date"`
	AnimalsNumber   int       `json:"animalsNumber" db:"animals_number"`
	AverageWeight   float64   `json:"averageWeight" db:"average_weight"`
	WeightVariation float64   `json:"weightVariation" db:"weight_variation"`
	AverageGain     float64   `json:"averageGain" db:"average_gain"`
	GainVariation   float64   `json:"gainVariation" db:"gain_variation"`
}

type AverageWeightGain struct {
	EntryDate   time.Time `json:"entryDate" db:"entry_date"`
	AverageGain float64   `json:"averageGain" db:"average_gain"`
}

type CardWeightGain struct {
	Trend   float64             `json:"trend"`
	Current float64             `json:"current"`
	Hist    []AverageWeightGain `json:"hist"`
}

type AverageWeight struct {
	EntryDate     time.Time `json:"entryDate" db:"entry_date"`
	AverageWeight float64   `json:"averageWeight" db:"average_weight"`
}

type CardWeight struct {
	Trend   float64         `json:"trend"`
	Current float64         `json:"current"`
	Hist    []AverageWeight `json:"hist"`
}

type AnimalRating struct {
	AnimalName     string  `json:"animalName" db:"animal_name"`
	AverageGain    float64 `json:"averageGain" db:"avg_gain"`
	GainTrend      float64 `json:"gainTrend" db:"gain_trend"`
	ChildrenNumber int     `json:"childrenNumber" db:"children_number"`
}
