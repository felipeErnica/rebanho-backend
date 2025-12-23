package lactation

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateGroupUpdate(db *sqlx.DB, entry LactationGroupSave) *apiError.APIError {
	query := `
		select exists (
			select 1
			from milk_entries
			where entry_date = :entry_date
				and user_id = :user_id
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.NamedPrimitive(db, query, &exists, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já exite uma marcação com esta data!")
	}

	return nil
}

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

func validateMilkEntry(db *sqlx.DB, entry MilkEntrySave) *apiError.APIError {
	err := isOnLac(db, entry)
	if err != nil {
		return err
	}

	err = entryExists(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func isOnLac(db *sqlx.DB, entry MilkEntrySave) *apiError.APIError {
	lacQuery := `
		select exists (
			select 1
			from lactations l
			where l.animal_id = :animal_id
			and l.start_date <= :entry_date
				and l.end_date is null 
				and l.deleted_at is null
				and l.user_id = :user_id
		)
	`

	var hasLac bool
	err := repositoriesUtil.NamedPrimitive(db, lacQuery, &hasLac, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if !hasLac {
		return apiError.IncorrectEntityAPIError("Esta vaca não está lactando! Inicie uma lactação antes de " + 
												"adicionar marcações de leite!")
	}

	return nil
}

func entryExists(db *sqlx.DB, entry MilkEntrySave) *apiError.APIError {

	existsQuery := `
		select exists (
			select 1
			from milk_entries m
			where m.animal_id = :animal_id
				and m.entry_date = :entry_date
				and m.user_id = :user_id
				and m.deleted_at is null
		)
	`
	var exists bool
	err := repositoriesUtil.NamedPrimitive(db, existsQuery, &exists, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning(`Esta marcação já existe! Deseja substituí-la por esta?`)
	}

	return nil
}

func isDiferentPasture(db *sqlx.DB, entry MilkEntrySave) *apiError.APIError {

	query := `
		select coalesce(pasture_id <> :pasture_id, false) as other_pasture
		from pasture_entries
		where animal_id = :animal_id and deleted_at is null
		order by entry_date desc
		limit 1
	`

	var otherPasture bool
	err := repositoriesUtil.NamedPrimitive(db, query, &otherPasture, entry)
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
