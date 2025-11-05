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
