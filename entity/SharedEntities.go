package entity

type Page[E any] struct {
	HasNextPage bool   `json:"hasNextPage"`
	NextCursor  string `json:"nextCursor"`
	List        *[]E   `json:"list"`
}

type SearchEntity struct {
    Id    string `json:"id" db:"id"`
	Label string `json:"label" db:"label"`
}
