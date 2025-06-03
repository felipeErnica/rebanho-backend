package animalDashboard

import "time"

type TotalBySex struct {
	TotalAnimals int `json:"totalAnimals" db:"total_animals"`
	TotalMales   int `json:"totalMales" db:"total_males"`
	TotalFemales int `json:"totalFemales" db:"total_females"`
}

type AnimalsByAgeAndFarm struct {
	FarmName      string `json:"farmName" db:"farm_name"`
	NewbornMale   int    `json:"newbornMale" db:"newborn_male"`
	NewbornFemale int    `json:"newbornFemale" db:"newborn_female"`
	BabyMale      int    `json:"babyMale" db:"baby_male"`
	BabyFemale    int    `json:"babyFemale" db:"baby_female"`
	ChildMale     int    `json:"childMale" db:"child_male"`
	ChildFemale   int    `json:"childFemale" db:"child_female"`
	YoungMale     int    `json:"youngMale" db:"young_male"`
	YoungFemale   int    `json:"youngFemale" db:"young_female"`
	AdultMale     int    `json:"adultMale" db:"adult_male"`
	AdultFemale   int    `json:"adultFemale" db:"adult_female"`
	OldMale       int    `json:"oldMale" db:"old_male"`
	OldFemale     int    `json:"oldFemale" db:"old_female"`
	TotalMale     int    `json:"totalMale" db:"total_male"`
	TotalFemale   int    `json:"totalFemale" db:"total_female"`
}

type AnimalsDashboardFilter struct {
	IsFiltered     bool       `json:"isFiltered"`
	MinWeaningDate *time.Time `json:"minWeaningDate" db:"weaning_date"`
	MaxWeaningDate *time.Time `json:"maxWeaningDate" db:"weaning_date"`
	MinBirthDate   *time.Time `json:"minBirthDate" db:"birth_date"`
	MaxBirthDate   *time.Time `json:"maxBirthDate" db:"birth_date"`
	Farms          *[]string  `json:"farms" db:"farm_id"`
}
