package entity

import "time"

type EmbryoTransfer struct {
	Id             string     `json:"id"`
	ReceiverId     string     `json:"receiverId"`
	ReceiverNumber string     `json:"receiverNumber"`
	ReceiverName   string     `json:"receiverName"`
	DonorId        string     `json:"donorId"`
	DonorName      string     `json:"donorName"`
	BullId         string     `json:"bullId"`
	BullName       string     `json:"bullName"`
	TransferDate   time.Time  `json:"transferDate"`
	Status         string     `json:"status"`
	Observation    *string    `json:"observation"`
	CalfId         string     `json:"calfId"`
	LossId         string     `json:"lossId"`
	CreatedAt      time.Time  `json:"createdAt"`
	DeletedAt      *time.Time `json:"DeletedAt"`
	UserId         string     `json:"userId"`
}
