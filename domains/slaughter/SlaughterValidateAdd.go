package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateAdd(db *sqlx.DB, entry *SlaughterEntrySave) *apiError.APIError {

	err := hasDeathDate(db, entry)
	if err != nil {
		return err
	}


	err = slaughterExists(db, entry)
	if err != nil {
		return err
	}

	return nil

}

func slaughterExists(db *sqlx.DB, entry *SlaughterEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from slaughter_entries
			where animal_id = $1 
				and entry_date = $2
				and user_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning("Já existe um registro de abate deste animal nesta data! Deseja substituir por este?")
	}

	return nil
}

func hasDeathDate(db *sqlx.DB, entry *SlaughterEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals
			where id = $1 
				and death_date <> $2
				and user_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError("O animal consta como morto em uma data diferente!")
	}

	return nil
}
