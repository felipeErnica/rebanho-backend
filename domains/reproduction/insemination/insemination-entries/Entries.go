package inseminationEntries

import "time"

type InseminationEntry struct {
    Id           string     `json:"id" db:"id"`
    AnimalId     string     `json:"animalId" db:"animal_id"`
    AnimalNumber string     `json:"animalNumber" db:"animal_number"`
    AnimalOrder  string     `json:"animalOrder" db:"animal_order"`
    AnimalName   string     `json:"animalName" db:"animal_name"`
    GroupId      string     `json:"groupId" db:"group_id"`
    GroupDate    time.Time  `json:"groupDate" db:"group_date"`
    BullName     string     `json:"bullName" db:"bull_name"`
    Observation  *string     `json:"observation" db:"observation"`
    Status       string     `json:"status" db:"status"`
    LossId       *string     `json:"lossId" db:"loss_id"`
    CalfId       *string     `json:"calfId" db:"calf_id"`
    CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
    DeletedAt    *time.Time `json:"deletedAt" db:"deleted_at"`
    UserId       string     `json:"userId" db:"user_id"`
}

type InseminationEntrySave struct {
    Id          string     `json:"id" db:"id"`
    AnimalId    string     `json:"animal" db:"animal_id"`
    GroupId     string     `json:"groupId" db:"group_id"`
    Observation *string     `json:"observation" db:"observation"`
    Status      string     `json:"status" db:"status"`
    LossId      *string     `json:"loss" db:"loss_id"`
    CalfId      *string     `json:"calf" db:"calf_id"`
    CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
    DeletedAt   *time.Time `json:"deletedAt" db:"deleted_at"`
    UserId      string     `json:"userId" db:"user_id"`
}
