package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateButcherUpdate(db *sqlx.DB, entry *ButcherSave) *log.APIError {

	err := nameExists(db, entry)
	if err != nil {
		return err
	}

	err = cnpjExists(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func nameExists(db *sqlx.DB, entry *ButcherSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM butchers
			WHERE name = :name
				AND user_id = :user_id
				AND id <> :id
				AND deleted_at IS NULL
		)
	`
	var exists bool
	err := util.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIError("Este nome já existe!")
	}

	return nil
}

func cnpjExists(db *sqlx.DB, entry *ButcherSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM butchers
			WHERE cnpj = :cnpj
				AND user_id = :user_id
				AND id <> :id
				AND deleted_at IS NULL
		)
	`
	var exists bool
	err := util.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIError("Este CNPJ já existe!")
	}

	return nil
}
