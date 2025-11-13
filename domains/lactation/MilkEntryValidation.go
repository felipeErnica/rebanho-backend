package lactation

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func ValidateMilkEntryUpdate(db *sqlx.DB, entry MilkEntry, userId string) *apiError.APIError {

	existsQuery := `
		select exists (
			select 1
			from milk_entries m
			where m.animal_id = $1
				and m.entry_date = $2
				and m.id <> $3
				and m.deleted_at is null
				and m.user_id = $4
		)
	`

	exists := false
	err := repositoriesUtil.GetPrimitive(db, existsQuery, &exists, entry.AnimalId, entry.EntryDate, entry.Id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já há uma marcação desta vaca na mesma data! Exclua uma das duas marcações!")
	}

	return nil
}

func validateMilkEntry(db *sqlx.DB, entry MilkEntry) *apiError.APIError {
	err := isOnLac(db, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return err
	}

	err = entryExists(db, entry.AnimalId, entry.EntryDate, entry.UserId)
	if err != nil {
		return err
	}

	return nil
}

func isOnLac(db *sqlx.DB, animalId string, entryDate time.Time, userId string) *apiError.APIError {
	lacQuery := `
		select exists (
			select 1
			from lactations l
			where l.animal_id = $1
				and l.start_date <= $2 
				and l.end_date is null 
				and l.deleted_at is null
				and l.user_id = $3
		)
	`

	hasLac := false
	err := db.Get(&hasLac, lacQuery, animalId, entryDate, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if !hasLac {
		return apiError.IncorrectEntityAPIError(
			"Esta vaca não está lactando! Inicie uma lactação antes de adicionar marcações de leite!", 
		)
	}

	return nil
}

func entryExists(db *sqlx.DB, animalId string, entryDate time.Time, userId string) *apiError.APIError {

	existsQuery := `
		select exists (
			select 1
			from milk_entries m
			where m.animal_id = $1
				and m.entry_date = $2
				and m.deleted_at is null
				and m.user_id = $3
		)
	`
	exists := false
	err := db.Get(&exists, existsQuery, animalId, entryDate, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning(`Esta marcação já existe! Deseja substituí-la por esta?`)
	}

	return nil
}

func isDiferentPasture(db *sqlx.DB, entry AddMilkEntryStruct) *apiError.APIError {

	query := `
		select coalesce(pasture_id <> $1, false) as other_pasture
		from pasture_entries
		where animal_id = $2 and deleted_at is null
		order by entry_date desc
		limit 1
	`

	var otherPasture bool
	err := repositoriesUtil.GetPrimitive(db, query, &otherPasture, entry.PastureId, entry.AnimalId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if otherPasture {
		return apiError.NewAPIWarning(
			"Lote diferente!",
			"A vaca não está no Lote informado! Deseja transferi-la?",
			apiError.TRANSFER_WARNING,
		)
	}

	return nil

}
