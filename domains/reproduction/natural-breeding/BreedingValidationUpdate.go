package naturalBreeding

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateUpdate(db *sqlx.DB, oldEntry *BreedingEntrySave, newEntry *BreedingEntrySave) *apiError.APIError {

	err := repeatedInfo(db, newEntry)
	if err != nil {
		return err
	}

	err = hasChildrenError(db, oldEntry)
	if err != nil {
		return err
	}

	err = isPregnantUpdate(db, oldEntry)
	if err != nil {
		return err
	}

	return nil
}

func validateUpdateBatch(db *sqlx.DB, group *BreedingGroup) *apiError.APIError {

	query := `
		select exists (
			select 1
			from breeding_entries
			where breeding_date = $1
				and user_id = $2
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, group.BreedingDate, group.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já existem registros de inseminação nesta mesma data! " + 
			"Para evitar conflitos, altere os registros existentes ou escolha outra data.",
		)
	}

	return nil
}

func repeatedInfo(db *sqlx.DB, entry *BreedingEntrySave) *apiError.APIError {
	query := `
		select exists (
			select 1
			from breeding_entries
			where animal_id = $1
				and breeding_date = $2
				and user_id = $3
				and id <> $4
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.BreedingDate, entry.UserId, entry.Id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já existe uma inseminação desta vaca na mesma data!")
	}

	return nil
}

func isPregnantUpdate(db *sqlx.DB, entry *BreedingEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from pregnancy_tests t
			where t.pregnancy_status = 'SUCCESS'
				and t.animal_id = $1
				and t.test_date > $2
				and age(t.test_date, $2) <= interval '340 days'
				and t.user_id = $3
				and t.deleted_at is null
		) and not exists (
			select 1
			from pregnancy_tests t
			where t.pregnancy_status = 'FAILED'
				and t.animal_id = $1
				and t.test_date > $2
				and age(t.test_date, $2) <= interval '340 days'
				and t.user_id = $3
				and t.deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.BreedingDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning("A vaca possui uma prenhez ligada a esta cobertura! Deseja alterar mesmo assim?")
	}

	return nil
}

func hasChildrenError(db *sqlx.DB, entry *BreedingEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals a
			where a.deleted_at is null 
				and a.user_id = $1
				and a.mother_id = $2
				and a.birth_date > $3
				and age(a.birth_date, $3) between interval '240 days' and interval '340 days'
				and a.father_id = $4
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.UserId, entry.AnimalId, entry.BreedingDate, entry.BullId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.IncorrectEntityAPIError(
			"A vaca possui uma cria, cujo o pai está ligado a esta cobertura. Altere ou exclua a parição antes de continuar!",
		)
	}

	return nil
}
