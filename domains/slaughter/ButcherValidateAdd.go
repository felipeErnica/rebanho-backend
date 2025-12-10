package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateButcherAdd(db *sqlx.DB, entry *ButcherSave) *apiError.APIError {

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

func nameExistsAdd(db *sqlx.DB, entry *ButcherSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from butchers
			where name = :name
				and cnpj is not null
				and cnpj <> :cnpj
				and user_id = :user_id
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

func cnpjExistsAdd(db *sqlx.DB, entry *ButcherSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from butchers
			where cnpj = :cnpj
				and user_id = :user_id
				and deleted_at is null
		)
	`
	var exists bool
	err := repositoriesUtil.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning("Este CNPJ já existe. Deseja substituir as informações por estas?")
	}

	return nil
}
