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
	Total         int    `json:"total" db:"total"`
}

type AnimalsByAge struct {
	AgeCategory  string    `json:"ageCategory" db:"age_category"`
	MinBirthDate time.Time `json:"minBirthDate" db:"min_birth_date"`
	MaxBirthDate time.Time `json:"maxBirthDate" db:"max_birth_date"`
	Male         int       `json:"male" db:"male"`
	Female       int       `json:"female" db:"female"`
}

type TotalByYear struct {
	Year         int `json:"year" db:"year"`
	TotalAnimals int `json:"totalAnimals" db:"total_animals"`
}

type AnimalByType struct {
	BeefCattle          int `json:"beefCattle" db:"beef_cattle"`
	DairyCattle         int `json:"dairyCattle" db:"dairy_cattle"`
	ReproductionAnimals int `json:"reproductionAnimals" db:"reproduction_animals"`
	Offspring           int `json:"offspring" db:"offspring"`
}

type AnimalEntriesFilter struct {
	AnimalsDashboardFilter
	MinEntryDate *time.Time `json:"minEntryDate" db:"min_entry_date"`
	MaxExitDate  *time.Time `json:"maxExitDate" db:"max_exit_date"`
}

type AnimalsDashboardFilter struct {
	IsFiltered   bool       `json:"isFiltered"`
	FarmId       *string    `json:"farms" db:"farm_id"`
	PastureId    *string    `json:"pastureId" db:"pasture_id"`
	MinBirthDate *time.Time `json:"minBirthDate" db:"birth_date"`
	MaxBirthDate *time.Time `json:"maxBirthDate" db:"birth_date"`
	AnimalType   *string    `json:"animalType" db:"type"`
	IsActive     *bool      `json:"isActive" db:"is_active"`
}
