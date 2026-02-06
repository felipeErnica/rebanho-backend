package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateButcherAdd(db *sqlx.DB, entry *ButcherSave) *log.APIError {

	err := nameExistsAdd(db, entry)
	if err != nil {
		return err
	}

	err = cnpjExistsAdd(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func nameExistsAdd(db *sqlx.DB, entry *ButcherSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM butchers
			WHERE name = :name
				AND cnpj IS NOT NULL
				AND cnpj <> :cnpj
				AND user_id = :user_id
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

func cnpjExistsAdd(db *sqlx.DB, entry *ButcherSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM butchers
			WHERE cnpj = :cnpj
				AND user_id = :user_id
				AND deleted_at IS NULL
		)
	`
	var exists bool
	err := util.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIWarning("Este CNPJ já existe. Deseja substituir as informações por estas?")
	}

	return nil
}
