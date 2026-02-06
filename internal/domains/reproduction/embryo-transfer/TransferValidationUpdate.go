package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(db *sqlx.DB, entry *EmbryoTransferSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM embryo_transfer
			WHERE receiver_id = $1
				AND transfer_date = $2
				AND id <> $3
				AND user_id = $4
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.ReceiverId, entry.TransferDate, entry.Id, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIError("Há um registro desta vaca na mesma data! Altere as informações antes de continuar!")
	}

	return nil
}

func validateUpdateGroups(db *sqlx.DB, entry *TransferGroupSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM embryo_transfer
			WHERE transfer_date = $1
				AND user_id = $2
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.TransferDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIError("Já existem registros nesta data! Altere as informações antes de continuar!")
	}

	return nil
}
