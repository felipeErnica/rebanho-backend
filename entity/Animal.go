package entity

import "time"

//type Animal struct {
    //Id                   string          `json:"id"`
    //Name                 *string         `json:"name"`
    //IdentificationNumber *string         `json:"identification_number"`
    //RingOrder            int             `json:"ring_order"`  
    //Sex                  string          `json:"sex"`
    //Father               AnimalShort     `json:"father"`
    //Mother               AnimalShort     `json:"mother"`
    //BirthDate            *time.Time      `json:"birth_date"`
    //DeathDate            *time.Time      `json:"death_date"`
    //Pasture              PastureShort    `json:"pasture"`
    //Status               *string         `json:"status"`
    //Isr                  float32         `json:"isr"`
    //AvarageProd          float32         `json:"avarage_prod"`
    //AvarageBirthInterval float32         `json:"avarage_birth_interval"`
    //MaxPeak              float32         `json:"max_peak"`
    //ChildrenQuantity     int             `json:"children_quantity"`
    //CreatedAt            time.Time       `json:"created_at"`
    //DeletedAt            *time.Time      `json:"deleted_at"`
//}

type Animal struct {
    Id                   string          `json:"id"`
    Name                 *string         `json:"name"`
    IdentificationNumber *string         `json:"identification_number"`
    RingOrder            int             `json:"ring_order"`  
    Sex                  string          `json:"sex"`
    FatherId             *string         `json:"father_id"`
    MotherId             *string         `json:"mother_id"`
    BirthDate            *time.Time      `json:"birth_date"`
    DeathDate            *time.Time      `json:"death_date"`
    PastureId            *string         `json:"pasture_id"`
    Status               *string         `json:"status"`
    Isr                  float32         `json:"isr"`
    AvarageProd          float32         `json:"avarage_prod"`
    AvarageBirthInterval float32         `json:"avarage_birth_interval"`
    MaxPeak              float32         `json:"max_peak"`
    ChildrenQuantity     int             `json:"children_quantity"`
    CreatedAt            time.Time       `json:"created_at"`
    DeletedAt            *time.Time      `json:"deleted_at"`
}

type AnimalShort struct {
    Id                   *string
    IdentificationNumber *string
    Name                 *string
}
