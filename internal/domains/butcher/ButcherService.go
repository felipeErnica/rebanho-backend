package butcher

import "github.com/felipeErnica/rebanho-backend/internal/log"

type ButcherService struct {
	Repo *ButcherRepository
}

func NewService(repo *ButcherRepository) *ButcherService {
	return &ButcherService{repo}
}

func (s *ButcherService) FindAll(userId string) (*[]ButcherEntry, error) {
	return s.Repo.FindAll(userId)
}

func (s *ButcherService) FindById(id string, userId string) (*ButcherEntry, error) {
	return s.Repo.FindById(id, userId)
}

func (s *ButcherService) Add(entry *ButcherSave) *log.APIError {
	validate, err := s.Repo.CheckSaveConflicts(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if validate.NameExists {
		return log.ConflictAPIError("Este nome já existe!")
	}

	if validate.CnpjExists {
		return log.ConflictAPIError("Este CNPJ já existe!")
	}

	if !entry.IgnoreAddress && validate.AddressExists {
		return log.ConflictAPIWarning("Este endereço já existe! Deseja adicionar mesmo assim?")
	}

	err = s.Repo.Add(entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}

func (s *ButcherService) Update(entry *ButcherSave) (*ButcherEntry, *log.APIError) {
	validate, err := s.Repo.CheckSaveConflicts(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	if validate.NameExists {
		return nil, log.ConflictAPIError("Este nome já existe!")
	}

	if validate.CnpjExists {
		return nil, log.ConflictAPIError("Este CNPJ já existe!")
	}

	if !entry.IgnoreAddress && validate.AddressExists {
		return nil, log.ConflictAPIWarning("Este endereço já existe! Deseja adicionar mesmo assim?")
	}

	response, err := s.Repo.Update(entry)
	if err != nil {
		return nil, log.InternalServerAPIError(err)
	}

	return response, nil
}

func (s *ButcherService) Delete(params *ButcherDelete) *log.APIError {
	if !params.Override {
		validate, err := s.Repo.CheckDelete(params)
		if err != nil {
			return log.InternalServerAPIError(err)
		}

		if validate {
			return log.DeleteWarning("Existem registros de abate que serão excluídos com o frigorífico! Deseja continuar?")
		}
	}

	err := s.Repo.Delete(params)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
