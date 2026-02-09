package milk

import "time"

type MilkDTO struct {
	Id        string    `json:"id"`
	Cow       Cow       `json:"cow"`
	Pasture   *Pasture  `json:"pasture"`
	EntryDate time.Time `json:"entryDate"`
	Quantity  float64   `json:"quantity"`
}

type Cow struct {
	Id   string  `json:"id"`
	Tag  *string `json:"tag"`
	Name *string `json:"name"`
}

type Pasture struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Farm Farm   `json:"farm"`
}

type Farm struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
