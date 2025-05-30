package slaughterEntry

import "time"

type SlaughterEntry struct {
	Id             string     `json:"id" db:"id"`
	AnimalId       *string    `json:"animalId" db:"animal_id"`
	AnimalName     *string    `json:"animalName" db:"animal_name"`
	AnimalNumber   *string    `json:"animalNumber" db:"animal_number"`
	AnimalOrder    *string    `json:"animalOrder" db:"animal_order"`
	AnimalBirth    *time.Time `json:"animalBirth" db:"animal_birth"`
	GroupId        string     `json:"groupId" db:"group_id"`
	GroupDate      time.Time  `json:"groupDate" db:"group_date"`
	Slaughterhouse string     `json:"slaughterhouse" db:"slaughterhouse"`
	Weight         float64    `json:"weight" db:"weight"`
	DeadWeight     float64    `json:"deadWeight" db:"dead_weight"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt      *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId         string     `json:"userId" db:"user_id"`
}

type SlaughterEntrySave struct {
	Id               string     `json:"id" db:"id"`
	AnimalId         *string    `json:"animalId" db:"animal_id"`
	GroupId          string     `json:"groupId" db:"group_id"`
	SlaughterhouseId string     `json:"slaughterhouseId" db:"slaughterhouse_id"`
	Weight           float64    `json:"weight" db:"weight"`
	DeadWeight       float64    `json:"deadWeight" db:"dead_weight"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt        *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId           string     `json:"userId" db:"user_id"`
}

type SlaughterEntryFilter struct {
	AnimalId       *[]string  `json:"animalId" db:"animal_id"`
	MinAnimalBirth *time.Time `json:"minAnimalBirth" db:"animal_birth"`
	MaxAnimalBirth *time.Time `json:"maxAnimalBirth" db:"animal_birth"`
	MinGroupDate   *time.Time `json:"minGroupDate" db:"group_date"`
	MaxGroupDate   *time.Time `json:"maxGroupDate" db:"group_date"`
	MinWeight      *float64   `json:"minWeight" db:"weight"`
	MaxWeight      *float64   `json:"maxWeight" db:"weight"`
	MinDeadWeight  *float64   `json:"minDeadWeight" db:"dead_weight"`
	MaxDeadWeight  *float64   `json:"maxDeadWeight" db:"dead_weight"`
}
