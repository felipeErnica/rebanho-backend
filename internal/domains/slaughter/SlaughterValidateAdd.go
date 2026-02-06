package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateAdd(db *sqlx.DB, entry *SlaughterEntrySave) *log.APIError {

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

func slaughterExists(db *sqlx.DB, entry *SlaughterEntrySave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM slaughter_entries
			WHERE animal_id = $1 
				AND entry_date = $2
				AND user_id = $3
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIWarning("Já existe um registro de abate deste animal nesta data! Deseja substituir por este?")
	}

	return nil
}

func hasDeathDate(db *sqlx.DB, entry *SlaughterEntrySave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM animals
			WHERE id = $1 
				AND death_date <> $2
				AND user_id = $3
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError("O animal consta como morto em uma data diferente!")
	}

	return nil
}
