package entity

import "time"

type SlaughterEntry struct {
	Id             string     `json:"id"`
	AnimalId       *string    `json:"animalId"`
	AnimalName     *string    `json:"animalName"`
	AnimalNumber   *string    `json:"animalNumber"`
	AnimalBirth    *time.Time `json:"animalBirth"`
	GroupId        string     `json:"groupId"`
	Slaughterhouse string     `json:"slaughterhouse"`
	SlaughterDate  time.Time  `json:"slaughterDate"`
	Weight         float64    `json:"weight"`
	DeadWeight     float64    `json:"deadWeight"`
	CreatedAt      time.Time  `json:"createdAt"`
	DeletedAt      *time.Time `json:"deletedAt"`
	UserId         string     `json:"userId"`
}
