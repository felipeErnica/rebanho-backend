package naturalMating

import "time"

type NaturalMating struct {
	Id           string     `json:"id" db:"id"`
	AnimalId     string     `json:"animalId" db:"animal_id"`
	AnimalNumber string     `json:"animalNumber" db:"animal_number"`
	AnimalName   string     `json:"animalName" db:"animal_name"`
	MatingDate   time.Time  `json:"matingDate" db:"mating_date"`
	BullId       string     `json:"bullId" db:"bull_id"`
	BullName     string     `json:"bullName" db:"bull_name"`
	Observation  *string    `json:"observation" db:"observation"`
	Status       string     `json:"status" db:"status"`
	LossId       *string    `json:"lossId" db:"loss_id"`
	CalfId       *string    `json:"calfId" db:"calf_id"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string     `json:"userId" db:"user_id"`
}

type NaturalMatingSave struct {
	Id          string     `json:"id" db:"id"`
	AnimalId    string     `json:"animalId" db:"animal_id"`
	MatingDate  time.Time  `json:"matingDate" db:"mating_date"`
	BullId      string     `json:"bullId" db:"bull_id"`
	Observation *string    `json:"observation" db:"observation"`
	Status      string     `json:"status" db:"status"`
	LossId      *string    `json:"lossId" db:"loss_id"`
	CalfId      *string    `json:"calfId" db:"calf_id"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt   *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId      string     `json:"userId" db:"user_id"`
}
