package lactation

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateAddLacation(db *sqlx.DB, lac LactationHist) *apiError.APIError {

	err := lacExists(db, lac)
	if err != nil {
		return err
	}

	err = invalidBirthDate(db, lac)
	if err != nil {
		return err
	}

	err = invalidDates(lac)
	if err != nil {
		return err
	}

	err = invalidStartDate(db, lac)
	if err != nil {
		return err
	}

	err = invalidEmptyEndDate(db, lac)
	if err != nil {
		return err
	}

	err = invalidEndDate(db, lac)
	if err != nil {
		return err
	}

	return nil
}

func lacExists(db *sqlx.DB, lac LactationHist) *apiError.APIError {
	query := `
		select exists (
			select 1
			from lactations
			where animal_id = $1 
				and start_date = $2
				and deleted_at is null
				and user_id = $3
		)
	`
	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.AnimalId, lac.StartDate, lac.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Esta lactação já existe! Deseja substitui-lá por esta?")
	}

	return nil
}

func invalidBirthDate(db *sqlx.DB, lac LactationHist) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals a
			where a.id = $1 and a.birth_date >= $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.CalfId, lac.StartDate)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError("A data de nascimento do bezerro não pode ser maior que a data de início!")
	}

	return nil
}

func invalidStartDate(db *sqlx.DB, lac LactationHist) *apiError.APIError {
	query := `
		select exists (
			select 1 
			from lactations l
			where l.animal_id = $1
				and l.start_date < $2
				and l.end_date >= $2
				and l.deleted_at is null
				and l.user_id = $3
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.AnimalId, lac.StartDate, lac.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			A data de início informada está em conflito com a data final
			da lactação anterior. A data final anterior é maior que a data de início informada! 
		`)
	}

	return nil
}

func invalidEmptyEndDate(db *sqlx.DB, lac LactationHist) *apiError.APIError {

	if lac.EndDate != nil {
		return nil
	}

	query := `
		select exists (
			select 1 
			from lactations l
			where l.animal_id = $1
				and l.start_date > $2
				and l.deleted_at is null
				and l.user_id = $3
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.AnimalId, lac.StartDate, lac.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			Não é possível adicionar uma lactação em aberto (sem data final), pois já existe uma lactação posterior!
		`)
	}

	return nil
}

func invalidEndDate(db *sqlx.DB, lac LactationHist) *apiError.APIError {

	if lac.EndDate == nil {
		return nil
	}

	query := `
		select exists (
			select 1 
			from lactations l
			where l.animal_id = $1
				and l.start_date > $2
				and l.start_date <= $3
				and l.deleted_at is null
				and l.user_id = $4
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.AnimalId, lac.StartDate, lac.EndDate, lac.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			A data de fim informada está em conflito com a data de início
			de uma lactação posterior. A data inicial posterior é menor que a data de fim informada! 
		`)
	}

	return nil
}

func invalidDates(lac LactationHist) *apiError.APIError {

	if lac.EndDate == nil {
		return nil
	}

	if lac.StartDate.After(*lac.EndDate) {
		return apiError.IncorrectEntityAPIError("A data final não pode ser maior que a inicial!")
	}

	return nil
}

func validateUpdateLacation(db *sqlx.DB, lac LactationHist) *apiError.APIError {

	err := invalidDates(lac)
	if err != nil {
		return err
	}

	if lac.CalfId != nil {
		err = invalidBirthDate(db, lac)
		if err != nil {
			return err
		}

		err = invalidCalf(db, lac)
		if err != nil {
			return err
		}
	}

	err = invalidUpdateStartDate(db, lac)
	if err != nil {
		return err
	}

	err = invalidUpdateEmptyEndDate(db, lac)
	if err != nil {
		return err
	}

	err = invalidUpdateEndDate(db, lac)
	if err != nil {
		return err
	}

	return nil
}

func invalidCalf(db *sqlx.DB, lac LactationHist) *apiError.APIError {
	query := `
		select exists (
			select 1
			from lactations 	
			where id <> $1
				and calf_id = $2
				and deleted_at is null
				and user_id = $3
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.Id, lac.CalfId, lac.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError("Este animal já está vinculado a outra lactação!")
	}

	return nil
}

func invalidUpdateStartDate(db *sqlx.DB, lac LactationHist) *apiError.APIError {
	query := `
		select exists (
			select 1 
			from lactations l
			where l.animal_id = $1
				and l.start_date < $2
				and l.end_date >= $2
				and l.deleted_at is null
				and l.user_id = $3
				and l.id <> $4
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.AnimalId, lac.StartDate, lac.UserId, lac.Id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			A data de início informada está em conflito com a data final
			da lactação anterior. A data final anterior é maior que a data de início informada! 
		`)
	}

	return nil
}

func invalidUpdateEmptyEndDate(db *sqlx.DB, lac LactationHist) *apiError.APIError {

	if lac.EndDate != nil {
		return nil
	}

	query := `
		select exists (
			select 1 
			from lactations l
			where l.animal_id = $1
				and l.start_date > $2
				and l.deleted_at is null
				and l.user_id = $3
				and l.id = $4
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.AnimalId, lac.StartDate, lac.UserId, lac.Id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			Não é possível adicionar uma lactação em aberto (sem data final), pois já existe uma lactação posterior!
		`)
	}

	return nil
}

func invalidUpdateEndDate(db *sqlx.DB, lac LactationHist) *apiError.APIError {

	if lac.EndDate == nil {
		return nil
	}

	query := `
		select exists (
			select 1 
			from lactations l
			where l.animal_id = $1
				and l.start_date > $2
				and l.start_date <= $3
				and l.deleted_at is null
				and l.user_id = $4
				and l.id = $5
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, lac.AnimalId, lac.StartDate, lac.EndDate, lac.UserId, lac.Id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			A data de fim informada está em conflito com a data de início
			de uma lactação posterior. A data inicial posterior é menor que a data de fim informada! 
		`)
	}

	return nil
}

func validateDeleteLactation(db *sqlx.DB, id string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from milk_entries m
				join lactations l on m.animal_id = l.animal_id
					and m.entry_date between l.start_date and coalesce(l.end_date, now())
			where l.id = $1
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteWarning(`
			Existem marcações de leite relacionadas a esta lactação. Deseja proceder com a exclusão?
			Caso sim, todas as marcações relacionadas com a lactação também serão excluídas!
		`)
	}

	return nil
}
