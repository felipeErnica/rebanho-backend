package entity

import  "time"

type SlaughterGroup struct {
    Id               string               `json:"id"`
    Slaughterhouse   SlaughterhouseShort  `json:"slaughterhouse"`
    WeightDecrease   float32              `json:"weight_decrease"`
    SlaughterDate    time.Time            `json:"slaughter_date"`
    CreatedAt        time.Time            `json:"created_at"`
    DeletedAt        *time.Time           `json:"deleted_at"`
    UserId           string               `json:"user_id"`
}

type SlaughterGroupShort struct {
    Id               string               `json:"id"`
    WeightDecrease   float32              `json:"weight_decrease"`
    SlaughterDate    time.Time            `json:"slaughter_date"`
}
