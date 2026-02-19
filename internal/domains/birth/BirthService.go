package birth

import (
	"encoding/json"
	"fmt"

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

func (s *BirthService) listToDTO(list []BirthDB) []BirthDTO {
	listDTO := make([]BirthDTO, 0)
	for _, entry := range list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return listDTO
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
	card := util.NewCardPercentage(*intervalHist)
	return card, nil
}

func (s *BirthService) GetLastBirthsNumber(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetLastBirthsNumber(userId)
	if err != nil {
		return nil, err
	}
	card := util.NewCardInt(*results)
	return card, nil
}

func (s *BirthService) GetYearBirthsNumber(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetYearBirthsNumber(userId)
	if err != nil {
		return nil, err
	}
	card := util.NewCardPercentage(*results)
	return card, nil
}

func (s *BirthService) GetYearDeathsNumber(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetYearDeathsNumber(userId)
	if err != nil {
		return nil, err
	}
	card := util.NewCardPercentage(*results)
	return card, nil
}

func (s *BirthService) GetDeathIndex(userId string) (*util.CardStats, error) {
	results, err := s.Repo.GetDeathIndex(userId)
	if err != nil {
		return nil, err
	}
	card := util.NewCardPercentage(*results)
	return card, nil
}

func (s *BirthService) GetLastBirths(userId string) (*[]BirthDTO, error) {
	list, err := s.Repo.GetLastBirths(userId)
	if err != nil {
		return nil, err
	}
	listDTO := s.listToDTO(*list)
	return &listDTO, nil
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

	listDto := s.listToDTO(*list)
	json, _ := json.MarshalIndent(listDto, "", "	")
	fmt.Println(string(json))
	page := util.NewPage(listDto, newCursor, limit)
	return page, err
}
