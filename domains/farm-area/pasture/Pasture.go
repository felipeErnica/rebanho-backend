package pasture

import (
	"time"

	"github.com/google/uuid"
)

type Pasture struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	BullId    *uuid.UUID `json:"bullId" db:"bull_id"`
	BullName  *string    `json:"bullName" db:"bull_name"`
	Name      string     `json:"name" db:"name"`
	FarmId    uuid.UUID  `json:"farmId" db:"farm_id"`
	FarmName  string     `json:"farmName" db:"farm_name"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId    uuid.UUID  `json:"userId" db:"user_id"`
}

type PastureAnimal struct {
	Id         uuid.UUID  `json:"id" db:"id"`
	Name       *string    `json:"name" db:"name"`
	RingNumber *string    `json:"ringNumber" db:"ring_number"`
	Sex        string     `json:"sex" db:"sex"`
	FatherId   *uuid.UUID `json:"fatherId" db:"father_id"`
	FatherName *string    `json:"fatherName" db:"father_name"`
	MotherId   *uuid.UUID `json:"motherId" db:"mother_id"`
	MotherName *string    `json:"motherName" db:"mother_name"`
	BirthDate  *time.Time `json:"birthDate" db:"birth_date"`
	DeathDate  *time.Time `json:"deathDate" db:"death_date"`
	AnimalType string     `json:"animalType" db:"animal_type"`
	UserId     uuid.UUID  `json:"userId" db:"user_id"`
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
