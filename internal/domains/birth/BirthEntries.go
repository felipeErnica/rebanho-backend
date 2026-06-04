package birth

import "time"

type BirthDB struct {
	CalfId          string     `db:"calf_id"`
	CalfName        *string    `db:"calf_name"`
	CalfTag         *string    `db:"calf_tag"`
	CalfBirthDate   time.Time  `db:"calf_birth_date"`
	CalfSex         string     `db:"calf_sex"`
	CalfObservation *string    `db:"calf_observation"`
	MotherId        string     `db:"mother_id"`
	MotherName      *string    `db:"mother_name"`
	MotherTag       *string    `db:"mother_tag"`
	MotherOrder     int        `db:"mother_order"`
	FatherId        *string    `db:"father_id"`
	FatherName      *string    `db:"father_name"`
	FatherTag       *string    `db:"father_tag"`
	BirthInterval   *int       `db:"birth_interval"`
	CreatedAt       time.Time  `db:"created_at"`
	DeletedAt       *time.Time `db:"deleted_at"`
	UserId          string     `db:"user_id"`
}

type BirthEntryFilter struct {
	Mothers          *[]string  `json:"mothers" db:"mother_id"`
	MinBirthDate     *time.Time `json:"minBirthDate" db:"calf_birth_date"`
	MaxBirthDate     *time.Time `json:"maxBirthDate" db:"calf_birth_date"`
	Sex              *string    `json:"sex" db:"calf_sex"`
	Fathers          *[]string  `json:"fathers" db:"calf_father_id"`
	MinBirthInterval *int       `json:"minBirthInterval" db:"birth_interval"`
	MaxBirthInterval *int       `json:"maxBirthInterval" db:"birth_interval"`
}

type BirthEntrySave struct {
	Id          string    `json:"id" db:"id"`
	Tag         *string   `json:"tag" db:"tag"`
	MotherId    *string   `json:"motherId" db:"mother_id"`
	FatherId    *string   `json:"fatherId" db:"father_id"`
	PastureId   *string   `json:"pastureId" db:"pasture_id"`
	BirthDate   time.Time `json:"birthDate" db:"birth_date"`
	Sex         string    `json:"sex" db:"sex"`
	Observation *string   `json:"observation" db:"observation"`
	Overwrite   bool      `json:"overwrite"`
	IgnoreTag   bool      `json:"ignoreTag"`
	UserId      string    `json:"-" db:"user_id"`
}

type BirthFooter struct {
	Total           int      `json:"total"`
	IntervalAverage *float64 `json:"intervalAverage" db:"interval_average"`
}

type TotalBirthsBySex struct {
	BirthMonth time.Time `json:"birthMonth" db:"birth_month"`
	Males      int       `json:"males" db:"males"`
	Females    int       `json:"females" db:"females"`
}

type BirthsBySex struct {
	Males   int `json:"males" db:"males"`
	Females int `json:"females" db:"females"`
}

type BirthsByDate struct {
	Date       time.Time `json:"date" db:"date"`
	BirthTotal int       `json:"birthTotal" db:"birth_total"`
	DeathTotal int       `json:"deathTotal" db:"death_total"`
}

type IntervalAnimal struct {
	AnimalName      string  `json:"animalName" db:"animal_name"`
	BirthNumbers    int     `json:"birthNumbers" db:"birth_numbers"`
	IntervalAverage float64 `json:"intervalAverage" db:"interval_average"`
	AverageRate     float64 `json:"averageRate" db:"average_rate"`
	Score           float64 `json:"-" db:"reproductive_score"`
}
