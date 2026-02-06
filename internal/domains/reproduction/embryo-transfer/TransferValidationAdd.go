package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateAdd(db *sqlx.DB, entry *EmbryoTransferSave) *log.APIError {

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

func isRepeated(db *sqlx.DB, entry *EmbryoTransferSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM embryo_transfer
			WHERE receiver_id = $1
				AND transfer_date = $2
				AND user_id = $3
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.ReceiverId, entry.TransferDate, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIWarning(
			"Esta vaca já está registrada como doadora na mesma data! " +
				"Deseja substituir este registro?",
		)
	}

	return nil
}

func isPregnant(db *sqlx.DB, entry *EmbryoTransferSave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM pregnancy_tests ta
			WHERE ta.test_date < $1
				AND ta.pregnancy_status = 'SUCCESS'
				AND ta.deleted_at IS NULL
				AND ta.animal_id = $2
				AND ta.user_id = $3
				AND age($1::timestamptz, ta.test_date) <= INTERVAL '340 days'
				AND NOT EXISTS (
					SELECT 1
					FROM pregnancy_tests tb
					WHERE tb.pregnancy_status = 'FAILED'
						AND tb.test_date BETWEEN ta.test_date AND $1
						AND tb.animal_id = $2
						AND tb.user_id = $3
						AND tb.deleted_at IS NULL
				)
				AND NOT EXISTS (
					SELECT 1
					FROM animals a
					WHERE a.birth_date BETWEEN ta.test_date AND $1
						AND a.mother_id = $2
						AND a.user_id = $3
						AND a.deleted_at IS NULL
				)
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, entry.TransferDate, entry.ReceiverId, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.IncorrectEntityAPIError(
			"Esta vaca consta como prenha e, portanto, não pode ser receptora. Atualize os registros de toque ou selecione outra " +
				"vaca para receptação.",
		)
	}

	return nil

}
