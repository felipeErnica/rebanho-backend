package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(db *sqlx.DB, entry *EmbryoTransferSave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from embryo_transfer
			where receiver_id = $1
				and transfer_date = $2
				and id <> $3
				and user_id = $4
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.ReceiverId, entry.TransferDate, entry.Id, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Há um registro desta vaca na mesma data! Altere as informações antes de continuar!")
	}

	return nil
}

func validateUpdateGroups(db *sqlx.DB, entry *TransferGroup) *apiError.APIError {

	query := `
		select exists (
			select 1
			from embryo_transfer
			where transfer_date = $1
				and user_id = $2
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.TransferDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já existem registros nesta data! Altere as informações antes de continuar!")
	}

	return nil
}
