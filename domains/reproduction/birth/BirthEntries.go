package birth

import "time"

type BirthEntry struct {
	Id            string     `json:"id" db:"id"`
	MotherId      string     `json:"motherId" db:"mother_id"`
	MotherName    string     `json:"-" db:"mother_name"`
	MotherInfo    string     `json:"motherInfo" db:"mother_info"`
	MotherOrder   int        `json:"motherOrder" db:"mother_order"`
	CalfId        string     `json:"calfId" db:"calf_id"`
	CalfName      string     `json:"calfName" db:"calf_name"`
	CalfBirthDate time.Time  `json:"calfBirthDate" db:"calf_birth_date"`
	CalfSex       string     `json:"calfSex" db:"calf_sex"`
	CalfFatherId  *string    `json:"calfFatherId" db:"calf_father_id"`
	CalfFather    *string    `json:"calfFather" db:"calf_father"`
	BirthInterval *int       `json:"birthInterval" db:"birth_interval"`
	Observation   *string    `json:"observation" db:"observation"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt     *time.Time `json:"deletedAt" db:"deleted_at"`
	UserId        string     `json:"userId" db:"user_id"`
}

type BirthEntryFilter struct {
	IsFiltered       bool       `json:"isFiltered"`
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
	MotherId    *string   `json:"motherId" db:"mother_id"`
	FatherId    *string   `json:"fatherId" db:"father_id"`
	BirthDate   time.Time `json:"birthDate" db:"birth_date"`
	Sex         string    `json:"sex" db:"sex"`
	Observation *string   `json:"observation" db:"observation"`
	UserId      string    `json:"-" db:"user_id"`
}

type BirthFooter struct {
	Total           int     `json:"total"`
	IntervalAverage float64 `json:"intervalAverage" db:"interval_average"`
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

type BirthsNumberEntry struct {
	EntryDate  time.Time `json:"entryDate" db:"entry_date"`
	BirthTotal float64   `json:"birthTotal" db:"birth_total"`
}

type DeathsNumberEntry struct {
	EntryDate   time.Time `json:"entryDate" db:"entry_date"`
	DeathsTotal float64   `json:"deathsTotal" db:"deaths_total"`
}

type BirthIntervalHist struct {
	BirthDate       time.Time `json:"birthDate" db:"birth_date"`
	IntervalAverage float64   `json:"intervalAverage" db:"interval_average"`
}

type DeathIndexHist struct {
	EntryDate  time.Time `json:"entryDate" db:"date"`
	DeathIndex float64   `json:"deathIndex" db:"death_index"`
}

type LossHist struct {
	Month  time.Time `json:"month" db:"month"`
	Losses int       `json:"losses" db:"losses"`
}

type IntervalAnimal struct {
	AnimalName      string  `json:"animalName" db:"animal_name"`
	BirthNumbers    int     `json:"birthNumbers" db:"birth_numbers"`
	IntervalAverage float64 `json:"intervalAverage" db:"interval_average"`
	AverageRate     float64 `json:"averageRate" db:"average_rate"`
	Score           float64 `json:"-" db:"reproductive_score"`
}

type BirthStats struct {
	CurrentInterval     float64             `json:"currentInterval"`
	IntervalTrend       float64             `json:"intervalTrend"`
	IntervalHist        []BirthIntervalHist `json:"intervalHist"`
	DeathIndex          float64             `json:"deathIndex"`
	DeathTrend          float64             `json:"deathTrend"`
	DeathIndexHist      []DeathIndexHist    `json:"deathIndexHist"`
	CurrentBirthNumbers int                 `json:"currentBirthNumbers"`
	BirthNumbersTrend   int                 `json:"birthNumbersTrend"`
	PregnantsNumber     int                 `json:"pregnantsNumber"`
	Losses              int                 `json:"losses"`
	LossTrend           int                 `json:"lossTrend"`
	LossHist            []LossHist          `json:"lossHist"`
	BirthHistory        []BirthsByDate      `json:"birthHistory" db:"birth_history"`
}

type IntervalStats struct {
	CurrentInterval   float64             `json:"currentInterval"`
	IntervalTrend     float64             `json:"intervalTrend"`
	BirthIntervalHist []BirthIntervalHist `json:"intervalHist"`
}

type DeathStats struct {
	DeathIndex      float64          `json:"currentDeathIndex"`
	DeathIndexTrend float64          `json:"deathTrend"`
	DeathIndexHist  []DeathIndexHist `json:"deathIndexHist"`
}

type LossStats struct {
	LossNumbers int        `json:"lossNumbers"`
	LossTrend   int        `json:"lossTrend"`
	LossHist    []LossHist `json:"lossHist"`
}

type CurrentStats struct {
	Current float64 `json:"current"`
	Trend   float64 `json:"trend"`
	Hist    any     `json:"hist"`
}
