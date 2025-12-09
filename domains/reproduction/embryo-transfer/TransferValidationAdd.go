package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateAdd(db *sqlx.DB, entry *EmbryoTransferSave) *apiError.APIError {

	err := isPregnant(db, entry)
	if err != nil {
		return err
	}

	err = isRepeated(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func isRepeated(db *sqlx.DB, entry *EmbryoTransferSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from embryo_transfer
			where receiver_id = $1
				and transfer_date = $2
				and user_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.ReceiverId, entry.TransferDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning(
			"Esta vaca já está registrada como doadora na mesma data! " +
				"Deseja substituir este registro?",
		)
	}

	return nil
}

func isPregnant(db *sqlx.DB, entry *EmbryoTransferSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from pregnancy_tests ta
			where ta.test_date < $1
				and ta.pregnancy_status = 'SUCCESS'
				and ta.deleted_at is null
				and ta.animal_id = $2
				and ta.user_id = $3
				and age($1::timestamptz, ta.test_date) <= interval '340 days'
				and not exists (
					select 1
					from pregnancy_tests tb
					where tb.pregnancy_status = 'FAILED'
						and tb.test_date between ta.test_date and $1
						and tb.animal_id = $2
						and tb.user_id = $3
						and tb.deleted_at is null
				)
				and not exists (
					select 1
					from animals a
					where a.birth_date between ta.test_date and $1
						and a.mother_id = $2
						and a.user_id = $3
						and a.deleted_at is null
				)
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.TransferDate, entry.ReceiverId, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(
			"Esta vaca consta como prenha e, portanto, não pode ser receptora. Atualize os registros de toque ou selecione outra " +
			"vaca para receptação.",
		)
	}

	return nil

}
