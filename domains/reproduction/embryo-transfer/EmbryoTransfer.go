package embryoTransfer

import "time"

type EmbryoTransfer struct {
	Id               string     `json:"id" db:"id"`
	ReceiverId       string     `json:"receiverId" db:"receiver_id"`
	ReceiverOrder    int        `json:"receiverOrder" db:"receiver_order"`
	ReceiverName     string     `json:"receiverName" db:"receiver_name"`
	DonorId          string     `json:"donorId" db:"donor_id"`
	DonorOrder       int        `json:"donorOrder" db:"donor_order"`
	DonorName        string     `json:"donorName" db:"donor_name"`
	BullId           string     `json:"bullId" db:"bull_id"`
	BullName         string     `json:"bullName" db:"bull_name"`
	TransferDate     time.Time  `json:"transferDate" db:"transfer_date"`
	BirthStatus      string     `json:"birthStatus" db:"birth_status"`
	PregnancyStatus  string     `json:"pregnancyStatus" db:"pregnancy_status"`
	Observation      *string    `json:"observation" db:"observation"`
	ChildInformation *string    `json:"childInformation" db:"child_information"`
	CreatedAt        time.Time  `json:"-" db:"created_at"`
	DeletedAt        *time.Time `json:"-" db:"deleted_at"`
	UserId           string     `json:"-" db:"user_id"`
}

type TransferEntryFilter struct {
	IsFiltered          bool       `json:"isFiltered"`
	Animals             *[]string  `json:"animals" db:"animal_id"`
	Bulls               *[]string  `json:"bulls" db:"bull_id"`
	MinInseminationDate *time.Time `json:"minMatingDate" db:"transfer_date"`
	MaxInseminationDate *time.Time `json:"maxMatingDate" db:"transfer_date"`
	BirthStatus         *string    `json:"birthStatus" db:"birth_status"`
	PregnancyStatus     *string    `json:"pregnancyStatus" db:"pregnancy_status"`
}

type TransferFoot struct {
	Totals               int     `json:"totals" db:"totals"`
	AverageBirthRate     float64 `json:"averageBirthRate" db:"average_birth_rate"`
	AveragePregnancyRate float64 `json:"averagePregnancyRate" db:"average_pregnancy_rate"`
}

type CardEntry struct {
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
	Hist    any     `json:"hist"`
}

type BirthRateEntry struct {
	TransferDate time.Time `json:"transferDate" db:"transfer_date"`
	BirthRate    float64   `json:"birthRate" db:"birth_rate"`
}

type PregnancyRateEntry struct {
	TransferDate  time.Time `json:"transferDate" db:"transfer_date"`
	PregnancyRate float64   `json:"pregnancyRate" db:"pregnancy_rate"`
}

type AnimalsNumberEntry struct {
	TransferDate  time.Time `json:"transferDate" db:"transfer_date"`
	AnimalsNumber float64   `json:"animalsNumber" db:"animals_number"`
}

type TransferHist struct {
	TransferDate      time.Time `json:"transferDate" db:"transfer_date"`
	AnimalsNumber     int       `json:"animalsNumber" db:"animals_number"`
	PregnanciesNumber int       `json:"pregnanciesNumber" db:"pregnancies_number"`
	BirthsNumber      int       `json:"birthsNumber" db:"births_number"`
}

type FutureBirths struct {
	BirthForecast time.Time `json:"birthForecast" db:"birth_forecast"`
	BirthsNumber  int       `json:"birthsNumber" db:"births_number"`
}

type BestAnimals struct {
	AnimalName                string  `json:"animalName" db:"animal_name"`
	Total                   int     `json:"total" db:"total"`
	BirthRate               float64 `json:"birthRate" db:"birth_rate"`
	PregnancyRate           float64 `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthComparisonRate     float64 `json:"birthComparisonRate" db:"birth_comparison_rate"`
	PregnancyComparisonRate float64 `json:"pregnancyComparisonRate" db:"pregnancy_comparison_rate"`
}

type TransferGroup struct {
	TransferDate            time.Time `json:"transferDate" db:"transfer_date"`
	CowNumber               float64   `json:"cowNumber" db:"cow_number"`
	BirthRate               float64   `json:"birthRate" db:"birth_rate"`
	PregnancyRate           float64   `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthComparisonRate     float64   `json:"birthComparisonRate" db:"birth_comparison_rate"`
	PregnancyComparisonRate float64   `json:"pregnancyComparisonRate" db:"pregnancy_comparison_rate"`
}

type LastEntry struct {
	TransferDate time.Time        `json:"transferDate"`
	Entries      []EmbryoTransfer `json:"entries"`
}
