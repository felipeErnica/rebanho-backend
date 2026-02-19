package tests

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type TestService struct {
	Repo *TestEntryRepository
}

func NewService(repo *TestEntryRepository) *TestService {
	return &TestService{repo}
}

func (s *TestService) toDTO(entry TestDB) TestDTO {
	dto := TestDTO{
		Id:              entry.Id,
		TestDate:        entry.TestDate,
		BirthForecast:   entry.BirthForecast,
		BirthStatus:     entry.BirthStatus,
		PregnancyStatus: entry.PregnancyStatus,
		Observation:     entry.Observation,
		Cow: Cow{
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
			BirthDate: *entry.CalfBirthDate,
			DeathDate: entry.CalfDeathDate,
		}
	}

	return dto
}

func (s *TestService) listToDTO(list *[]TestDB) *[]TestDTO {
	listDTO := make([]TestDTO, 0)
	for _, entry := range *list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return &listDTO
}

func (s *TestService) GetPregnancyRates(userId string) (*util.CardStats, error) {
	result, err := s.Repo.GetPregnancyRate(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*result), nil
}

func (s *TestService) GetAnimalsNumber(userId string) (*util.CardStats, error) {
	result, err := s.Repo.GetAnimalsNumber(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*result), nil
}

func (s *TestService) GetBirthRates(userId string) (*util.CardStats, error) {
	result, err := s.Repo.GetBirthRate(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*result), nil
}

func (s *TestService) GetTestHist(userId string) (*[]PregnancyTestHist, error) {
	return s.Repo.GetPregnancyTestHist(userId)
}

func (s *TestService) GetLastEntries(userId string) (*[]TestDTO, error) {
	list, err := s.Repo.GetLastEntries(userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *TestService) GetLastGroups(userId string) (*[]TestGroups, error) {
	return s.Repo.GetLastGroups(userId)
}

func (s *TestService) GetNextBirths(userId string) (*[]util.GraphData, error) {
	return s.Repo.GetNextBirths(userId)
}

func (s *TestService) GetRankedResults(rankBy string, userId string) (*[]TestAnimal, error) {
	switch rankBy {
	case "best-results":
		return s.Repo.GetBestResults(userId)
	case "worst-results":
		return s.Repo.GetWorstResults(userId)
	default:
		return nil, fmt.Errorf("A expressão %s não pe válida!", rankBy)
	}
}

func (s *TestService) FindEntriesPage(
	filter *TestFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*util.Page[TestDTO], error) {
	sort = util.AddCommonFields(sort)
	list, err := s.Repo.FindEntriesPage(filter, sort, order, cursor, limit, userId)
	if err != nil {
		return nil, err
	}

	newCursor := util.CreateCursorKey(sort, *list)
	listDTO := s.listToDTO(list)
	page := util.NewPage(*listDTO, newCursor, limit)
	return page, nil
}

func (s *TestService) GetEntriesFoot(filter *TestFilter, userId string) (*TestFoot, error) {
	return s.Repo.GetEntriesFoot(filter, userId)
}

func (s *TestService) FindGroups(userId string) (*[]TestGroups, error) {
	return s.Repo.FindGroups(userId)
}

func (s *TestService) FindEntriesByGroup(
	sort string,
	order string,
	testDate time.Time,
	userId string,
) (*[]TestDTO, error) {
	list, err := s.Repo.FindEntriesByGroup(sort, order, testDate, userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *TestService) GetEntriesByGroupFoot(testDate time.Time, userId string) (*TestFoot, error) {
	return s.Repo.GetEntriesByGroupFoot(testDate, userId)
}

func (s *TestService) Add(entry *TestSave) *log.APIError {
	exists, err := s.Repo.CheckEntryExistence(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIWarning("Já existe um toque desta vaca na mesma data! Deseja substituí-lo?")
	}

	return s.Repo.Add(entry)
}

func (s *TestService) Update(entry *TestSave) (*TestDTO, *log.APIError) {
	exists, err := s.Repo.CheckEntryExistence(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if exists {
		return nil, log.ConflictAPIError("Já existe um toque desta vaca na mesma data!")
	}

	resp, err := s.Repo.Update(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	dto := s.toDTO(*resp)
	return &dto, nil
}

func (s *TestService) UpdateGroup(group *TestGroupSave) (*TestGroups, *log.APIError) {
	exists, err := s.Repo.CheckGroupExistence(group)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if exists {
		return nil, log.ConflictAPIError("Já existem registros de toque nesta data. Para evitar conflitos" +
			" altere a data escolhida ou modifique os registros em conflito!")
	}

	return s.Repo.UpdateGroup(group)
}

func (s *TestService) Delete(id string, userId string) *log.APIError {
	return s.Repo.Delete(id, userId)
}

func (s *TestService) DeleteGroup(testDate time.Time, userId string) *log.APIError {
	return s.Repo.DeleteGroup(testDate, userId)
}
