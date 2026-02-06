package lactation

import (
	"bytes"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type LactationService struct {
	Repo *LactationRepository
}

func NewLactationService(repo *LactationRepository) *LactationService {
	return &LactationService{repo}
}

func (s *LactationService) AddLactation(lac *LactationHistSave) *log.APIError {

	validate, err := s.Repo.CheckLactationConflicts(*lac)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	validateErr := s.saveValidation(validate, lac)
	if validateErr != nil {
		return validateErr
	}

	err = s.Repo.AddLactation(lac)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *LactationService) UpdateLactation(lac *LactationHistSave) (*LactationHist, *log.APIError) {

	validate, err := s.Repo.CheckLactationConflicts(*lac)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	validateErr := s.saveValidation(validate, lac)
	if validateErr != nil {
		return nil, validateErr
	}

	updatedLac, err := s.Repo.UpdateLactation(lac)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return updatedLac, nil
}

func (s *LactationService) saveValidation(validate *SaveValidation, lac *LactationHistSave) *log.APIError {

	apiErrors := make([]string, 0)
	if lac.EndDate != nil && lac.StartDate.After(*lac.EndDate) {
		apiErrors = append(apiErrors, "A data de início não pode ser maior que a data de encerramento!")
	}

	if validate.InvalidStart {
		apiErrors = append(apiErrors,
			"A data de início informada está em conflito com a data de encerramento "+
				"da lactação anterior. A data de início informada é menor que a data de encerramento anterior!",
		)
	}

	if validate.InvalidNew {
		apiErrors = append(apiErrors, "Não é possível adicionar uma nova lactação enquanto a antiga não for encerrada!")
	}

	if validate.InvalidEmptyEnd {
		apiErrors = append(apiErrors, "Não é possível adicionar uma lactação em aberto (sem encerramento), pois já existe uma lactação posterior!")
	}

	if validate.InvalidEnd {
		apiErrors = append(apiErrors,
			"A data de encerramento informada está em conflito com a data de início "+
				"de uma lactação posterior. A data de encerramento informada é maior que a data de início da lactação posterior!",
		)
	}

	if validate.InvalidCalf {
		apiErrors = append(apiErrors, "O bezerro selecionado está vinculado a outra lactação!")
	}

	if len(apiErrors) != 0 {
		var errBuff bytes.Buffer
		for i, msg := range apiErrors {
			errPoint := fmt.Sprintf("\n%d - %s", i+1, msg)
			errBuff.WriteString(errPoint)
		}
		errMsg := fmt.Sprintf("Os seguintes erros foram encontrados: %s", errBuff.String())
		return log.IncorrectEntityAPIError(errMsg)
	}

	if lac.PastureId != nil && validate.DifferentPasture && !lac.TransferPasture {
		return log.NewAPIWarning(
			"Pasto diferente!",
			"A vaca não está no pasto selecionado! Deseja transferí-la?",
			"PastureWarning",
		)
	}

	if validate.LactationExists && lac.Id != nil {
		return log.ConflictAPIError("Já existe uma lactação desta vaca na mesma data!")
	}

	if validate.LactationExists && !lac.Overwrite {
		return log.ConflictAPIWarning("Esta lactação já existe! Deseja substituí-la por esta?")
	}

	return nil
}

func (s *LactationService) GetLastLactating(userId string) (*CardContainer, error) {
	averageHist, err := s.Repo.GetLastLactatingEntries(userId)
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
		current = (*averageHist)[0].AnimalsNumber
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].AnimalsNumber
		previous = (*averageHist)[lenght-2].AnimalsNumber
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *LactationService) GetLastDry(userId string) (*CardContainer, error) {
	averageHist, err := s.Repo.GetLastDryEntries(userId)
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
		current = (*averageHist)[0].AnimalsNumber
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].AnimalsNumber
		previous = (*averageHist)[lenght-2].AnimalsNumber
		trend = current - previous
	}

	return &CardContainer{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}
