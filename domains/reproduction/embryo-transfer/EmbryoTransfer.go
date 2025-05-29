package embryoTransfer

import "time"

type EmbryoTransfer struct {
	Id             string     `json:"id" db:"id"`
	ReceiverId     string     `json:"receiverId" db:"receiver_id"`
	ReceiverNumber string     `json:"receiverNumber" db:"receiver_number"`
	ReceiverName   string     `json:"receiverName" db:"receiver_name"`
	DonorId        string     `json:"donorId" db:"donor_id"`
	DonorName      string     `json:"donorName" db:"donor_name"`
	BullId         string     `json:"bullId" db:"bull_id"`
	BullName       string     `json:"bullName" db:"bull_name"`
	TransferDate   time.Time  `json:"transferDate" db:"transfer_date"`
	Status         string     `json:"status" db:"status"`
	Observation    *string    `json:"observation" db:"observation"`
	CalfId         string     `json:"calfId" db:"calf_id"`
	LossId         string     `json:"lossId" db:"loss_id"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt      *time.Time `json:"DeletedAt" db:"deleted_at"`
	UserId         string     `json:"userId" db:"user_id"`
}

type EmbryoTransferSave struct {
	Id             string     `json:"id" db:"id"`
	ReceiverId     string     `json:"receiverId" db:"receiver_id"`
	DonorId        string     `json:"donorId" db:"donor_id"`
	BullId         string     `json:"bullId" db:"bull_id"`
	TransferDate   time.Time  `json:"transferDate" db:"transfer_date"`
	Status         string     `json:"status" db:"status"`
	Observation    *string    `json:"observation" db:"observation"`
	CalfId         string     `json:"calfId" db:"calf_id"`
	LossId         string     `json:"lossId" db:"loss_id"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
	DeletedAt      *time.Time `json:"DeletedAt" db:"deleted_at"`
	UserId         string     `json:"userId" db:"user_id"`
}
