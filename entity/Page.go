package entity


type Page[E IDTO] struct {
    HasNextPage     bool    `json:"hasNextPage"`
    NextCursor      string  `json:"nextCursor"`
    List            *[]E    `json:"list"`
}
