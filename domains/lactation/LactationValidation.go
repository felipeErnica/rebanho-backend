package lactation

import (
	"bytes"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type SaveValidation struct {
	LactationExists  *string `db:"lac_exists"`
	InvalidNew       bool    `db:"invalid_new"`
	InvalidStart     bool    `db:"invalid_start"`
	InvalidEnd       bool    `db:"invalid_end"`
	InvalidCalf      bool    `db:"invalid_calf"`
	InvalidEmptyEnd  bool    `db:"invalid_empty_end"`
	DifferentPasture bool    `db:"different_pasture"`
}

func ValidateSave(db *sqlx.DB, lac LactationHistSave) *apiError.APIError {

	query := `
		select 
			(
				select id
				from lactations 
				where animal_id = :animal_id
					and id is distinct from :id
					and start_date = :start_date
					and user_id = :user_id
					and deleted_at is null
			) as lac_exists,
			exists (
				select 1 
				from lactations
				where :end_date is null
					and id is distinct from :id
					and animal_id = :animal_id
					and end_date is null
					and  user_id = :user_id
					and deleted_at is null
			) as invalid_new,
			exists (
				select 1 
				from lactations l
				where l.animal_id = :animal_id
					and l.id is distinct from :id
					and l.start_date < :start_date
					and l.end_date >= :start_date
					and l.deleted_at is null
					and l.user_id = :user_id
			) as invalid_start,
			exists (
				select 1 
				from lactations l
				where :end_date is null
					and l.animal_id = :animal_id
					and l.id is distinct from :id
					and l.start_date > :start_date
					and l.user_id = :user_id
					and l.deleted_at is null
			) as invalid_empty_end,
			exists (
				select 1 
				from lactations l
				where l.animal_id = :animal_id
					and l.id is distinct from :id
					and l.start_date > :start_date
					and l.start_date <= :end_date
					and l.deleted_at is null
					and l.user_id = :user_id
			) as invalid_end,
			exists (
				select 1
				from lactations 	
				where id is distinct from :id
					and calf_id = :calf_id
					and user_id = :user_id
					and deleted_at is null
			) as invalid_calf,
			(
				select coalesce(pasture_id <> :pasture_id, false)
				from pasture_entries
				where animal_id = :animal_id
					and exit_date is null
					and user_id = :user_id
					and deleted_at is null
			) as different_pasture
	`
	validate, err := repositoriesUtil.NamedGet(db, query, SaveValidation{}, lac)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	apiErrors := make([]string, 0)

	if lac.EndDate != nil && lac.StartDate.After(*lac.EndDate) {
		apiErrors = append(apiErrors, "A data de início não pode ser maior que a data de encerramento!")
	}

	if validate.InvalidStart {
		apiErrors = append(apiErrors,
			"A data de início informada está em conflito com a data de encerramento "+
				"da lactação anterior. A data de início informada é menor que a data de encerramento anterior!",
		)
	}

	if validate.InvalidNew {
		apiErrors = append(apiErrors, "Não é possível adicionar uma nova lactação enquanto a antiga não for encerrada!")
	}

	if validate.InvalidEmptyEnd {
		apiErrors = append(apiErrors, "Não é possível adicionar uma lactação em aberto (sem encerramento), pois já existe uma lactação posterior!")
	}

	if validate.InvalidEnd {
		apiErrors = append(apiErrors,
			"A data de encerramento informada está em conflito com a data de início " +
			"de uma lactação posterior. A data de encerramento informada é maior que a data de início da lactação posterior!",
		)
	}

	if validate.InvalidCalf {
		apiErrors = append(apiErrors, "O bezerro selecionado está vinculado a outra lactação!")
	}

	if len(apiErrors) != 0 {
		var errBuff bytes.Buffer
		for i, msg := range apiErrors {
			errPoint := fmt.Sprintf("\n%d - %s", i + 1, msg)
			errBuff.WriteString(errPoint)
		}
		errMsg := fmt.Sprintf("Os seguintes erros foram encontrados: %s", errBuff.String())
		return apiError.IncorrectEntityAPIError(errMsg)
	}

	if validate.DifferentPasture && !lac.TransferPasture {
		return apiError.NewAPIWarning(
			"Pasto diferente!",
			"A vaca não está no pasto selecionado! Deseja transferí-la?",
			"PastureWarning",
		)
	}

	if validate.LactationExists != nil && lac.Id != nil {
		return apiError.ConflictAPIError("Já existe uma lactação desta vaca na mesma data!")
	}

	if validate.LactationExists != nil && !lac.Overwrite {
		return apiError.ReplaceAPIWarning("Esta lactação já existe! Deseja substituí-la por esta?", validate.LactationExists)
	}

	return nil
}
