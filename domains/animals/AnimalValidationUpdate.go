package animals

import (
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	err := nameConflict(db, entry)
	if err != nil {
		return err
	}

	err = numberConflict(db, entry)
	if err != nil {
		return err
	}

	err = infoExistsUpdate(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func nameConflict(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	if entry.Name != nil {
		return nil
	}

	query := `
		select exists (
			select 1
			from animals
			where name = $1
				and death_date is null
				and id <> $2
				and deleted_at is null
				and user_id = $3
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.Name, entry.Id, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já há um animal vivo com este nome. Altere o nome antes de continuar!")
	}

	return nil
}

func numberConflict(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	if entry.Name != nil {
		return nil
	}

	query := `
		select exists (
			select 1
			from animals
			where ring_number = $1
				and death_date is null
				and name is not null
				and id <> $2
				and deleted_at is null
				and user_id = $3
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.RingNumber, entry.Id, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já há um animal vivo com este brinco. Altere o brinco antes de continuar!")
	}

	return nil
}

func infoExistsUpdate(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	if entry.Name != nil {
		return nil
	}

	query := `
		select exists (
			select 1
			from animals
			where user_id = $1
				and ring_number = $2
				and name is not null
				and death_date is not null
				and id <> $4
				and deleted_at is null
		) as same_number,
		select exists (
			select 1
			from animals
			where user_id = $1
				and name = $3
				and id <> $4
				and death_date is not null
				and deleted_at is null
		) as same_name
	`

	res, err := repositoriesUtil.GetOne[InfoValidation](db, query, entry.UserId, entry.RingNumber, entry.Name, entry.Id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	warnings := []string{}
	if res.SameName {
		warnings = append(warnings, "Há um animal morto com o mesmo nome.")
	}

	if res.SameNumber {
		warnings = append(warnings, "Há um animal morto com o mesmo brinco.")
	}

	if len(warnings) > 0 {
		warning := strings.Join(warnings, "\n")
		msg := fmt.Sprintf("Verifique os avisos antes de continuar: \n%s \nDeseja continuar?", warning)
		return apiError.NewAPIWarning("Informações já existem!", msg, "IgnoreWarning")
	}

	return nil
}
