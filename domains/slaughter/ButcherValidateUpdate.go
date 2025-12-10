package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateButcherUpdate(db *sqlx.DB, entry *ButcherSave) *apiError.APIError {

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

func nameExists(db *sqlx.DB, entry *ButcherSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from butchers
			where name = :name
				and user_id = :user_id
				and id <> :id
				and deleted_at is null
		)
	`
	var exists bool
	err := repositoriesUtil.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Este nome já existe!")
	}

	return nil
}

func cnpjExists(db *sqlx.DB, entry *ButcherSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from butchers
			where cnpj = :cnpj
				and user_id = :user_id
				and id <> :id
				and deleted_at is null
		)
	`
	var exists bool
	err := repositoriesUtil.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Este CNPJ já existe!")
	}

	return nil
}
