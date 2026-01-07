package animals

import "time"

type Animal struct {
	Id                   string     `json:"id" db:"id"`
	Name                 *string    `json:"name" db:"name"`
	RingNumber           *string    `json:"ringNumber" db:"ring_number"`
	AnimalOrder          int        `json:"animalOrder" db:"animal_order"`
	WeightBirth          *float64   `json:"weightBirth" db:"weight_birth"`
	Sex                  string     `json:"sex" db:"sex"`
	BirthDate            *time.Time `json:"birthDate" db:"birth_date"`
	DeathDate            *time.Time `json:"deathDate" db:"death_date"`
	WeaningDate          *time.Time `json:"weaningDate" db:"weaning_date"`
	FatherId             *string    `json:"fatherId" db:"father_id"`
	FatherName           *string    `json:"fatherName" db:"father_name"`
	MotherId             *string    `json:"motherId" db:"mother_id"`
	MotherName           *string    `json:"motherName" db:"mother_name"`
	PastureId            *string    `json:"pastureId" db:"pasture_id"`
	PastureName          *string    `json:"pastureName" db:"pasture_name"`
	AnimalType           string     `json:"animalType" db:"animal_type"`
	IsBreedingBull       bool       `json:"isBreedingBull" db:"is_breeding_bull"`
	IsInseminationBull   bool       `json:"isInseminationBull" db:"is_insemination_bull"`
	IsTransferBull       bool       `json:"isTransferBull" db:"is_transfer_bull"`
	IsEmbryoDonor        bool       `json:"isEmbryoDonor" db:"is_embryo_donor"`
	IsOutsideAnimal      bool       `json:"isOutsideAnimal" db:"is_outside_animal"`
	AverageProd          *float64   `json:"averageProd" db:"average_prod"`
	AverageLacInterval   *float64   `json:"averageLacInterval" db:"average_lac_interval"`
	AverageBirthInterval *float64   `json:"averageBirthInterval" db:"average_birth_interval"`
	AveragePeak          *float64   `json:"averagePeak" db:"average_peak"`
	Observation          *string    `json:"observation" db:"observation"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	UserId               string     `json:"-" db:"user_id"`
}

type AnimalFoot struct {
	Total                int      `json:"total" db:"total"`
	AverageProd          *float64 `json:"averageProd" db:"average_prod"`
	AverageBirthInterval *float64 `json:"averageBirthInterval" db:"average_birth_interval"`
	AverageLacInterval   *float64 `json:"averageLacInterval" db:"average_lac_interval"`
	AveragePeak          *float64 `json:"averagePeak" db:"average_peak"`
}

type AnimalSave struct {
	Id                   *string    `json:"id" db:"id"`
	Name                 *string    `json:"name" db:"name"`
	WeightBirth          *float64   `json:"weightBirth" db:"weight_birth"`
	RingNumber           string     `json:"ringNumber" db:"ring_number"`
	Sex                  string     `json:"sex" db:"sex"`
	WeaningDate          *time.Time `json:"weaningDate" db:"weaning_date"`
	FatherId             *string    `json:"fatherId" db:"father_id"`
	MotherId             *string    `json:"motherId" db:"mother_id"`
	BirthDate            time.Time  `json:"birthDate" db:"birth_date"`
	DeathDate            *time.Time `json:"deathDate" db:"death_date"`
	AnimalType           string     `json:"animalType" db:"animal_type"`
	IsBreedingBull       bool       `json:"isBreedingBull" db:"is_breeding_bull"`
	IsInsemininationBull bool       `json:"isInseminationBull" db:"is_insemination_bull"`
	IsTransferBull       bool       `json:"isTransferBull" db:"is_transfer_bull"`
	IsEmbryoDonor        bool       `json:"isEmbryoDonor" db:"is_embryo_donor"`
	IsOutsideAnimal      bool       `json:"isOutsideAnimal" db:"is_outside_animal"`
	Observation          *string    `json:"observation" db:"observation"`
	UserId               string     `json:"-" db:"user_id"`
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
	MinAverageProd          *float64   `json:"minAverageProd" db:"average_prod"`
	MaxAverageProd          *float64   `json:"maxAverageProd" db:"average_prod"`
	MinAverageBirthInterval *float64   `json:"minAverageBirthInterval" db:"average_birth_interval"`
	MaxAverageBirthInterval *float64   `json:"maxAverageBirthInterval" db:"average_birth_interval"`
	MinAveragePeak          *float64   `json:"minAveragePeak" db:"average_peak"`
	MaxAveragePeak          *float64   `json:"maxAveragePeak" db:"average_peak"`
	MinChildrenQuantity     *int       `json:"minChildrenQuantity" db:"children_quantity"`
	MaxChildrenQuantity     *int       `json:"maxChildrenQuantity" db:"children_quantity"`
}

type CardEntry struct {
	Current int     `json:"current"`
	Trend   float64 `json:"trend"`
	Hist    any     `json:"hist"`
}

type AnimalsNumberHist struct {
	EntryDate     time.Time `json:"entryDate" db:"entry_date"`
	AnimalsNumber int       `json:"animalsNumber" db:"animals_number"`
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

type AnimalsByAge struct {
	Category string `json:"category" db:"category"`
	Female   int    `json:"female" db:"female"`
	Male     int    `json:"male" db:"male"`
}

type TotalByYear struct {
	Year         int `json:"year" db:"year"`
	TotalAnimals int `json:"totalAnimals" db:"total_animals"`
}

type AnimalByType struct {
	BeefAnimals         int `json:"beefAnimals" db:"beef_animals"`
	DairyAnimals        int `json:"dairyAnimals" db:"dairy_animals"`
	ReproductionAnimals int `json:"reproductionAnimals" db:"reproduction_animals"`
	Offspring           int `json:"offspring" db:"offspring"`
}
