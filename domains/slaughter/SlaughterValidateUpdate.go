package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(tx *sqlx.Tx, entry *SlaughterEntrySave) *apiError.APIError {
	query := `
		select exists (
			select 1
			from slaughter_entries
			where entry_date = $1
				and animal_id = $2
				and id <> $3
				and user_id = $4
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitiveTx(tx, query, &exists, entry.EntryDate, entry.AnimalId, entry.Id, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já há um abate registrado deste animal nesta data!")
	}

	return nil
}
