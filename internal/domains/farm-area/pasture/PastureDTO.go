package pasture

type PastureDTO struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Bull          *Bull  `json:"bull"`
	Farm          Farm   `json:"farm"`
	PastureSize   int    `json:"pastureSize"`
	AnimalsNumber int    `json:"animalsNumber"`
}

type Farm struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Bull struct {
	Id   string  `json:"id"`
	Tag  *string `json:"tag"`
	Name string  `json:"name"`
}
