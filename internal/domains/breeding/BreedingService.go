package breeding

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type BreedingService struct {
	Repo *BreedingRepository
}

func NewBreedingService(repo *BreedingRepository) *BreedingService {
	return &BreedingService{repo}
}

func (s *BreedingService) toDTO(entry BreedingDB) BreedingDTO {
	dto := BreedingDTO{
		Id:              entry.Id,
		BreedingDate:    entry.BreedingDate,
		BirthStatus:     entry.BirthStatus,
		PregnancyStatus: entry.PregnancyStatus,
		Observation:     entry.Observation,
		Cow: Parent{
			Id:   entry.AnimalId,
			Name: entry.AnimalName,
			Tag:  entry.AnimalTag,
		},
		Bull: Parent{
			Id:   entry.BullId,
			Name: entry.BullName,
			Tag:  entry.BullTag,
		},
	}

	if entry.ChildId != nil {
		dto.Child = &Child{
			Id:        *entry.ChildId,
			Name:      *entry.ChildName,
			Tag:       entry.ChildTag,
			Sex:       *entry.ChildSex,
			BirthDate: *entry.ChildBirthDate,
		}
	}

	return dto
}

func (s *BreedingService) GetBirthRateStats(userId string) (*util.CardStats, error) {
	birthRates, err := s.Repo.GetBirthRateStats(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*birthRates), nil
}

func (s *BreedingService) GetPregnancyRateStats(userId string) (*util.CardStats, error) {
	pregnancyRates, err := s.Repo.GetPregnancyRateStats(userId)
	if err != nil {
		return nil, err
	}
	return util.NewCardPercentage(*pregnancyRates), nil
}

func (s *BreedingService) GetAnimalsNumber(userId string) (*util.CardStats, error) {
	animalsNumbers, err := s.Repo.GetAnimalsNumber(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*animalsNumbers), nil
}

func (s *BreedingService) GetLastEntries(userId string) (*[]BreedingDTO, error) {
	list, err := s.Repo.GetLastEntries(userId)
	if err != nil {
		return nil, err
	}

	listDTO := make([]BreedingDTO, 0)
	for _, entry := range *list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}

	return &listDTO, nil
}

func (s *BreedingService) AddBreeding(entry *BreedingEntrySave) *log.APIError {
	exists, err := s.Repo.CheckBreedingAdd(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIError("Já existe uma entrada desta vaca na mesma data! Deseja substituí-la por esta?")
	}

	return s.Repo.AddBreeding(entry)
}

func (s *BreedingService) Delete(
	id string,
	ignorePregnancy bool,
	changeFather bool,
	userId string,
) *log.APIError {

	validate, err := s.Repo.CheckBreedingDelete(id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if validate.HasPregnancy && !ignorePregnancy {
		return log.DeleteWarningKind(
			"PregnancyWarning",
			"A vaca possui uma prenhez ligada ao touro desta cobertura. Deseja excluir mesmo assim?",
		)
	}

	if validate.HasChildren && !changeFather {
		return log.DeleteWarningKind(
			"ChildrenWarning",
			"A vaca possui uma cria do touro desta cobertura. Deseja substituir o pai?",
		)
	}

	err = s.Repo.Delete(id, changeFather, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *BreedingService) Update(newEntry *BreedingEntrySave) (*BreedingDB, *log.APIError) {
	validate, err := s.Repo.CheckBreedingSave(newEntry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if validate.Repeated {
		return nil, log.ConflictAPIError("Já existe uma inseminação desta vaca na mesma data!")
	}

	if !newEntry.SkipValidation {
		warnings := make([]string, 0)
		if validate.HasChildren {
			warnings = append(warnings, "A vaca possui uma cria, cujo o pai está ligado a esta cobertura.")
		}

		if validate.HasPregnancy {
			warnings = append(warnings, "A vaca possui uma prenhez ligada a esta cobertura!")
		}

		if len(warnings) != 0 {
			message := util.FormatWarningMessage(warnings...)
			return nil, log.ConflictAPIWarning(message + "\nDeseja alterar mesmo assim?")
		}
	}

	return s.Repo.Update(newEntry)
}

func (s *BreedingService) UpdateGroup(date time.Time, group *BreedingGroup) (*BreedingGroup, *log.APIError) {
	exists, err := s.Repo.CheckEntriesOnDate(group.BreedingDate, group.UserId)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}
	if exists {
		return nil, log.ConflictAPIError("Já existem registros de inseminação nesta mesma data! " +
			"Para evitar conflitos, altere os registros existentes ou escolha outra data.",
		)
	}

	return s.Repo.UpdateGroup(date, group)
}
