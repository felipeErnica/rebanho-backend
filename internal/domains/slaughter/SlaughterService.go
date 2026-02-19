package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type SlaughterService struct {
	Repo *SlaughterRepository
}

func NewService(repo *SlaughterRepository) *SlaughterService {
	return &SlaughterService{repo}
}

func (s *SlaughterService) toDTO(entry SlaughterDB) SlaughterDTO {
	dto := SlaughterDTO{
		Id:              entry.Id,
		EntryDate:       entry.EntryDate,
		DiscountRate:    entry.DiscountRate,
		Weight:          entry.Weight,
		DiscountWeight:  entry.DiscountWeight,
		DeadWeight:      entry.DeadWeight,
		PerformanceRate: entry.PerformanceRate,
		Butcher: Butcher{
			Id:       entry.ButcherId,
			Name:     entry.ButcherName,
			Discount: entry.ButcherDiscount,
		},
	}

	if entry.AnimalId != nil {
		dto.Animal = &Animal{
			Id:        *entry.AnimalId,
			Tag:       entry.AnimalTag,
			Name:      entry.AnimalName,
			Sex:       *entry.AnimalSex,
			BirthDate: entry.AnimalBirth,
		}

		if entry.FatherId != nil {
			dto.Animal.Father = &Parent{
				Id:   *entry.FatherId,
				Tag:  entry.FatherTag,
				Name: entry.FatherName,
			}
		}

		if entry.MotherId != nil {
			dto.Animal.Mother = &Parent{
				Id:   *entry.MotherId,
				Tag:  entry.MotherTag,
				Name: entry.MotherName,
			}
		}
	}

	return dto
}

func (s *SlaughterService) listToDTO(list *[]SlaughterDB) *[]SlaughterDTO {
	listDTO := make([]SlaughterDTO, 0)
	for _, entry := range *list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return &listDTO
}

func (s *SlaughterService) toGroupDTO(entry SlaughterGroupDB) SlaughterGroupDTO {
	return SlaughterGroupDTO{
		EntryDate:           entry.EntryDate,
		AnimalsNumber:       entry.AnimalsNumber,
		AverageWeight:       entry.AverageWeight,
		WeightVariation:     entry.WeightVariation,
		AverageDeadWeight:   entry.AverageDeadWeight,
		DeadWeightVariation: entry.DeadWeightVariation,
		AverageRate:         entry.AverageRate,
		RateVariation:       entry.RateVariation,
		Butcher: Butcher{
			Id:       entry.ButcherId,
			Name:     entry.ButcherName,
			Discount: entry.ButcherDiscount,
		},
	}
}

func (s *SlaughterService) listGroupToDTO(list *[]SlaughterGroupDB) *[]SlaughterGroupDTO {
	listDTO := make([]SlaughterGroupDTO, 0)
	for _, entry := range *list {
		dto := s.toGroupDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return &listDTO
}

func (s *SlaughterService) GetLastAverageWeight(userId string) (*util.CardStats, error) {
	result, err := s.Repo.GetLastAverageWeight(userId)
	if err != nil {
		return nil, err
	}
	return util.NewCardPercentage(*result), nil
}

func (s *SlaughterService) GetLastDeadWeight(userId string) (*util.CardStats, error) {
	result, err := s.Repo.GetLastDeadWeight(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*result), nil
}

func (s *SlaughterService) GetLastPerformance(userId string) (*util.CardStats, error) {
	result, err := s.Repo.GetLastPerformance(userId)
	if err != nil {
		return nil, err
	}

	return util.NewCardPercentage(*result), nil
}

func (s *SlaughterService) GetWeightHist(userId string) (*[]WeightHist, error) {
	return s.Repo.GetWeightHist(userId)
}

func (s *SlaughterService) GetRateHist(userId string) (*[]util.GraphData, error) {
	return s.Repo.GetRateHist(userId)
}

func (s *SlaughterService) GetBestFathers(userId string) (*[]TableRatings, error) {
	return s.Repo.GetBestFathers(userId)
}

func (s *SlaughterService) GetBestMothers(userId string) (*[]TableRatings, error) {
	return s.Repo.GetBestMothers(userId)
}

func (s *SlaughterService) GetBestButchers(userId string) (*[]TableRatings, error) {
	return s.Repo.GetBestButchers(userId)
}

func (s *SlaughterService) GetLastEntries(userId string) (*[]SlaughterDTO, error) {
	list, err := s.Repo.GetLastEntries(userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *SlaughterService) GetLastGroups(userId string) (*[]SlaughterGroupDTO, error) {
	list, err := s.Repo.GetLastGroups(userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listGroupToDTO(list)
	return listDTO, nil
}

func (s *SlaughterService) FindPage(
	sort string,
	order string,
	cursor string,
	filter *SlaughterFilter,
	limit int,
	userId string,
) (*util.Page[SlaughterDTO], error) {
	sort = util.AddCommonFields(sort)
	list, err := s.Repo.FindPage(sort, order, cursor, limit, filter, userId)
	if err != nil {
		return nil, err
	}

	newCursor := util.CreateCursorKey(sort, *list)
	listDTO := s.listToDTO(list)
	page := util.NewPage(*listDTO, newCursor, limit)

	return page, nil
}

func (s *SlaughterService) GetPageFoot(filter *SlaughterFilter, userId string) (*SlaughterFoot, error) {
	return s.Repo.GetPageFoot(filter, userId)
}

func (s *SlaughterService) FindButcherPage(
	sort string,
	order string,
	cursor string,
	filter *SlaughterFilter,
	butcherId string,
	limit int,
	userId string,
) (*util.Page[SlaughterDTO], error) {
	sort = util.AddCommonFields(sort)
	list, err := s.Repo.FindButcherPage(sort, order, cursor, filter, butcherId, limit, userId)
	if err != nil {
		return nil, err
	}

	newCursor := util.CreateCursorKey(sort, *list)
	listDTO := s.listToDTO(list)
	page := util.NewPage(*listDTO, newCursor, limit)
	return page, nil
}

func (s *SlaughterService) GetButcherPageFoot(
	filter *SlaughterFilter,
	butcherId string,
	userId string,
) (*SlaughterFoot, error) {
	return s.Repo.GetButcherPageFoot(filter, butcherId, userId)
}

func (s *SlaughterService) FindGroups(order string, userId string) (*[]SlaughterGroupDTO, error) {
	list, err := s.Repo.FindGroups(order, userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listGroupToDTO(list)
	return listDTO, nil
}

func (s *SlaughterService) FindEntries(
	sort string,
	order string,
	filter *SlaughterFilter,
	userId string,
) (*[]SlaughterDTO, error) {
	list, err := s.Repo.FindEntries(sort, order, filter, userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *SlaughterService) GetEntriesFoot(filter *SlaughterFilter, userId string) (*SlaughterFoot, error) {
	return s.Repo.GetEntriesFoot(filter, userId)
}

func (s *SlaughterService) Delete(id string, userId string) *log.APIError {
	err := s.Repo.Delete(id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *SlaughterService) DeleteBatch(ids []string, userId string) *log.APIError {
	err := s.Repo.DeleteBatch(ids, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *SlaughterService) Update(entry *SlaughterSave) (*SlaughterDTO, *log.APIError) {
	validate, err := s.Repo.CheckSave(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if validate.Exists {
		return nil, log.ConflictAPIError("Já existe um registro deste animal nesta data!")
	}

	if validate.HasDeath && !entry.IgnoreDeath {
		return nil, log.ConflictAPIWarning(
			"O animal já está morto! Deseja continuar mesmo assim?" +
				"\nATENÇÃO: A data de morte será alterada para a data do abate.",
		)
	}

	result, err := s.Repo.Update(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	dto := s.toDTO(*result)
	return &dto, nil
}

func (s *SlaughterService) UpdateBatch(batch *SlaughterSaveBatch) *log.APIError {
	err := s.Repo.UpdateBatch(batch)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *SlaughterService) Add(entry *SlaughterSave) *log.APIError {
	validate, err := s.Repo.CheckSave(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if validate.Exists && !entry.Overwrite {
		return log.ConflictAPIWarning("Já existe um registro deste animal nesta data! Deseja substitui-lo por este?")
	}

	if validate.HasDeath && !entry.IgnoreDeath {
		return log.NewAPIWarning(
			"Animal Morto.",
			"O animal já está morto! Deseja continuar mesmo assim?"+
				"\nATENÇÃO: A data de morte será alterada para a data do abate.",
			"death_warning",
		)
	}

	err = s.Repo.Add(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
