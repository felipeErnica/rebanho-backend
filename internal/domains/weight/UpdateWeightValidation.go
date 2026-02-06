package weight

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(db *sqlx.DB, entry *WeightEntrySave) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM weight_entries
			WHERE id <> :id
				AND animal_id = :animal_id
				AND entry_date = :entry_date
				AND user_id = :user_id
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	return nil
}
