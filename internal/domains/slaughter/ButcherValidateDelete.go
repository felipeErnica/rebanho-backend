package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

func validateButcherDelete(db *sqlx.DB, id string, userId string) *log.APIError {

	query := `
		SELECT EXISTS (
			SELECT 1
			FROM slaughter_entries
			WHERE butcher_id = $1
				AND user_id = $2
				AND deleted_at IS NULL
		)
	`

	var exists bool
	err := util.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return log.InternalServerAPIError(err)
	}

	if exists {
		return log.DeleteAPIError("Não é possível excluir este frigorífico pois possui abates registrados nele!" +
			"\n\nOBS.: Apague ou altere os registros de abate antes de excluir o frigorífico.")
	}

	return nil
}
