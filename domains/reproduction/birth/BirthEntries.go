package birth

import "time"

type BirthEntry struct {
	Id            string     `json:"id" db:"id"`
	AnimalId      string     `json:"animalId" db:"animal_id"`
	AnimalName    string     `json:"animalName" db:"animal_name"`
	AnimalNumber  string     `json:"animalNumber" db:"animal_number"`
	CalfId        string     `json:"calfId" db:"calf_id"`
	CalfBirthDate string     `json:"calfBirthDate" db:"calf_birth_date"`
	CalfSex       string     `json:"calfSex" db:"calf_sex"`
	CalfFather    *string    `json:"calfFather" db:"calf_father"`
	Observation   *string    `json:"observation" db:"observation"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        string     `json:"userId" db:"user_id"`
}

type BirthEntrySave struct {
	Id            string     `json:"id" db:"id"`
	AnimalId      string     `json:"animalId" db:"animal_id"`
	CalfId        string     `json:"calfId" db:"calf_id"`
	Observation   *string    `json:"observation" db:"observation"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        string     `json:"userId" db:"user_id"`
}
