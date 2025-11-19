package insemination

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func inseminationExists(db *sqlx.DB, entry *InseminationEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from insemination_entries
			where animal_id = $1
				and insemination_date = $2
				and user_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.InseminationDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIError("Já existe uma entrada desta vaca na mesma data! Deseja substituí-la por esta?")
	}
	
	return nil;
}

func validateDelete(db *sqlx.DB, entry *InseminationEntrySave) *apiError.APIError {

	err := hasChildren(db, entry)
	if err != nil {
		return err
	}

	err = isPregnant(db, entry)
	if err != nil {
		return err
	}

	return nil
}

func isPregnant(db *sqlx.DB, entry *InseminationEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from pregnancy_tests t
			where t.pregnancy_status = 'SUCCESS'
				and t.animal_id = $1
				and t.test_date < $2
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
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.AnimalId, entry.InseminationDate, entry.UserId)
	if err != nil { 
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning("A vaca possui uma prenhez ligada a esta inseminação. Deseja excluir mesmo assim?")
	}

	return  nil
}

func hasChildren(db *sqlx.DB, entry *InseminationEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals a
			where a.deleted_at is null 
				and a.user_id = $1
				and a.mother_id = $2
				and a.birth_date > $3
				and a.father_id = $4
				and age(a.birth_date, $3) between interval '240 days' and interval '340 days'
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.UserId, entry.AnimalId, entry.InseminationDate, entry.BullId)
	if err != nil { 
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.NewAPIWarning(
			"Inseminação gerou crias!",
			"A vaca possui uma cria, cujo o pai está ligado a esta inseminação. Deseja excluir mesmo assim?" +
			"\n\nOBS.: Ao aceitar, o pai do bezerro será modificado automaticamente.",
			"ChildrenWarning",
		)
	}

	return  nil
}
