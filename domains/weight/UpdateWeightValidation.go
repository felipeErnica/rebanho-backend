package weight

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(db *sqlx.DB, entry *WeightEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from weight_entries
			where id <> :id
				and animal_id = :animal_id
				and entry_date = :entry_date
				and user_id = :user_id
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}
