package birth

import (
	"math"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type BirthService struct {
	Repo *BirthRepository
}

func NewBirthService(repo *BirthRepository) *BirthService {
	return &BirthService{repo}
}

func (s *BirthService) toDTO(entry BirthDB) BirthDTO {
	dto := BirthDTO{
		Calf: Calf{
			Id:          entry.CalfId,
			Name:        entry.CalfName,
			Tag:         entry.CalfTag,
			BirthDate:   entry.CalfBirthDate,
			Sex:         entry.CalfSex,
			Observation: entry.CalfObservation,
		},
		Mother: Parent{
			Id:   entry.MotherId,
			Name: entry.MotherName,
			Tag:  entry.MotherTag,
		},
		BirthInterval: entry.BirthInterval,
	}

	if entry.FatherId != nil {
		dto.Father = &Parent{
			Id:   *entry.FatherId,
			Name: entry.FatherName,
			Tag:  entry.FatherTag,
		}
	}

	return dto
}

func (s *BirthService) AddBirth(entry *BirthEntrySave) *log.APIError {
	validation, err := s.Repo.CheckBirthConflicts(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if validation.BirthExists && !entry.Overwrite {
		return log.ConflictAPIWarning("Este nascimento já existe. Deseja substitui-lo?")
	}

	if validation.RingExists && !entry.IgnoreTag {
		return log.NewAPIWarning(
			"Brinco já existe!",
			"Este brinco já existe. Deseja adicionar mesmo assim?",
			"TagWarning",
		)
	}

	if validation.InvalidPreviousBirth {
		return log.IncorrectEntityAPIError("O intervalo em relação ao nascimento anterior é muito pequeno. O intervalo deve ser maior que 240 dias!")
	}

	if validation.InvalidNextBirth {
		return log.IncorrectEntityAPIError("O intervalo em relação ao nascimento posterior é muito pequeno. O intervalo deve ser maior que 240 dias!")
	}

	if err := s.Repo.AddBirth(entry); err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *BirthService) UpdateBirth(entry *BirthEntrySave) (*BirthDB, *log.APIError) {
	validation, err := s.Repo.CheckBirthConflicts(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if validation.BirthExists {
		return nil, log.ConflictAPIError("Este nascimento já existe!")
	}

	if validation.InvalidPreviousBirth {
		return nil, log.IncorrectEntityAPIError("O intervalo em relação ao nascimento anterior é muito pequeno. O intervalo deve ser maior que 240 dias!")
	}

	if validation.InvalidNextBirth {
		return nil, log.IncorrectEntityAPIError("O intervalo em relação ao nascimento posterior é muito pequeno. O intervalo deve ser maior que 240 dias!")
	}

	if validation.RingExists && !entry.IgnoreTag {
		return nil, log.NewAPIWarning(
			"Brinco já existe!",
			"Este brinco já existe. Deseja salvar mesmo assim?",
			"TagWarning",
		)
	}

	res, err := s.Repo.UpdateBirth(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return res, nil
}

func (s *BirthService) GetPotentialFather(entry *BirthEntrySave) (*BirthEntrySave, *log.APIError) {
	fatherId, err := s.Repo.GetPotentialFather(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if fatherId != "" {
		entry.FatherId = &fatherId
	}

	return entry, nil
}

func (s *BirthService) GetBirthIntervalHistory(userId string) (*util.CardStats, error) {
	intervalHist, err := s.Repo.GetBirthIntervalHistory(userId)
	if err != nil {
		return nil, err
	}

	var currentInterval, previousInterval, intervalTrend float64
	hist := *intervalHist

	switch length := len(hist); length {
	case 0:
		currentInterval = 0
		previousInterval = 0
		intervalTrend = 0
	case 1:
		currentInterval = hist[length-1].Value
		previousInterval = 0
		intervalTrend = 0
	default:
		currentInterval = hist[length-1].Value
		previousInterval = hist[length-2].Value
		intervalTrend = ((currentInterval / previousInterval) - 1) * 100
	}

	// Handle NaN if previousInterval is 0
	if math.IsNaN(intervalTrend) {
		intervalTrend = 0
	}

	card := util.NewCardStats(hist, intervalTrend, currentInterval)
	return card, nil
}

func (s *BirthService) GetLastBirthsNumber(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetLastBirthsNumber(userId)
	if err != nil {
		return nil, err
	}

	hist := *results
	var current, previous, trend float64

	switch length := len(hist); length {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[length-1].Value
		previous = 0
		trend = 0
	default:
		current = hist[length-1].Value
		previous = hist[length-2].Value
		trend = current - previous
	}

	card := util.NewCardStats(hist, trend, current)
	return card, nil
}

func (s *BirthService) GetYearBirthsNumber(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetYearBirthsNumber(userId)
	if err != nil {
		return nil, err
	}

	hist := *results
	var current, previous, trend float64

	switch length := len(hist); length {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[length-1].Value
		previous = 0
		trend = 0
	default:
		current = hist[length-1].Value
		previous = hist[length-2].Value
		trend = util.CalculatePercentageTrend(current, previous)
	}

	card := util.NewCardStats(hist, trend, current)
	return card, nil
}

func (s *BirthService) GetYearDeathsNumber(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetYearDeathsNumber(userId)
	if err != nil {
		return nil, err
	}

	hist := *results
	var current, previous, trend float64

	switch length := len(hist); length {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[length-1].Value
		previous = 0
		trend = 0
	default:
		current = hist[length-1].Value
		previous = hist[length-2].Value
		trend = util.CalculatePercentageTrend(current, previous)
	}

	card := util.NewCardStats(hist, trend, current)
	return card, nil
}

func (s *BirthService) GetDeathIndex(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetDeathIndex(userId)
	if err != nil {
		return nil, err
	}

	indexHist := *results
	var currentIndex, previousIndex, indexTrend float64

	switch length := len(indexHist); length {
	case 0:
		currentIndex = 0
		previousIndex = 0
		indexTrend = 0
	case 1:
		currentIndex = indexHist[length-1].Value
		previousIndex = 0
		indexTrend = 0
	default:
		currentIndex = indexHist[length-1].Value
		previousIndex = indexHist[length-2].Value
		indexTrend = ((currentIndex / previousIndex) - 1) * 100
	}

	if math.IsNaN(indexTrend) {
		indexTrend = 0
	}

	card := util.NewCardStats(indexHist, indexTrend, currentIndex)
	return card, nil
}

func (s *BirthService) FindPage(
	userId string, 
	sort string, 
	order string, 
	filter *BirthEntryFilter, 
	cursor string, 
	limit int,
) (*util.Page[BirthDTO], error) {
	sort = util.AddNewFields(sort, "id")
	list, err := s.Repo.FindPage(userId, sort, order, filter, cursor, limit)
	if err != nil {
		return nil, err
	}

	newCursor := util.CreateCursorKey(sort, *list)

	listDto := make([]BirthDTO, 0)
	for _, entry := range *list {
		dto := s.toDTO(entry)
		listDto = append(listDto, dto)
	}

	page := util.NewPage(listDto, newCursor, limit)
	return page, err
}
