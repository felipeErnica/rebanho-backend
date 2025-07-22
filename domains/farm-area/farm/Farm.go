package farm

import (
	"time"

	"github.com/google/uuid"
)

type Farm struct {
	Id        uuid.UUID  `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId    uuid.UUID  `json:"userId" db:"user_id"`
}

type FarmAnimal struct {
	Id          uuid.UUID  `json:"id" db:"id"`
	Name        *string    `json:"name" db:"name"`
	RingNumber  *string    `json:"ringNumber" db:"ring_number"`
	Sex         string     `json:"sex" db:"sex"`
	FatherId    *uuid.UUID `json:"fatherId" db:"father_id"`
	FatherName  *string    `json:"fatherName" db:"father_name"`
	MotherId    *uuid.UUID `json:"motherId" db:"mother_id"`
	MotherName  *string    `json:"motherName" db:"mother_name"`
	BirthDate   *time.Time `json:"birthDate" db:"birth_date"`
	DeathDate   *time.Time `json:"deathDate" db:"death_date"`
	PastureId   *uuid.UUID `json:"pastureId" db:"pasture_id"`
	PastureName *string    `json:"pastureName" db:"pasture_name"`
	AnimalType  string     `json:"animalType" db:"animal_type"`
	UserId      uuid.UUID  `json:"userId" db:"user_id"`
}

type FarmAnimalFilter struct {
	IsFiltered              bool       `json:"isFiltered"`
	Name                    *string    `json:"name" db:"name"`
	Number                  *string    `json:"ringNumber" db:"ring_number"`
	Sex                     *string    `json:"sex" db:"sex"`
	Fathers                 *[]string  `json:"fathers" db:"father_id"`
	Mothers                 *[]string  `json:"mothers" db:"mother_id"`
	MinBirthDate            *time.Time `json:"minBirthDate" db:"birth_date"`
	MaxBirthDate            *time.Time `json:"maxBirthDate" db:"birth_date"`
	MinDeathDate            *time.Time `json:"minDeathDate" db:"death_date"`
	MaxDeathDate            *time.Time `json:"maxDeathDate" db:"death_date"`
	Pastures                *[]string  `json:"pastures" db:"pasture_id"`
	Types                   *[]string  `json:"types" db:"type"`
}

type FarmFilter struct {
	Name      *string   `json:"name" db:"name"`
	State     *[]string `json:"state" db:"state"`
	City      *string   `json:"city" db:"city"`
	TaxNumber *string   `json:"taxNumber" db:"tax_number"`
	Status    *string   `json:"status" db:"status"`
}
