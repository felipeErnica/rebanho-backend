package entity

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/enums"
	"github.com/google/uuid"
)

type PregnancyLoss struct {
	Id       string
	AnimalId string
	LossType enums.LossType
	LossDate time.Time
}

func (p *PregnancyLoss) New(create *CreatePregnancyLoss) *PregnancyLoss {
    p = &PregnancyLoss{
        Id: uuid.NewString(),
        AnimalId: create.AnimalId,
        LossType: create.LossType,
        LossDate: create.LossDate,
    }
    return p
}

type CreatePregnancyLoss struct {
	AnimalId string
	LossType enums.LossType
	LossDate time.Time
}
