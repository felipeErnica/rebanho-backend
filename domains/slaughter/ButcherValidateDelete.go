package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateButcherDelete(db *sqlx.DB, id string, userId string) *apiError.APIError {

	query := `
		select exists (
			select 1
			from slaughter_entries
			where butcher_id = $1
				and user_id = $2
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.DeleteAPIError("Não é possível excluir este frigorífico pois possui abates registrados nele!" +
										"\n\nOBS.: Apague ou altere os registros de abate antes de excluir o frigorífico.")
	}

	return nil
}
