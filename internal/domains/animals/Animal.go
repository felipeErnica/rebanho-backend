package animals

import "time"

type AnimalDB struct {
	Id                   string     `db:"id"`
	Name                 *string    `db:"name"`
	Tag                  *string    `db:"tag"`
	AnimalOrder          int        `db:"animal_order"`
	WeightBirth          *float64   `db:"weight_birth"`
	Sex                  string     `db:"sex"`
	BirthDate            *time.Time `db:"birth_date"`
	DeathDate            *time.Time `db:"death_date"`
	WeaningDate          *time.Time `db:"weaning_date"`
	FatherId             *string    `db:"father_id"`
	FatherName           *string    `db:"father_name"`
	FatherTag            *string    `db:"father_tag"`
	MotherId             *string    `db:"mother_id"`
	MotherName           *string    `db:"mother_name"`
	MotherTag            *string    `db:"mother_tag"`
	PastureId            *string    `db:"pasture_id"`
	PastureName          *string    `db:"pasture_name"`
	FarmId               *string    `db:"farm_id"`
	FarmName             *string    `db:"farm_name"`
	AnimalType           string     `db:"animal_type"`
	IsBreedingBull       bool       `db:"is_breeding_bull"`
	IsInseminationBull   bool       `db:"is_insemination_bull"`
	IsTransferBull       bool       `db:"is_transfer_bull"`
	IsEmbryoDonor        bool       `db:"is_embryo_donor"`
	IsOutsideAnimal      bool       `db:"is_outside_animal"`
	AverageProd          *float64   `db:"average_prod"`
	AverageLacInterval   *float64   `db:"average_lac_interval"`
	AverageBirthInterval *float64   `db:"average_birth_interval"`
	AveragePeak          *float64   `db:"average_peak"`
	ChildrenNumber       *int       `db:"children_number"`
	Observation          *string    `db:"observation"`
	IsLactating          *bool      `db:"is_lactating"`
	CreatedAt            time.Time  `db:"created_at"`
	UserId               string     `db:"user_id"`
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
	IgnoreDead           bool       `json:"ignoreDead"`
	Observation          *string    `json:"observation" db:"observation"`
	UserId               string     `json:"-" db:"user_id"`
}

type AnimalFilter struct {
	Name                    *string    `schema:"name" db:"name"`
	Number                  *string    `schema:"ringNumber" db:"ring_number"`
	Sex                     *string    `schema:"sex" db:"sex"`
	MinWeaningDate          *time.Time `schema:"minWeaningDate" db:"weaning_date"`
	MaxWeaningDate          *time.Time `schema:"maxWeaningDate" db:"weaning_date"`
	Fathers                 *[]string  `schema:"fathers" db:"father_id"`
	Mothers                 *[]string  `schema:"mothers" db:"mother_id"`
	MinBirthDate            *time.Time `schema:"minBirthDate" db:"birth_date"`
	MaxBirthDate            *time.Time `schema:"maxBirthDate" db:"birth_date"`
	MinDeathDate            *time.Time `schema:"minDeathDate" db:"death_date"`
	MaxDeathDate            *time.Time `schema:"maxDeathDate" db:"death_date"`
	Pastures                *[]string  `schema:"pastures" db:"pasture_id"`
	Farms                   *[]string  `schema:"farms" db:"farm_id" table:"pastures"`
	Types                   *[]string  `schema:"types" db:"animal_type"`
	MinAverageProd          *float64   `schema:"minAverageProd" db:"average_prod"`
	MaxAverageProd          *float64   `schema:"maxAverageProd" db:"average_prod"`
	MinAverageBirthInterval *float64   `schema:"minAverageBirthInterval" db:"average_birth_interval"`
	MaxAverageBirthInterval *float64   `schema:"maxAverageBirthInterval" db:"average_birth_interval"`
	MinAveragePeak          *float64   `schema:"minAveragePeak" db:"average_peak"`
	MaxAveragePeak          *float64   `schema:"maxAveragePeak" db:"average_peak"`
	MinChildrenNumber       *int       `schema:"minChildrenNumber" db:"children_number"`
	MaxChildrenNumber       *int       `schema:"maxChildrenNumber" db:"children_number"`
	HasDeath                *bool      `schema:"isAlive" db:"death_date"`
	HasName                 *bool      `schema:"hasName" db:"name"`
	IsBreedingBull          *bool      `schema:"isBreedingBull" db:"is_breeding_bull"`
	IsInsemininationBull    *bool      `schema:"isInseminationBull" db:"is_insemination_bull"`
	IsTransferBull          *bool      `schema:"isTransferBull" db:"is_transfer_bull"`
	IsEmbryoDonor           *bool      `schema:"isEmbryoDonor" db:"is_embryo_donor"`
	IsOutsideAnimal         *bool      `schema:"isOutsideAnimal" db:"is_outside_animal"`
	IsLactating             *bool      `schema:"isLactating" db:"is_lactating"`
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
