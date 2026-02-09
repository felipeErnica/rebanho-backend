package milk

import "time"

type MilkDB struct {
	Id          string     `db:"id"`
	AnimalId    string     `db:"animal_id"`
	AnimalTag   *string    `db:"animal_tag"`
	AnimalName  *string    `db:"animal_name"`
	AnimalOrder int        `db:"animal_order"`
	PastureId   *string    `db:"pasture_id"`
	PastureName *string     `db:"pasture_name"`
	FarmId      *string     `db:"farm_id"`
	FarmName    *string     `db:"farm_name"`
	EntryDate   time.Time  `db:"entry_date"`
	Quantity    float64    `db:"quantity"`
	CreatedAt   time.Time  `db:"created_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	UserId      string     `db:"user_id"`
}

type MilkEntryFilter struct {
	Animals      *[]string  `json:"animals" db:"animal_id"`
	Pastures     *[]string  `json:"pastures" db:"pasture_id" table:"pe"`
	MinEntryDate *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate *time.Time `json:"maxEntryDate" db:"entry_date"`
	MinQuantity  *float64   `json:"minQuantity" db:"quantity"`
	MaxQuantity  *float64   `json:"maxQuantity" db:"quantity"`
}

type MilkEntryFoot struct {
	AnimalsNumber int     `json:"animalsNumber" db:"animals_number"`
	TotalMilk     float64 `json:"totalMilk" db:"total_milk"`
	AverageMilk   float64 `json:"averageMilk" db:"avg_milk"`
}

type MilkEntrySave struct {
	Id              *string   `json:"id" db:"id"`
	AnimalId        string    `json:"animalId" db:"animal_id"`
	PastureId       *string   `json:"pastureId" db:"pasture_id"`
	EntryDate       time.Time `json:"entryDate" db:"entry_date"`
	Quantity        float64   `json:"quantity" db:"quantity"`
	Overwrite       bool      `json:"overwrite"`
	TransferPasture bool      `json:"transferPasture"`
	UserId          string    `json:"-" db:"user_id"`
}

type LactationGroup struct {
	EntryDate        time.Time `json:"entryDate" db:"entry_date"`
	AnimalsNumber    int       `json:"animalsNumber" db:"animals_number"`
	TotalMilk        float64   `json:"totalMilk" db:"total_milk"`
	AverageMilk      float64   `json:"averageMilk" db:"avg_milk"`
	NumberDifference int       `json:"numberDifference" db:"number_difference"`
	AverageRate      float64   `json:"averageRate" db:"avg_rate"`
	TotalRate        float64   `json:"totalRate" db:"total_rate"`
}

type LactationGroupSave struct {
	OldEntry  time.Time `json:"oldEntry" db:"old_entry"`
	EntryDate time.Time `json:"entryDate" db:"entry_date"`
	UserId    string    `json:"-" db:"user_id"`
}

type LactationGroupFilter struct {
	MinEntryDate *time.Time `json:"minEntryDate" db:"entry_date"`
	MaxEntryDate *time.Time `json:"maxEntryDate" db:"entry_date"`
}
