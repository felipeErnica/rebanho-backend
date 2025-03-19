package entity

type SlaughterEntry struct {
    Id         string  `json:"id"`
    GroupId    string  `json:"group_id"`
    Weight     string  `json:"weight"`
    DeadWeight string  `json:"dead_weight"`
}
