package animals

import "time"

type AnimalDTO struct {
	Id                   string     `json:"id"`
	Name                 *string    `json:"name"`
	Tag                  *string    `json:"tag"`
	WeightBirth          *float64   `json:"weightBirth"`
	Sex                  string     `json:"sex"`
	BirthDate            *time.Time `json:"birthDate"`
	DeathDate            *time.Time `json:"deathDate"`
	WeaningDate          *time.Time `json:"weaningDate"`
	AnimalType           string     `json:"animalType"`
	Father               *Parent    `json:"father"`
	Mother               *Parent    `json:"mother"`
	Pasture              *Pasture   `json:"pasture"`
	IsBreedingBull       bool       `json:"isBreedingBull"`
	IsInseminationBull   bool       `json:"isInseminationBull"`
	IsTransferBull       bool       `json:"isTransferBull"`
	IsEmbryoDonor        bool       `json:"isEmbryoDonor"`
	IsOutsideAnimal      bool       `json:"isOutsideAnimal"`
	AverageProd          *float64   `json:"averageProd"`
	AverageLacInterval   *float64   `json:"averageLacInterval"`
	AverageBirthInterval *float64   `json:"averageBirthInterval"`
	AveragePeak          *float64   `json:"averagePeak"`
	ChildrenNumber       *int       `json:"childrenNumber"`
	Observation          *string    `json:"observation"`
}

type Parent struct {
	Id   string  `json:"id"`
	Name *string `json:"name"`
	Tag  *string `json:"tag"`
}

type Pasture struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Farm Farm   `json:"farm"`
}

type Farm struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
