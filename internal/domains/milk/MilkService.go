package milk

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type MilkService struct {
	Repo *MilkRepository
}

func NewService(repo *MilkRepository) *MilkService {
	return &MilkService{repo}
}

func (s *MilkService) UpdateGroup(groupEntry *LactationGroupSave) (*LactationGroup, *log.APIError) {

	exists, err := s.Repo.CheckGroupUpdateConflicts(*groupEntry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if exists {
		return nil, log.ConflictAPIError("Já exitem marcações nesta data!")
	}

	res, err := s.Repo.UpdateGroup(groupEntry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return res, nil
}

func (s *MilkService) Add(entry *MilkEntrySave) *log.APIError {

	validate, err := s.Repo.CheckMilkEntryConflicts(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if validate.HasLac {
		return log.IncorrectEntityAPIError("Não há nenhuma lactação correspondente a esta marcação!")
	}

	if validate.EntryExist && entry.Id != nil {
		return log.ConflictAPIError("Já exite uma marcação com esta data!")
	}

	if validate.EntryExist && !entry.Overwrite {
		return log.ConflictAPIWarning(`Esta marcação já existe! Deseja substituí-la por esta?`)
	}

	if entry.PastureId != nil && validate.IsDifferentPasture && !entry.TransferPasture {
		return log.NewAPIWarning(
			"Pasto diferente!",
			"A vaca não está no pasto informado! Deseja transferi-la?",
			"PastureWarning",
		)
	}

	err = s.Repo.Add(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *MilkService) Update(entry *MilkEntrySave) (*MilkEntry, *log.APIError) {

	validate, err := s.Repo.CheckMilkEntryConflicts(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if validate.HasLac {
		return nil, log.IncorrectEntityAPIError("Não há nenhuma lactação correspondente a esta marcação!")
	}

	if validate.EntryExist && entry.Id != nil {
		return nil, log.ConflictAPIError("Já exite uma marcação com esta data!")
	}

	if validate.EntryExist && !entry.Overwrite {
		return nil, log.ConflictAPIWarning(`Esta marcação já existe! Deseja substituí-la por esta?`)
	}

	if entry.PastureId != nil && validate.IsDifferentPasture && !entry.TransferPasture {
		return nil, log.NewAPIWarning(
			"Pasto diferente!",
			"A vaca não está no pasto informado! Deseja transferi-la?",
			"PastureWarning",
		)
	}

	res, err := s.Repo.Update(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return res, nil
}

func (s *MilkService) GetLastMilk(userId string) (*CardContainer, error) {
	averageHist, err := s.Repo.GetLastMilkEntries(userId)
	if err != nil {
		return nil, err
	}

	var current, previous, trend float64

	switch lenght := len(*averageHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = (*averageHist)[0].TotalMilk
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].TotalMilk
		previous = (*averageHist)[lenght-2].TotalMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *MilkService) GetYearMilk(userId string) (*CardContainer, error) {
	averageHist, err := s.Repo.GetYearMilkEntries(userId)
	if err != nil {
		return nil, err
	}

	var current, previous, trend float64

	switch lenght := len(*averageHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = (*averageHist)[0].TotalMilk
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].TotalMilk
		previous = (*averageHist)[lenght-2].TotalMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *MilkService) GetLastAverageMilk(userId string) (*CardContainer, error) {
	averageHist, err := s.Repo.GetLastAverageMilkEntries(userId)
	if err != nil {
		return nil, err
	}

	var current, previous, trend float64

	switch lenght := len(*averageHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = (*averageHist)[0].AverageMilk
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].AverageMilk
		previous = (*averageHist)[lenght-2].AverageMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *MilkService) GetYearAverageMilk(userId string) (*CardContainer, error) {
	averageHist, err := s.Repo.GetYearAverageMilkEntries(userId)
	if err != nil {
		return nil, err
	}

	var current, previous, trend float64

	switch lenght := len(*averageHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = (*averageHist)[0].AverageMilk
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].AverageMilk
		previous = (*averageHist)[lenght-2].AverageMilk
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}
