package animals

import "time"

type Animal struct {
	Id                   string     `json:"id" db:"id"`
	Name                 *string    `json:"name" db:"name"`
	RingNumber           *string    `json:"ringNumber" db:"ring_number"`
	AnimalOrder          int        `json:"animalOrder" db:"animal_order"`
	WeightBirth          float64    `json:"weightBirth" db:"weight_birth"`
	Sex                  string     `json:"sex" db:"sex"`
	WeaningDate          *time.Time `json:"weaningDate" db:"weaning_date"`
	FatherId             *string    `json:"fatherId" db:"father_id"`
	FatherName           *string    `json:"fatherName" db:"father_name"`
	MotherId             *string    `json:"motherId" db:"mother_id"`
	MotherName           *string    `json:"motherName" db:"mother_name"`
	BirthDate            *time.Time `json:"birthDate" db:"birth_date"`
	DeathDate            *time.Time `json:"deathDate" db:"death_date"`
	PastureId            *string    `json:"pastureId" db:"pasture_id"`
	PastureName          *string    `json:"pastureName" db:"pasture_name"`
	FarmId               *string    `json:"farmId" db:"farm_id"`
	FarmName             *string    `json:"farmName" db:"farm_name"`
	AnimalType           string     `json:"animalType" db:"animal_type"`
	IsBreedingBull       bool       `json:"isBreedingBull" db:"is_breeding_bull"`
	IsInseminationBull   bool       `json:"isInseminationBull" db:"is_insemination_bull"`
	IsTransferBull       bool       `json:"isTransferBull" db:"is_transfer_bull"`
	IsEmbryoDonor        bool       `json:"isEmbryoDonor" db:"is_embryo_donor"`
	IsOutsideAnimal      bool       `json:"isOutsideAnimal" db:"is_outside_animal"`
	AverageProd          *float64   `json:"averageProd" db:"average_prod"`
	AverageProdInterval  *float64   `json:"averageProdInterval" db:"average_prod_interval"`
	AverageBirthInterval *float64   `json:"averageBirthInterval" db:"average_birth_interval"`
	AveragePeak          *float64   `json:"averagePeak" db:"average_peak"`
	Observation          *string    `json:"observation" db:"observation"`
	UserId               string     `json:"-" db:"user_id"`
}

type AnimalSave struct {
	Id                 *string    `json:"id" db:"id"`
	Name               *string    `json:"name" db:"name"`
	WeightBirth        *float64   `json:"weightBirth" db:"weight_birth"`
	RingNumber         string     `json:"ringNumber" db:"ring_number"`
	Sex                string     `json:"sex" db:"sex"`
	WeaningDate        *time.Time `json:"weaningDate" db:"weaning_date"`
	FatherId           *string    `json:"fatherId" db:"father_id"`
	MotherId           *string    `json:"motherId" db:"mother_id"`
	BirthDate          time.Time  `json:"birthDate" db:"birth_date"`
	DeathDate          *time.Time `json:"deathDate" db:"death_date"`
	AnimalType         string     `json:"animalType" db:"animal_type"`
	IsBreedingBull     bool       `json:"isBreedingBull" db:"is_breeding_bull"`
	IsInseminationBull bool       `json:"isInseminationBull" db:"is_insemination_bull"`
	IsTransferBull     bool       `json:"isTransferBull" db:"is_transfer_bull"`
	IsEmbryoDonor      bool       `json:"isEmbryoDonor" db:"is_embryo_donor"`
	IsOutsideAnimal    bool       `json:"isOutsideAnimal" db:"is_outside_animal"`
	Observation        *string    `json:"observation" db:"observation"`
	UserId             string     `json:"-" db:"user_id"`
}

type AnimalFilter struct {
	IsFiltered              bool       `json:"isFiltered"`
	Name                    *string    `json:"name" db:"name"`
	Number                  *string    `json:"ringNumber" db:"ring_number"`
	Sex                     *string    `json:"sex" db:"sex"`
	MinWeaningDate          *time.Time `json:"minWeaningDate" db:"weaning_date"`
	MaxWeaningDate          *time.Time `json:"maxWeaningDate" db:"weaning_date"`
	Fathers                 *[]string  `json:"fathers" db:"father_id"`
	Mothers                 *[]string  `json:"mothers" db:"mother_id"`
	MinBirthDate            *time.Time `json:"minBirthDate" db:"birth_date"`
	MaxBirthDate            *time.Time `json:"maxBirthDate" db:"birth_date"`
	MinDeathDate            *time.Time `json:"minDeathDate" db:"death_date"`
	MaxDeathDate            *time.Time `json:"maxDeathDate" db:"death_date"`
	Pastures                *[]string  `json:"pastures" db:"pasture_id"`
	Farms                   *[]string  `json:"farms" db:"farm_id" table:"pastures"`
	Types                   *[]string  `json:"types" db:"type"`
	MinIsr                  *float64   `json:"minIsr" db:"isr"`
	MaxIsr                  *float64   `json:"maxIsr" db:"isr"`
	MinAverageProd          *float64   `json:"minAverageProd" db:"average_prod"`
	MaxAverageProd          *float64   `json:"maxAverageProd" db:"average_prod"`
	MinAverageBirthInterval *float64   `json:"minAverageBirthInterval" db:"average_birth_interval"`
	MaxAverageBirthInterval *float64   `json:"maxAverageBirthInterval" db:"average_birth_interval"`
	MinAveragePeak          *float64   `json:"minAveragePeak" db:"average_peak"`
	MaxAveragePeak          *float64   `json:"maxAveragePeak" db:"average_peak"`
	MinChildrenQuantity     *int       `json:"minChildrenQuantity" db:"children_quantity"`
	MaxChildrenQuantity     *int       `json:"maxChildrenQuantity" db:"children_quantity"`
}

type DeleteAnimalStruct struct {
	Id                string `json:"id"`
	UserId            string `json:"-"`
	CheckLactation    bool   `json:"checkLactation"`
	CheckSlaughter    bool   `json:"checkSlaughter"`
	CheckInsemination bool   `json:"checkInsemination"`
	CheckBreeding     bool   `json:"checkBreeding"`
	CheckTransfer     bool   `json:"checkTransfer"`
}

type TotalBySex struct {
	TotalAnimals int `json:"totalAnimals" db:"total_animals"`
	TotalMales   int `json:"totalMales" db:"total_males"`
	TotalFemales int `json:"totalFemales" db:"total_females"`
}

type AnimalsByAgeAndFarm struct {
	FarmId        *string `json:"farmId" db:"farm_id"`
	FarmName      *string `json:"farmName" db:"farm_name"`
	NewbornMale   int     `json:"newbornMale" db:"newborn_male"`
	NewbornFemale int     `json:"newbornFemale" db:"newborn_female"`
	BabyMale      int     `json:"babyMale" db:"baby_male"`
	BabyFemale    int     `json:"babyFemale" db:"baby_female"`
	ChildMale     int     `json:"childMale" db:"child_male"`
	ChildFemale   int     `json:"childFemale" db:"child_female"`
	YoungMale     int     `json:"youngMale" db:"young_male"`
	YoungFemale   int     `json:"youngFemale" db:"young_female"`
	AdultMale     int     `json:"adultMale" db:"adult_male"`
	AdultFemale   int     `json:"adultFemale" db:"adult_female"`
	OldMale       int     `json:"oldMale" db:"old_male"`
	OldFemale     int     `json:"oldFemale" db:"old_female"`
	TotalMale     int     `json:"totalMale" db:"total_male"`
	TotalFemale   int     `json:"totalFemale" db:"total_female"`
	Total         int     `json:"total" db:"total"`
}

type AnimalsByAge struct {
	AgeCategory  string     `json:"ageCategory" db:"age_category"`
	MinBirthDate *time.Time `json:"minBirthDate" db:"min_birth_date"`
	MaxBirthDate *time.Time `json:"maxBirthDate" db:"max_birth_date"`
	Male         int        `json:"male" db:"male"`
	Female       int        `json:"female" db:"female"`
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
