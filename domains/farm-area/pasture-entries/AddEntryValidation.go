package pastureEntries

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateAddEntry(db *sqlx.DB, entry PastureEntry) *apiError.APIError {

	err := entryExists(db, entry)
	if err != nil {
		return err
	}

	err = samePasture(db, entry)
	if err != nil {
		return err
	}

	err = invalidEntryDate(db, entry)
	if err != nil {
		return err
	}

	if entry.ExitDate != nil {
		err = invalidExitDate(db, entry)
		if err != nil {
			return err
		}

		err = invalidDates(entry)
		if err != nil {
			return err
		}
	}

	err = invalidEmptyExitDate(db, entry)
	if err != nil {
		return err
	}

	return nil

}

func validateTransferEntry(db *sqlx.DB, entry PastureEntry) *apiError.APIError {

	err := transferExists(db, entry)
	if err != nil {
		return err
	}

	err = invalidEntryDate(db, entry)
	if err != nil {
		return err
	}

	err = samePasture(db, entry)
	if err != nil {
		return err
	}


	if entry.ExitDate != nil {
		err = invalidExitDate(db, entry)
		if err != nil {
			return err
		}

		err = invalidDates(entry)
		if err != nil {
			return err
		}
	}

	return nil

}

func transferExists(db *sqlx.DB, entry PastureEntry) *apiError.APIError {
	query := `
		select exists (
			select 1
			from pasture_entries
			where entry_date = $1
				and animal_id = $2
				and pasture_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.EntryDate, entry.AnimalId, entry.PastureId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Esta entrada já existe!")
	}

	return nil
}

func entryExists(db *sqlx.DB, entry PastureEntry) *apiError.APIError {
	query := `
		select exists (
			select 1
			from pasture_entries
			where entry_date = $1
				and animal_id = $2
				and pasture_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.EntryDate, entry.AnimalId, entry.PastureId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning("Esta entrada já existe! Deseja substituí-la?")
	}

	return nil
}

func samePasture(db *sqlx.DB, entry PastureEntry) *apiError.APIError {
	query := `
		select pasture_id = $1 as same_pasture
		from pasture_entries
		where animal_id = $2 and deleted_at is null
		order by entry_date desc
		limit 1
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.PastureId, entry.AnimalId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError("O animal já esta neste Lote. A transferência é desnecessária.")
	}

	return nil

}

func invalidEntryDate(db *sqlx.DB, entry PastureEntry) *apiError.APIError {
	query := `
		select exists (
			select 1 
			from pasture_entries 
			where animal_id = $1
				and entry_date between $2 and exit_date
				and deleted_at is null
				and user_id = $3
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			A data de entrada informada está em conflito com a data de saída
			do Lote anterior. A data de saída anterior é maior que a data de entrada informada! 
		`)
	}

	return nil
}

func invalidEmptyExitDate(db *sqlx.DB, entry PastureEntry) *apiError.APIError {

	query := `
		select exists (
			select 1 
			from pasture_entries 
			where animal_id = $1
				and exit_date is null
				and deleted_at is null
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			Não é possível adicionar uma entrada nesta data, pois existe uma entrada em aberto!
		`)
	}

	return nil
}

func invalidExitDate(db *sqlx.DB, entry PastureEntry) *apiError.APIError {

	query := `
		select exists (
			select 1 
			from pasture_entries
			where animal_id = $1
				and start_date > $2
				and start_date <= $3
				and deleted_at is null
				and user_id = $4
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.ExitDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(`
			A data de saída informada está em conflito com a data de entrada
			de uma entrada posterior. A data de entrada posterior é menor que a data de saída informada! 
		`)
	}

	return nil
}

func invalidDates(entry PastureEntry) *apiError.APIError {

	if entry.ExitDate == nil {
		return nil
	}

	if entry.EntryDate.After(*entry.ExitDate) {
		return apiError.IncorrectEntityAPIError("A data final não pode ser maior que a inicial!")
	}

	return nil
}
