package entity

type Page[E any] struct {
    HasNextPage     bool    `json:"hasNextPage"`
    NextCursor      string  `json:"nextCursor"`
    List            *[]E    `json:"list"`
}
