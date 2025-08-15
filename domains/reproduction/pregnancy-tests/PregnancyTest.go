package pregnancyTests

import "time"

type TestEntry struct {
	Id            string     `json:"id" db:"id"`
	GroupId       string     `json:"groupId" db:"group_id"`
	GroupDate     string     `json:"groupDate" db:"group_date"`
	AnimalId      string     `json:"animalId" db:"animal_id"`
	AnimalName    string     `json:"animalName" db:"animal_name"`
	AnimalOrder   int        `json:"animalOrder" db:"animal_order"`
	AnimalNumber  string     `json:"animalNumber" db:"animal_number"`
	IsPregnant    bool       `json:"isPregnant" db:"is_pregnant"`
	BirthForecast *time.Time `json:"birthForecast" db:"birth_forecast"`
	Observation   *string    `json:"observation" db:"observation"`
	Status        string     `json:"status" db:"status"`
	LossId        *string    `json:"lossId" db:"loss_id"`
	CalfId        *string    `json:"calfId" db:"calf_id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        string     `json:"userId" db:"user_id"`
}

type TestEntrySave struct {
	Id            string     `json:"id" db:"id"`
	GroupId       string     `json:"groupId" db:"group_id"`
	AnimalId      string     `json:"animalId" db:"animal_id"`
	IsPregnant    bool       `json:"isPregnant" db:"is_pregnant"`
	BirthForecast *time.Time `json:"birthForecast" db:"birth_forecast"`
	Observation   *string    `json:"observation" db:"observation"`
	Status        string     `json:"status" db:"status"`
	LossId        *string    `json:"lossId" db:"loss_id"`
	CalfId        *string    `json:"calfId" db:"calf_id"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        string     `json:"userId" db:"user_id"`
}

type PregnancyHist struct {
	TestDate      time.Time `json:"testDate" db:"test_date"`
	PregnancyRate float64   `json:"pregnancyRate" db:"pregnancy_rate"`
}

type PregnancyStats struct {
	Current float64         `json:"current"`
	Trend   float64         `json:"trend"`
	Hist    []PregnancyHist `json:"hist"`
}

type BirthHist struct {
	TestDate  time.Time `json:"testDate" db:"test_date"`
	BirthRate float64   `json:"birthRate" db:"birth_rate"`
}

type BirthStats struct {
	Current float64     `json:"current"`
	Trend   float64     `json:"trend"`
	Hist    []BirthHist `json:"hist"`
}

type PregnancyTestHist struct {
	TestDate      time.Time `json:"testDate" db:"test_date"`
	Totals        int       `json:"totals" db:"totals"`
	PregnancyRate float64   `json:"pregnancyRate" db:"pregnancy_rate"`
	BirthRate     float64   `json:"birthRate" db:"birth_rate"`
}
