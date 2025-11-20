package pregnancyTests

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateAdd(db *sqlx.DB, entry *TestEntrySave) *apiError.APIError {
	query := `
		select exists (
			select 1
			from pregnancy_tests
			where test_date = $1
				and animal_id = $2
				and user_id = $3
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.TestDate, entry.AnimalId, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning("Já existe um toque desta vaca na mesma data! Deseja substituí-lo?")
	}

	return nil
}

func validateUpdate(db *sqlx.DB, entry *TestEntrySave) *apiError.APIError {
	query := `
		select exists (
			select 1
			from pregnancy_tests
			where test_date = $1
				and animal_id = $2
				and user_id = $3
				and id <> $4
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.TestDate, entry.AnimalId, entry.UserId, entry.Id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já existe um toque desta vaca na mesma data!")
	}

	return nil
}

func validateUpdateBatch(db *sqlx.DB, group *TestGroups) *apiError.APIError {

	query := `
		select exists (
			select 1
			from pregnancy_tests
			where test_date = $1 
				and deleted_at is null
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, group.TestDate, group.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já existem registros de toque nesta data. Para evitar conflitos" + 
			" altere a data escolhida ou modifique os registros em conflito!")
	}

	return nil 
}


