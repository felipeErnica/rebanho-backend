package milk

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"time"
)

type MilkService struct {
	Repo *MilkRepository
}

func NewService(repo *MilkRepository) *MilkService {
	return &MilkService{repo}
}

func (s *MilkService) toDTO(entry MilkDB) MilkDTO {
	dto := MilkDTO{
		Id:        entry.Id,
		EntryDate: entry.EntryDate,
		Quantity:  entry.Quantity,
		Cow: Cow{
			Id:   entry.AnimalId,
			Name: entry.AnimalName,
			Tag:  entry.AnimalTag,
		},
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

func (s *MilkService) listToDTO(list *[]MilkDB) *[]MilkDTO {
	listDTO := make([]MilkDTO, 0)
	for _, entry := range *list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return &listDTO
}

func (s *MilkService) GetLastGroups(userId string) (*[]LactationGroup, error) {
	return s.Repo.GetLastGroups(userId)
}

func (s *MilkService) GetLastEntries(userId string) (*[]MilkDTO, error) {
	list, err := s.Repo.GetLastEntries(userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *MilkService) GetMilkProduction(userId string) (*[]util.GraphData, error) {
	return s.Repo.GetMilkProduction(userId)
}

func (s *MilkService) GetLastMilk(userId string) (*util.CardStats, error) {
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
		current = (*averageHist)[0].Value
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].Value
		previous = (*averageHist)[lenght-2].Value
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &util.CardStats{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *MilkService) GetYearMilk(userId string) (*util.CardStats, error) {
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
		current = (*averageHist)[0].Value
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].Value
		previous = (*averageHist)[lenght-2].Value
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &util.CardStats{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *MilkService) GetLastAverageMilk(userId string) (*util.CardStats, error) {
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
		current = (*averageHist)[0].Value
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].Value
		previous = (*averageHist)[lenght-2].Value
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &util.CardStats{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *MilkService) GetYearAverageMilk(userId string) (*util.CardStats, error) {
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
		current = (*averageHist)[0].Value
		previous = 0
		trend = 0
	default:
		current = (*averageHist)[lenght-1].Value
		previous = (*averageHist)[lenght-2].Value
		trend = util.CalculatePercentageTrend(current, previous)
	}

	return &util.CardStats{
		Current: current,
		Trend:   trend,
		Hist:    averageHist,
	}, nil
}

func (s *MilkService) FindPage(
	filter *MilkEntryFilter,
	sort string,
	order string,
	cursor string,
	limit int,
	userId string,
) (*util.Page[MilkDTO], error) {
	sort = util.AddCommonFields(sort)
	list, err := s.Repo.FindPage(filter, sort, order, cursor, limit, userId)
	if err != nil {
		return nil, err
	}

	newCursor := util.CreateCursorKey(sort, *list)
	listDTO := s.listToDTO(list)
	page := util.NewPage(*listDTO, newCursor, limit)
	return page, nil
}

func (s *MilkService) GetPageFoot(filter *MilkEntryFilter, userId string) (*MilkEntryFoot, error) {
	return s.Repo.GetPageFoot(filter, userId)
}

func (s *MilkService) FindGroupsPage(
	filter *LactationGroupFilter, 
	order string, 
	cursor string,
	limit int,
	userId string,
) (*util.Page[LactationGroup], error) {
	list, err := s.Repo.FindGroupsPage(filter, order, cursor, limit, userId)
	if err != nil {
		return nil, err
	}

	newCursor := util.CreateCursorKey("entry_date", *list)
	page := util.NewPage(*list, newCursor, limit)

	return page, err
}

func (s *MilkService) GetGroupEntries(userId string, entryDate time.Time) (*[]MilkDTO, error) {
	list, err :=  s.Repo.GetGroupEntries(userId, entryDate)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *MilkService) GetGroupEntriesFoot(userId string, entryDate time.Time) (*MilkEntryFoot, error) {
	return s.Repo.GetGroupEntriesFoot(userId, entryDate)
}

func (s *MilkService) GetLactationEntriesFoot(lacId string) (*MilkEntryFoot, error) {
	return s.Repo.GetLactationEntriesFoot(lacId)
}

func (s *MilkService) GetLactationEntries(lacId string) (*[]MilkDTO, error) {
	list, err :=  s.Repo.GetLactationEntries(lacId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
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

func (s *MilkService) Update(entry *MilkEntrySave) (*MilkDB, *log.APIError) {

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

func (s *MilkService) Delete(id string) error {
	return s.Repo.Delete(id)
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

func (s *MilkService) DeleteGroup(entryDate time.Time, userId string) error {
	return s.Repo.DeleteGroup(entryDate, userId)
}

