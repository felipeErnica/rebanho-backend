package animals

import (
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type AnimalService struct {
	Repo *AnimalRepository
}

func NewService(repository *AnimalRepository) *AnimalService {
	return &AnimalService{Repo: repository}
}

const DELETE_OBSERVATION = "\n\nOBS.: A exclusão de animais só é recomendada em caso de erros de registro. " +
	"Em caso de morte e/ou abate, faça o registro apropriado. " +
	"Lembre-se, apagar um animal do sistema apagará, PERMANENTEMENTE, todas as informações ligadas a ele."

func (s *AnimalService) toDTO(entry AnimalDB) AnimalDTO {

	dto := AnimalDTO{
		Id:                   entry.Id,
		Name:                 entry.Name,
		Tag:                  entry.Tag,
		Sex:                  entry.Sex,
		AnimalType:           entry.AnimalType,
		BirthDate:            entry.BirthDate,
		DeathDate:            entry.DeathDate,
		WeightBirth:          entry.WeightBirth,
		WeaningDate:          entry.WeaningDate,
		IsBreedingBull:       entry.IsBreedingBull,
		IsTransferBull:       entry.IsTransferBull,
		IsInseminationBull:   entry.IsInseminationBull,
		IsEmbryoDonor:        entry.IsEmbryoDonor,
		IsOutsideAnimal:      entry.IsOutsideAnimal,
		AverageProd:          entry.AverageProd,
		AverageLacInterval:   entry.AverageLacInterval,
		AverageBirthInterval: entry.AverageBirthInterval,
		AveragePeak:          entry.AveragePeak,
		ChildrenNumber:       entry.ChildrenNumber,
		Observation:          entry.Observation,
	}

	if entry.FatherId != nil {
		dto.Father = &Parent{
			Id:   *entry.FatherId,
			Name: entry.FatherName,
			Tag:  entry.FatherTag,
		}
	}

	if entry.MotherId != nil {
		dto.Mother = &Parent{
			Id:   *entry.MotherId,
			Name: entry.MotherName,
			Tag:  entry.MotherTag,
		}
	}

	if entry.PastureId != nil {
		dto.Pasture = &Pasture{
			Id:   *entry.PastureId,
			Name: *entry.PastureName,
			Farm: Farm{
				Id:   *entry.FarmId,
				Name: *entry.FarmName,
			},
		}
	}

	return dto

}

func (s *AnimalService) listToDTO(list []AnimalDB) []AnimalDTO {
	listDTO := make([]AnimalDTO, 0)
	for _, entry := range list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return listDTO
}

func (s *AnimalService) GetDairyHist(userId string) (*util.CardStats, error) {
	hist, err := s.Repo.GetDairyHist(userId)
	if err != nil {
		return nil, err
	}
	return util.NewCardPercentage(*hist), nil
}

func (s *AnimalService) GetBirthHist(userId string) (*util.CardStats, error) {
	hist, err := s.Repo.GetBirthHist(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardInt(*hist), nil
}

func (s *AnimalService) GetDeathHist(userId string) (*util.CardStats, error) {
	hist, err := s.Repo.GetDeathHist(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardInt(*hist), nil
}

func (s *AnimalService) GetSlaughterHist(userId string) (*util.CardStats, error) {
	hist, err := s.Repo.GetSlaughterHist(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*hist), nil
}

func (s *AnimalService) GetAnimalTypes(userId string) (*AnimalByType, error) {
	return s.Repo.GetAnimalTypes(userId)
}

func (s *AnimalService) GetLastDeaths(userId string) (*[]AnimalDTO, error) {
	list, err := s.Repo.GetLastDeaths(userId)
	if err != nil {
		return nil, err
	}

	listDto := s.listToDTO(*list)
	return &listDto, nil
}

func (s *AnimalService) GetAgeAndSex(userId string) (*[]AnimalsByAge, error) {
	return s.Repo.GetAgeAndSex(userId)
}

func (s *AnimalService) FindPage(
	userId string,
	cursor string,
	sort string,
	order string,
	limit int,
	filter *AnimalFilter,
) (*util.Page[AnimalDTO], error) {
	sort = util.AddCommonFields(sort)
	list, err := s.Repo.FindPage(userId, cursor, sort, order, limit, filter)
	if err != nil {
		return nil, err
	}

	newCursor, err := util.CreateCursorKey(sort, *list)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(*list)

	page := util.NewPage(listDTO, newCursor, limit)
	return page, nil
}

func (s *AnimalService) GetPageFoot(userId string, filter *AnimalFilter) (*AnimalFoot, error) {
	return s.Repo.GetPageFoot(userId, filter)
}

func (s *AnimalService) FindById(id string, userId string) (*AnimalDTO, error) {
	entry, err := s.Repo.FindById(id, userId)
	if err != nil {
		return nil, err
	}

	dto := s.toDTO(*entry)
	return &dto, nil
}

func (s *AnimalService) Search(
	sort string,
	order string,
	filter *AnimalFilter,
	userId string,
) (*[]AnimalDTO, error) {

	list, err := s.Repo.Search(sort, order, filter, userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(*list)
	return &listDTO, nil
}

func (s *AnimalService) DeleteValidation(skipValidation bool, id string, userId string) *log.APIError {

	if !skipValidation {
		errResult, err := s.Repo.CheckDeleteErrorConditions(id, userId)
		if err != nil {
			return log.InternalServerAPIError(err)
		}

		messages := []string{}
		if errResult.HasChildren {
			messages = append(messages, "O animal possui filhos registrados no sistema.")
		}

		if errResult.IsCalfInLactation {
			messages = append(messages, "O animal está ligado, como bezerro, a uma lactação.")
		}

		if errResult.IsEmbryoDonor {
			messages = append(messages, "A vaca é uma doadora de embriões.")
		}

		if errResult.IsInseminationBull {
			messages = append(messages, "O touro possui registros ativos de inseminação.")
		}

		if errResult.IsBreedingBull {
			messages = append(messages, "O touro possui registros ativos de cobertura.")
		}

		if errResult.IsEmbryoBull {
			messages = append(messages, "O touro possui registros ativos de transferência embrionária.")
		}

		if len(messages) != 0 {
			warnMsg := "O animal não pode ser excluído, devido aos seguintes motivos:"
			formatedMsg := []string{}
			for i, message := range messages {
				msg := fmt.Sprintf("%d - %s", i+1, message)
				formatedMsg = append(formatedMsg, msg)
			}
			resultMsg := strings.Join(formatedMsg, "\n")
			return log.DeleteAPIError(warnMsg + "\n" + resultMsg + DELETE_OBSERVATION)
		}

		warnResult, err := s.Repo.CheckDeleteWarningConditions(id, userId)
		if err != nil {
			return log.InternalServerAPIError(err)
		}

		records := make([]string, 0)
		if warnResult.HasLactation {
			records = append(records, "lactações")
		}

		if warnResult.HasBreeding {
			records = append(records, "coberturas")
		}

		if warnResult.HasInsemination {
			records = append(records, "inseminações")
		}

		if warnResult.HasSlaughter {
			records = append(records, "abate")
		}

		if warnResult.HasTransfer {
			records = append(records, "receptação de embrião")
		}

		if warnResult.HasBullPastureEntries {
			records = append(records, "entradas em pasto")
		}

		if len(records) != 0 {
			recordStr := strings.Join(records, ", ")
			warnMsg := fmt.Sprintf("Este animal possui importantes registros de: %s. Deseja exclui-lo mesmo assim?", recordStr)
			return log.DeleteWarning(warnMsg + DELETE_OBSERVATION)
		}

	}

	return nil
}

func (s *AnimalService) Delete(skipValidation bool, id string, userId string) *log.APIError {
	err := s.DeleteValidation(skipValidation, id, userId)
	if err != nil {
		return nil
	}
	return s.Repo.Delete(id, userId)
}

func (s *AnimalService) Update(newEntry *AnimalSave) (*AnimalDTO, *log.APIError) {

	res, err := s.Repo.CheckSaveConflicts(newEntry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if res.NumberExists {
		return nil, log.ConflictAPIError("Já há um animal vivo com este brinco. Altere o brinco antes de continuar!")
	}

	if res.NameExists {
		return nil, log.ConflictAPIError("Já há um animal vivo com este nome. Altere o nome antes de continuar!")
	}

	if !newEntry.IgnoreDead {
		warnings := []string{}
		if res.DeadNameExists {
			warnings = append(warnings, "\nHá um animal morto com o mesmo nome.")
		}

		if res.DeadNumberExists {
			warnings = append(warnings, "\nHá um animal morto com o mesmo brinco.")
		}

		if len(warnings) > 0 {
			warning := strings.Join(warnings, "")
			msg := fmt.Sprintf("Os seguintes conflitos foram detectados: %s \nDeseja continuar?", warning)
			return nil, log.NewAPIWarning("Informações já existem!", msg, "IgnoreWarning")
		}
	}

	response, err := s.Repo.Update(newEntry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	dto := s.toDTO(*response)
	return &dto, nil
}

func (s *AnimalService) Add(newEntry *AnimalSave) *log.APIError {

	res, err := s.Repo.CheckSaveConflicts(newEntry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if res.NumberExists {
		return log.ConflictAPIError("Já há um animal vivo com este brinco. Altere o brinco antes de continuar!")
	}

	if res.NameExists {
		return log.ConflictAPIError("Já há um animal vivo com este nome. Altere o nome antes de continuar!")
	}

	if !newEntry.IgnoreDead {
		warnings := []string{}
		if res.DeadNameExists {
			warnings = append(warnings, "\nHá um animal morto com o mesmo nome.")
		}

		if res.DeadNumberExists {
			warnings = append(warnings, "\nHá um animal morto com o mesmo brinco.")
		}

		if len(warnings) > 0 {
			warning := strings.Join(warnings, "")
			msg := fmt.Sprintf("Os seguintes conflitos foram detectados: %s \nDeseja continuar?", warning)
			return log.NewAPIWarning("Informações já existem!", msg, "IgnoreWarning")
		}
	}

	err = s.Repo.Add(newEntry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil

}
