package pasture

type PastureService struct {
	Repo *PastureRepository
}

func NewService(repo *PastureRepository) *PastureService {
	return &PastureService{repo}
}

func (s *PastureService) toDTO(entry PastureDB) PastureDTO {
	dto := PastureDTO{
		Id:            entry.Id,
		Name:          entry.Name,
		PastureSize:   entry.PastureSize,
		AnimalsNumber: entry.AnimalsNumber,
		Farm: Farm{
			Id:   entry.FarmId,
			Name: entry.Name,
		},
	}

	if entry.BullId != nil {
		dto.Bull = &Bull{
			Id:   *entry.BullId,
			Name: *entry.BullName,
			Tag:  entry.BullTag,
		}
	}

	return dto
}

func (s *PastureService) toDTOList(list []PastureDB) []PastureDTO {
	listDto := make([]PastureDTO, 0)
	for _, entry := range list {
		dto := s.toDTO(entry)
		listDto = append(listDto, dto)
	}
	return listDto
}

func (s *PastureService) Search(filter *PastureFilter, userId string) (*[]PastureDTO, error) {
		list, err := s.Repo.Search(filter, userId)
	if err != nil {
		return nil, err
	}
	listDto := s.toDTOList(*list)
	return &listDto, nil
}

func (s *PastureService) FindAnimalsById(pastureId string, userId string, sort string, order string) (*[]PastureAnimal, error) {
	return s.Repo.FindAnimalsById(pastureId, userId, sort, order)
}
