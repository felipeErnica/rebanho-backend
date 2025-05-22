package entity

import "time"

type PregnancyTestEntry struct {
	Id            string     `json:"id"`
	GroupId       string     `json:"groupId"`
	GroupDate     string     `json:"groupDate"`
	AnimalId      string     `json:"animalId"`
	AnimalName    string     `json:"animalName"`
	AnimalOrder   int        `json:"animalOrder"`
	AnimalNumber  string     `json:"animalNumber"`
	IsPregnant    bool       `json:"isPregnant"`
	BirthForecast time.Time  `json:"birthForecast"`
	Observation   string     `json:"observation"`
	Status        string     `json:"status"`
	LossId        string     `json:"lossId"`
	CalfId        string     `json:"calfId"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
	UserId        string     `json:"userId"`
}

type PregnancyTestEntrySave struct {
	Id            string     `json:"id"`
	GroupId       string     `json:"groupId"`
	AnimalId      string     `json:"animalId"`
	IsPregnant    bool       `json:"isPregnant"`
	BirthForecast time.Time  `json:"birthForecast"`
	Observation   string     `json:"observation"`
	LossId        string     `json:"lossId"`
	CalfId        string     `json:"calfId"`
	CreatedAt     time.Time  `json:"createdAt"`
	DeletedAt     *time.Time `json:"deletedAt"`
	UserId        string     `json:"userId"`
}
