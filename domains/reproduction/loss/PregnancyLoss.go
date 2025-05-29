package loss

import "time"

type PregnancyLoss struct {
	Id           string     `json:"id" db:"id"`
	AnimalId     string     `json:"animalId" db:"animal_id"`
	AnimalOrder  int        `json:"animalOrder" db:"animal_order"`
	AnimalNumber string     `json:"animalNumber" db:"animal_number"`
	AnimalName   string     `json:"animalName" db:"animal_name"`
	LossType     string     `json:"lossType" db:"loss_type"`
	LossDate     time.Time  `json:"lossDate" db:"loss_date"`
	Observation  *string    `json:"observation" db:"observation"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId       string     `json:"userId" db:"user_id"`
}

type PregnancyLossSave struct {
	Id          string     `json:"id" db:"id"`
	AnimalId    string     `json:"animalId" db:"animal_id"`
	LossType    string     `json:"lossType" db:"loss_type"`
	LossDate    time.Time  `json:"lossDate" db:"loss_date"`
	Observation *string    `json:"observation" db:"observation"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt   *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId      string     `json:"userId" db:"user_id"`
}

type PregnancyLossFilter struct {
	AnimalId    []string     `json:"animalId" db:"animal_id"`
	LossType    []string   `json:"lossType" db:"loss_type"`
	MinLossDate time.Time  `json:"minLossDate" db:"loss_date"`
	MaxLossDate time.Time  `json:"maxLossDate" db:"loss_date"`
}
