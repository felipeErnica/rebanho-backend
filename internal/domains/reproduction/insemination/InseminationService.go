package insemination

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type InseminationService struct {
	Repo *InseminationRepository
}

func NewService(repo *InseminationRepository) *InseminationService {
	return &InseminationService{repo}
}

func (s *InseminationService) GetBirthRateStats(userId string) (*CardStats, *log.APIError) {
	result, err := s.Repo.GetBirthRateStats(userId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	birthRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(birthRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = birthRates[lenght-1].BirthRate
		previousRate = 0
		trend = 0
	default:
		currentRate = birthRates[lenght-1].BirthRate
		previousRate = birthRates[lenght-2].BirthRate
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardStats{
		Hist:    birthRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (s *InseminationService) GetPregnancyRateStats(userId string) (*CardStats, *log.APIError) {
	result, err := s.Repo.GetPregnancyRateStats(userId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	pregnancyRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(pregnancyRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = pregnancyRates[lenght-1].PregnancyRate
		previousRate = 0
		trend = 0
	default:
		currentRate = pregnancyRates[lenght-1].PregnancyRate
		previousRate = pregnancyRates[lenght-2].PregnancyRate
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardStats{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (s *InseminationService) GetAnimalsNumber(userId string) (*CardStats, error) {
	result, err := s.Repo.GetAnimalsNumber(userId)
	if err != nil {
		return nil, err
	}

	pregnancyRates := *result
	var currentRate, previousRate, trend float64

	switch lenght := len(pregnancyRates); lenght {
	case 0:
		currentRate = 0
		previousRate = 0
		trend = 0
	case 1:
		currentRate = pregnancyRates[lenght-1].AnimalsNumber
		previousRate = 0
		trend = 0
	default:
		currentRate = pregnancyRates[lenght-1].AnimalsNumber
		previousRate = pregnancyRates[lenght-2].AnimalsNumber
		trend = util.CalculatePercentageTrend(currentRate, previousRate)
	}

	stats := &CardStats{
		Hist:    pregnancyRates,
		Current: currentRate,
		Trend:   trend,
	}

	return stats, nil
}

func (s *InseminationService) Add(entry *InseminationEntrySave) *log.APIError {
	exists, err := s.Repo.CheckAddConflicts(*entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists && !entry.Overwrite {
		return log.ConflictAPIWarning("Já existe uma entrada desta vaca na mesma data! Deseja substituí-la por esta?")
	}

	err = s.Repo.Add(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *InseminationService) Update(entry *InseminationEntrySave) (*InseminationEntry, *log.APIError) {
	validation, err := s.Repo.CheckUpdateConflicts(*entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if validation.Repeated {
		return nil, log.ConflictAPIError("Já existe uma inseminação desta vaca na mesma data!")
	}

	if !entry.IgnoreWarnings {
		warnMsg := make([]string, 0)
		if validation.HasChildren {
			warnMsg = append(warnMsg, "A vaca possui uma cria, cujo o pai é o touro desta inseminação.")
		}

		if validation.IsPregnant {
			warnMsg = append(warnMsg, "A vaca possui uma prenhez ligada a esta inseminação.")
		}

		if len(warnMsg) != 0 {
			msgBody := util.FormatWarningMessage(warnMsg...)
			return nil, log.ConflictAPIWarning(msgBody + "\nDeseja alterar mesmo assim?")
		}
	}

	res, err := s.Repo.Update(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return res, nil
}

func (s *InseminationService) Delete(params *InseminationEntryDelete) *log.APIError {
	validation, err := s.Repo.CheckDeleteConflicts(params)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if validation.IsPregnant {
		return log.ConflictAPIWarning("A vaca possui uma prenhez ligada a esta inseminação. Deseja excluir mesmo assim?")
	}

	if validation.HasChildren {
		return log.NewAPIWarning(
			"Inseminação gerou crias!",
			"A vaca possui uma cria, cujo o pai é o touro desta inseminação. Deseja substituí-lo?",
			"ChildrenWarning",
		)
	}

	err = s.Repo.Delete(params)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *InseminationService) UpdateGroup(group *InseminationGroupSave) (*InseminationGroup, *log.APIError) {
	exists, err := s.Repo.CheckGroupConflicts(*group)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if exists {
		return nil, log.ConflictAPIError("Já existem registros de inseminação nesta mesma data! " +
			"Para evitar conflitos, altere os registros existentes ou escolha outra data.",
		)
	}

	res, err := s.Repo.UpdateGroup(group)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return res, nil
}

func (s *InseminationService) DeleteGroup(params *InseminationGroupDelete) *log.APIError {
	err := s.Repo.DeleteGroup(params)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
