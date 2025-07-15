package dashboard

import "github.com/google/uuid"

type FarmInfo struct {
	FarmId         uuid.UUID `json:"farmId" db:"farm_id"`
	FarmName       string    `json:"farmName" db:"farm_name"`
	PasturesNumber int       `json:"pasturesNumber" db:"pastures_number"`
	AnimalsNumber  int       `json:"animalsNumber" db:"animals_number"`
}

type PastureInfo struct {
	PastureName   string     `json:"pastureName" db:"pasture_name"`
	BullId        *uuid.UUID `json:"bullId" db:"bull_id"`
	BullName      *string    `json:"bullName" db:"bull_name"`
	AnimalsNumber int        `json:"animalsNumber" db:"animals_number"`
}
