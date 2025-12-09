package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateDelete(db *sqlx.DB, entry *EmbryoTransferSave) *apiError.APIError {

	err := hasChildren(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func hasChildren(db *sqlx.DB, entry *EmbryoTransferSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals a
			where birth_date > $1
				and age(birth_date, $1::timestamptz) between interval '240 days' and interval '340 days'
				and father_id = $2
				and mother_id = $3
				and not exists (
					select 1
					from pregnancies_test t
					where t.pregnancy_status = 'FAILED'
						and t.animal_id = $3
						and t.test_date between $1 and a.birth_date
				)
				and user_id = $4
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.TransferDate, entry.BullId, entry.ReceiverId, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(
			"Há uma parição relacionada a esta transferência embrionária! " +
			"Antes de excluir a transferência, altere o pai da parição ou altere os exames de toque.", 
		)
	}

	return nil
}
