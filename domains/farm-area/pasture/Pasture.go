package pasture

import "time"

type Pasture struct {
	Id        string     `json:"id" db:"id"`
	BullId    string     `json:"bullId" db:"bull_id"`
	BullName  string     `json:"bullName" db:"bull_name"`
	Name      string     `json:"name" db:"name"`
	FarmId    string     `json:"farmId" db:"farm_id"`
	FarmName  string     `json:"farmName" db:"farm_name"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId    string     `json:"userId" db:"user_id"`
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
