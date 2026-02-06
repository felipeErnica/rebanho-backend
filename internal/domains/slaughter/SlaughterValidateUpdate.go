package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(tx *sqlx.Tx, entry *SlaughterEntrySave) *log.APIError {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM slaughter_entries
			WHERE entry_date = $1
				AND animal_id = $2
				AND id <> $3
				AND user_id = $4
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitiveTx(tx, query, &exists, entry.EntryDate, entry.AnimalId, entry.Id, entry.UserId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.ConflictAPIError("Já há um abate registrado deste animal nesta data!")
	}

	return nil
}
