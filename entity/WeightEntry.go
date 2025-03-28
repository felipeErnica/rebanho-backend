package entity

import "time"

type WeightEntry struct {
    Id        string       `json:"id"`
    Animal    AnimalShort  `json:"animal"`
    GroupId   string       `json:"group_id"`
    Weight    float32      `json:"weight"`
    CreatedAt time.Time    `json:"created_at"`
    DeletedAt time.Time    `json:"deleted_at"`
    UserId    string       `json:"user_id"`
}
