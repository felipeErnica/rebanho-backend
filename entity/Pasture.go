package entity

import "time"

type Pasture struct {
    Id        string     `json:"id"`
    Bull      Bull       `json:"bull"`
    Name      string     `json:"name"`
    CreatedAt time.Time  `json:"created_at"`
}

type Bull struct {
    Id   *string  `json:"id"`
    Name *string  `json:"name"`
}

type PastureShort struct {
    Id   *string  `json:"id"`
    Name *string  `json:"name"`
}
