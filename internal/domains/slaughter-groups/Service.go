package slaughtergroups

type Service struct {
	Repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo}
}

func (s *Service) toDTO(entry GroupDB) DTO {
	return DTO{
		Id: entry.Id,
		EntryDate: entry.EntrytDate,
		Discount: entry.Discount,
		AnimalsNumber: entry.AnimalsNumber,
		AverageWeight: entry.AverageWeight,
		AverageDeadWeight: entry.AverageDeadWeight,
		AverageRate: entry.AverageRate,
		Butcher: Butcher{
			Id: entry.ButcherId,
			Name: entry.ButcherName,
			Discount: entry.ButcherDiscount,
		},
	}
}

func (s *Service) listToDTO(list *[]GroupDB) *[]DTO {
	listDTO := make([]DTO, 0)
	for _, entry := range *list {
		dto := s.toDTO(entry)
		listDTO = append(listDTO, dto)
	}
	return &listDTO
}

func (s *Service) FindAll(userId string) (*[]DTO, error) {
	list, err := s.Repo.FindAll(userId)
	if err != nil {
		return  nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *Service) FindLast(userId string) (*[]DTO, error) {
	list, err := s.Repo.FindLast(userId)
	if err != nil {
		return nil, err
	}

	listDTO := s.listToDTO(list)
	return listDTO, nil
}

func (s *Service) FindById(id string, userId string) (*DTO, error) {
	entry, err := s.Repo.FindById(id, userId)
	if err != nil {
		return  nil, err
	}

	dto := s.toDTO(*entry)
	return &dto, nil
}
