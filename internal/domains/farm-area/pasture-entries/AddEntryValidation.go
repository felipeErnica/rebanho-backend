package pastureEntries

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateAddEntry(db *sqlx.DB, entry PastureEntry) *log.APIError {

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

func validateTransferEntry(db *sqlx.DB, entry PastureEntry) *log.APIError {

	err := invalidEntryDate(db, entry)
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

func entryExists(db *sqlx.DB, entry PastureEntry) *log.APIError {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM pasture_entries
			WHERE entry_date = $1
				AND animal_id = $2
				AND pasture_id = $3
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.EntryDate, entry.AnimalId, entry.PastureId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIWarning("Esta entrada já existe! Deseja substituí-la?")
	}

	return nil
}

func samePasture(db *sqlx.DB, entry PastureEntry) *log.APIError {
	query := `
		SELECT pasture_id = $1 AS same_pasture
		FROM pasture_entries
		WHERE animal_id = $2 AND deleted_at IS NULL
		ORDER BY entry_date DESC
		LIMIT 1
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.PastureId, entry.AnimalId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError("O animal já esta neste Lote. A transferência é desnecessária.")
	}

	return nil

}

func invalidEntryDate(db *sqlx.DB, entry PastureEntry) *log.APIError {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM pasture_entries 
			WHERE animal_id = $1
				AND entry_date BETWEEN $2 AND exit_date
				AND deleted_at IS NULL
				AND user_id = $3
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError(`
			A data de entrada informada está em conflito com a data de saída
			do Lote anterior. A data de saída anterior é maior que a data de entrada informada! 
		`)
	}

	return nil
}

func invalidEmptyExitDate(db *sqlx.DB, entry PastureEntry) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM pasture_entries 
			WHERE animal_id = $1
				AND exit_date IS NULL
				AND deleted_at IS NULL
				AND user_id = $2
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError(`
			Não é possível adicionar uma entrada nesta data, pois existe uma entrada em aberto!
		`)
	}

	return nil
}

func invalidExitDate(db *sqlx.DB, entry PastureEntry) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM pasture_entries
			WHERE animal_id = $1
				AND start_date > $2
				AND start_date <= $3
				AND deleted_at IS NULL
				AND user_id = $4
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.ExitDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError(`
			A data de saída informada está em conflito com a data de entrada
			de uma entrada posterior. A data de entrada posterior é menor que a data de saída informada! 
		`)
	}

	return nil
}

func invalidDates(entry PastureEntry) *log.APIError {

	if entry.ExitDate == nil {
		return nil
	}

	if entry.EntryDate.After(*entry.ExitDate) {
		return log.IncorrectEntityAPIError("A data final não pode ser maior que a inicial!")
	}

	return nil
}

func validateTransferCalf(db *sqlx.DB, entry PastureEntry) *log.APIError {

	err := invalidCalfEntry(db, entry)
	if err != nil {
		return err
	}

	if entry.ExitDate != nil {
		err = invalidCalfExit(db, entry)
		if err != nil {
			return err
		}

		err = invalidCalfDates(entry)
		if err != nil {
			return err
		}
	}

	return nil

}

func cancelChangeCalf(db *sqlx.DB, entry PastureEntry) (bool, *log.APIError) {

	weaningQuery := `
		SELECT EXISTS(
			SELECT 1
			FROM animals
			WHERE id = $1 
				AND deleted_at IS NULL 
				AND animal_type <> 'OFFSPRING'
		)
	`

	var isNotWeaningCalf bool
	err := util.GetPrimitive(db, weaningQuery, &isNotWeaningCalf, entry.AnimalId)
	if err != nil {
		return false, log.InternalServerAPIError(err)
	}

	query := `
		SELECT pasture_id = $1 AS same_pasture
		FROM pasture_entries
		WHERE animal_id = $2 AND deleted_at IS NULL
		ORDER BY entry_date DESC
		LIMIT 1
	`

	var samePasture bool
	err = util.GetPrimitive(db, query, &samePasture, entry.PastureId, entry.AnimalId)
	if err != nil {
		return false, log.InternalServerAPIError(err)
	}

	return samePasture || isNotWeaningCalf, nil
}

func invalidCalfEntry(db *sqlx.DB, entry PastureEntry) *log.APIError {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM pasture_entries 
			WHERE animal_id = $1
				AND entry_date BETWEEN $2 AND exit_date
				AND deleted_at IS NULL
				AND user_id = $3
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError(`
			A data de entrada do bezerro(a) está em conflito com a data de saída anterior. 
			A data de saída anterior é maior que a data de entrada informada! 
		`)
	}

	return nil
}

func invalidCalfExit(db *sqlx.DB, entry PastureEntry) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM pasture_entries
			WHERE animal_id = $1
				AND start_date > $2
				AND start_date <= $3
				AND deleted_at IS NULL
				AND user_id = $4
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.AnimalId, entry.EntryDate, entry.ExitDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError(`
			A data de saída do bezerro(a) está em conflito com a data de entrada
			de uma entrada posterior. A data de entrada posterior é menor que a data de saída informada! 
		`)
	}

	return nil
}

func invalidCalfDates(entry PastureEntry) *log.APIError {

	if entry.ExitDate == nil {
		return nil
	}

	if entry.EntryDate.After(*entry.ExitDate) {
		return log.IncorrectEntityAPIError("A data de saída do bezerro(a) não pode ser menor que a de entrada!")
	}

	return nil
}
