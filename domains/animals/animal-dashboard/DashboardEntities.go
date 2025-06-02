package animalDashboard

type AnimalsByAgeAndFarm struct {
	FarmName      string `json:"farmName" db:"farm_name"`
	NewBornMale   int    `json:"newBornMale" db:"new_born_male"`
	NewBornFemale int    `json:"newBornFemale" db:"new_born_female"`
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
}
