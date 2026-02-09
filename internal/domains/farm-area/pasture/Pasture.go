package pasture

import "time"

type PastureDB struct {
	Id            string     `json:"id" db:"id"`
	BullId        *string    `json:"bullId" db:"bull_id"`
	BullTag       *string    `json:"bullTag" db:"bull_tag"`
	BullName      *string    `json:"bullName" db:"bull_name"`
	Name          string     `json:"name" db:"name"`
	FarmId        string     `json:"farmId" db:"farm_id"`
	FarmName      string     `json:"farmName" db:"farm_name"`
	PastureSize   int        `json:"pastureSize" db:"pasture_size"`
	AnimalsNumber int        `json:"animalsNumber" db:"animals_number"`
	CreatedAt     time.Time  `json:"-" db:"created_at"`
	DeletedAt     *time.Time `json:"-" db:"deleted_at"`
}

type PastureFilter struct {
	BullId    *string  `schema:"bullId" db:"bull_id"`
	Name      *string  `schema:"name" db:"name"`
	FarmId    *string  `schema:"farmId" db:"farm_id"`
	MinSize   *float64 `schema:"minPastureSize" db:"pasture_size"`
	MaxSize   *float64 `schema:"maxPastureSize" db:"pasture_size"`
	MinNumber *int     `schema:"minAnimalsNumber" db:"animals_number"`
	MaxNumber *int     `schema:"maxAnimalsNumber" db:"animals_number"`
}

type PastureAnimal struct {
	Id         string     `json:"id" db:"id"`
	Name       *string    `json:"name" db:"name"`
	Tag        *string    `json:"tag" db:"tag"`
	Sex        string     `json:"sex" db:"sex"`
	BirthDate  *time.Time `json:"birthDate" db:"birth_date"`
	AnimalType string     `json:"animalType" db:"animal_type"`
	UserId     string     `json:"userId" db:"user_id"`
}

type PastureSave struct {
	Id        string     `json:"id" db:"id"`
	BullId    string     `json:"bullId" db:"bull_id"`
	Name      string     `json:"name" db:"name"`
	FarmId    string     `json:"farmId" db:"farm_id"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId    string     `json:"userId" db:"user_id"`
}
