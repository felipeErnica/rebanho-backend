package entity


type Page[E IDTO] struct {
    HasNextPage     bool    `json:"has_next_page"`
    NextCursor      string  `json:"next_cursor"`
    List            *[]E    `json:"list"`
}
