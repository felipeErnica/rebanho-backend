package naturalBreeding

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func breedingExists(db *sqlx.DB, entry *BreedingEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from breeding_entries
			where animal_id = $1
				and breeding_date = $2
				and user_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.BreedingDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já existe uma entrada desta vaca na mesma data! Deseja substituí-la por esta?")
	}

	return nil
}
