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
func (s *LactationService) toDTO(entry LactationDB) LactationDTO {
	dto := LactationDTO{
		Id:                entry.Id,
		StartDate:         entry.StartDate,
		EndDate:           entry.EndDate,
		LacInterval:       entry.LacInterval,
		LacPeriod:         entry.LacPeriod,
		AverageProduction: entry.AverageProduction,
		TotalProduction:   entry.TotalProduction,
		Peak:              entry.Peak,
		Observation:       entry.Observation,
		Cow: Animal{
			Id:   entry.AnimalId,
			Tag:  entry.AnimalTag,
			Name: entry.AnimalName,
		},
	}

	if entry.CalfId != nil {
		dto.Calf = &Calf{
			Id:        *entry.CalfId,
			Tag:       entry.CalfTag,
			Name:      entry.CalfName,
			Sex:       *entry.CalfSex,
			BirthDate: entry.CalfBirthDate,
			DeathDate: entry.CalfDeathDate,
		}
	}

	return dto
}

func (s *LactationService) listToDTO(list *[]LactationDB) *[]LactationDTO {
	listDTO := make([]LactationDTO, 0)
	for _, entry := range *list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return &listDTO
}

func (s *LactationService) GetLastLactating(userId string) (*util.CardStats, error) {
	averageHist, err := s.Repo.GetLastLactatingEntries(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*averageHist), nil
}

func (s *LactationService) GetLastDry(userId string) (*util.CardStats, error) {
	averageHist, err := s.Repo.GetLastDryEntries(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*averageHist), nil
}

func (s *LactationService) GetDairyTypes(userId string) (*DairyTypes, error) {
	return s.Repo.GetDairyTypes(userId)
}

func (s *LactationService) GetBestAnimals(userId string) (*[]AnimalsRating, error) {
	return s.Repo.GetBestAnimals(userId)
}

func (s *LactationService) GetWorstAnimals(userId string) (*[]AnimalsRating, error) {
	return s.Repo.GetWorstAnimals(userId)
}

func (s *LactationService) GetBestMothers(userId string) (*[]ParentsRating, error) {
	return s.Repo.GetBestMothers(userId)
}

func (s *LactationService) GetWorstMothers(userId string) (*[]ParentsRating, error) {
	return s.Repo.GetWorstMothers(userId)
}

func (s *LactationService) GetBestFathers(userId string) (*[]ParentsRating, error) {
	return s.Repo.GetBestFathers(userId)
}

func (s *LactationService) GetWorstFathers(userId string) (*[]ParentsRating, error) {
	return s.Repo.GetWorstFathers(userId)
}

func (s *LactationService) FindLactationPage(
	filter *LactationHistFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*util.Page[LactationDTO], error) {
	sort = util.AddCommonFields(sort)
	list, err := s.Repo.FindLactationPage(filter, sort, order, cursor, limit, userId)
	if err != nil {
		return nil, err
	}

	listDto := s.listToDTO(list)
	page := util.NewPage(*listDto, cursor, limit)
	return page, nil
}

func (s *LactationService) GetLactationPageFoot(filter *LactationHistFilter, userId string) (*LactationHistFoot, error) {
	return s.Repo.GetLactationPageFoot(filter, userId)
}

func (s *LactationService) FindAnimalsPage(
	filter *AnimalFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*util.Page[AnimalDTO], error) {
	sort = util.AddCommonFields(sort)
	list, err := s.Repo.FindAnimalsPage(filter, sort, order, cursor, limit, userId)
	if err != nil {
		return nil, err
	}

	listDTO := make([]AnimalDTO, 0)
	for _, entry := range *list {
		dto := AnimalDTO{
			Id: entry.Id,
			Name: entry.Name,
			Tag: entry.Tag,
			BirthDate: entry.BirthDate,
		}

		if entry.LacId != nil {
			dto.Lactation = &LactationDTO{
				Id: *entry.LacId,
				StartDate: *entry.LacStart,
				EndDate: entry.LacEnd,
				LacInterval: entry.LacInterval,
				LacPeriod: entry.LacPeriod,
				AverageProduction: entry.LacAverage,
				TotalProduction: entry.LacTotal,
				Peak: entry.LacPeak,
				Observation: entry.LacObservation,
			}

			if entry.CalfId != nil {
				dto.Lactation.Calf = &Calf{
					Id: *entry.CalfId,
					Name: entry.CalfName,
					Tag: entry.CalfTag,
					Sex: *entry.CalfSex,
					BirthDate: entry.CalfBirthDate,
					DeathDate: entry.CalfDeathDate,
				}
			}
		}

		listDTO = append(listDTO, dto)
	}

	newCursor := util.CreateCursorKey(sort, *list)
	page := util.NewPage(listDTO, newCursor, limit)
	return page, nil
}

func (s *LactationService) GetAnimalsPageFoot(filter *AnimalFilter, userId string) (*LactationHistFoot, error) {
	return s.Repo.GetAnimalsPageFoot(filter, userId)
}

func (s *LactationService) FindById(id string, userId string) (*LactationDTO, error) {
	entry, err := s.Repo.FindById(id, userId)
	if err != nil {
		return nil, err
	}
	dto := s.toDTO(*entry)
	return &dto, nil
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

func (s *LactationService) UpdateLactation(lac *LactationHistSave) (*LactationDTO, *log.APIError) {

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

	dto := s.toDTO(*updatedLac)
	return &dto, nil
}

func (s *LactationService) DeleteLactation(id string, userId string) error {
	return s.Repo.DeleteLactation(id, userId)
}
