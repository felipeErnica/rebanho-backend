package pregnancyTests

import (
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

type PregnancyTestService struct {
	Repo *TestEntryRepository
}

func NewService(repo *TestEntryRepository) *PregnancyTestService {
	return &PregnancyTestService{repo}
}

func (s *PregnancyTestService) GetPregnancyRates(userId string) (*CardStats, error) {
	result, err := s.Repo.GetPregnancyRate(userId)

	if err != nil {
		return nil, err
	}
	pregnancyHist := *result
	var current, previous, trend float64

	switch lenght := len(pregnancyHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = pregnancyHist[lenght-1].PregnancyRate
		previous = 0
		trend = 0
	default:
		current = pregnancyHist[lenght-1].PregnancyRate
		previous = pregnancyHist[lenght-2].PregnancyRate
		trend = ((current / previous) - 1) * 100
	}

	stats := &CardStats{
		Trend:   trend,
		Current: current,
		Hist:    pregnancyHist,
	}

	return stats, nil
}

func (s *PregnancyTestService) GetAnimalsNumber(userId string) (*CardStats, error) {
	result, err := s.Repo.GetAnimalsNumber(userId)
	if err != nil {
		return nil, err
	}

	pregnancyHist := *result
	var current, previous, trend float64

	switch lenght := len(pregnancyHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = pregnancyHist[lenght-1].AnimalsNumber
		previous = 0
		trend = 0
	default:
		current = pregnancyHist[lenght-1].AnimalsNumber
		previous = pregnancyHist[lenght-2].AnimalsNumber
		trend = ((current / previous) - 1) * 100
	}

	stats := &CardStats{
		Trend:   trend,
		Current: current,
		Hist:    pregnancyHist,
	}

	return stats, nil
}

func (s *PregnancyTestService) GetBirthRates(userId string) (*BirthStats, error) {
	result, err := s.Repo.GetBirthRate(userId)
	if err != nil {
		return nil, err
	}

	birthHist := *result
	var current, previous, trend float64

	switch lenght := len(birthHist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = birthHist[lenght-1].BirthRate
		previous = 0
		trend = 0
	default:
		current = birthHist[lenght-1].BirthRate
		previous = birthHist[lenght-2].BirthRate
		trend = ((current / previous) - 1) * 100
	}

	stats := &BirthStats{
		Trend:   trend,
		Current: current,
		Hist:    birthHist,
	}

	return stats, nil
}

func (s *PregnancyTestService) GetTestHist(userId string) (*[]PregnancyTestHist, error) {
	return s.Repo.GetPregnancyTestHist(userId)
}

func (s *PregnancyTestService) GetLastEntries(userId string) (*LastEntries, error) {
	return s.Repo.GetLastEntries(userId)
}

func (s *PregnancyTestService) GetLastGroups(userId string) (*[]TestGroups, error) {
	return s.Repo.GetLastGroups(userId)
}

func (s *PregnancyTestService) GetNextBirths(userId string) (*[]NextBirths, error) {
	return s.Repo.GetNextBirths(userId)
}

func (s *PregnancyTestService) GetRankedResults(rankBy string, userId string) (*[]TestAnimal, error) {
	switch rankBy {
	case "best-results":
		return s.Repo.GetBestResults(userId)
	case "worst-results":
		return s.Repo.GetWorstResults(userId)
	default:
		return nil, fmt.Errorf("A expressão %s não pe válida!", rankBy)
	}
}

func (s *PregnancyTestService) FindEntriesPage(
	filter *TestEntryFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[TestEntry], error) {
	return s.Repo.FindEntriesPage(filter, sort, order, cursor, userId)
}

func (s *PregnancyTestService) GetEntriesFoot(filter *TestEntryFilter, userId string) (*TestEntryFoot, error) {
	return s.Repo.GetEntriesFoot(filter, userId)
}

func (s *PregnancyTestService) FindGroups(userId string) (*[]TestGroups, error) {
	return s.Repo.FindGroups(userId)
}

func (s *PregnancyTestService) FindEntriesByGroup(
	sort string,
	order string,
	testDate time.Time,
	userId string,
) (*[]TestEntry, error) {
	return s.Repo.FindEntriesByGroup(sort, order, testDate, userId)
}

func (s *PregnancyTestService) GetEntriesByGroupFoot(testDate time.Time, userId string) (*TestEntryFoot, error) {
	return s.Repo.GetEntriesByGroupFoot(testDate, userId)
}

func (s *PregnancyTestService) Add(entry *TestEntrySave) *log.APIError {
	exists, err := s.Repo.CheckEntryExistence(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIWarning("Já existe um toque desta vaca na mesma data! Deseja substituí-lo?")
	}

	return s.Repo.Add(entry)
}

func (s *PregnancyTestService) Update(entry *TestEntrySave) (*TestEntry, *log.APIError) {
	exists, err := s.Repo.CheckEntryExistence(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if exists {
		return nil, log.ConflictAPIError("Já existe um toque desta vaca na mesma data!")
	}

	return s.Repo.Update(entry)
}

func (s *PregnancyTestService) UpdateBatch(group *TestGroups) (*TestGroups, *log.APIError) {
	exists, err := s.Repo.CheckGroupExistence(group)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if exists {
		return nil, log.ConflictAPIError("Já existem registros de toque nesta data. Para evitar conflitos" +
			" altere a data escolhida ou modifique os registros em conflito!")
	}

	return s.Repo.UpdateBatch(group)
}

func (s *PregnancyTestService) Delete(id string, userId string) *log.APIError {
	return s.Repo.Delete(id, userId)
}

func (s *PregnancyTestService) DeleteBatch(testDate time.Time, userId string) *log.APIError {
	return s.Repo.DeleteBatch(testDate, userId)
}
