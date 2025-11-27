package animals

import (
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type InfoValidation struct {
	SameName   bool `db:"same_name"`
	SameNumber bool `db:"same_number"`
}

func validateAdd(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	err := numberExists(db, entry)
	if err != nil {
		return err
	}

	err = isRepeated(db, entry)
	if err != nil {
		return err
	}

	err = infoExists(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func validateReplace(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	err := infoExists(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func isRepeated(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals
			where name = $1
				and death_date is null
				and deleted_at is null
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.Name, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning(
			"Já há um animal vivo com este nome. Ao continuar, os dados deste " +
				"animal serão substituídos! Deseja continuar?",
		)
	}

	return nil
}

func numberExists(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals
			where ring_number = $1
				and death_date is null
				and name is not null
				and deleted_at is null
				and user_id = $2
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.Name, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já há um animal vivo com este brinco. Altere o brinco antes de continuar!")
	}

	return nil
}

func infoExists(db *sqlx.DB, entry *AnimalSave) *apiError.APIError {
	query := `
		select 
			exists (
				select 1
				from animals
				where user_id = $1
					and ring_number = $2
					and name is not null
					and death_date is not null
					and deleted_at is null
			) as same_number,
			exists (
				select 1
				from animals
				where user_id = $1
					and name = $3
					and death_date is not null
					and deleted_at is null
			) as same_name
	`

	res, err := repositoriesUtil.GetOne[InfoValidation](db, query, entry.UserId, entry.RingNumber, entry.Name)
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
