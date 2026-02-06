package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateDelete(db *sqlx.DB, entry *EmbryoTransferSave) *log.APIError {

	err := hasChildren(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func hasChildren(db *sqlx.DB, entry *EmbryoTransferSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM animals a
			WHERE birth_date > $1
				AND age(birth_date, $1::timestamptz) BETWEEN INTERVAL '240 days' AND INTERVAL '340 days'
				AND father_id = $2
				AND mother_id = $3
				AND NOT EXISTS (
					SELECT 1
					FROM pregnancies_test t
					WHERE t.pregnancy_status = 'FAILED'
						AND t.animal_id = $3
						AND t.test_date BETWEEN $1 AND a.birth_date
				)
				AND user_id = $4
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.TransferDate, entry.BullId, entry.ReceiverId, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError(
			"Há uma parição relacionada a esta transferência embrionária! " +
				"Antes de excluir a transferência, altere o pai da parição ou altere os exames de toque.",
		)
	}

	return nil
}
