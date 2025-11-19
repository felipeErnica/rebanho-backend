package birth

import (
	"database/sql"
	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func validateAddBirth(db *sqlx.DB, entry *BirthEntrySave) *apiError.APIError {

	err := hasValidInterval(db, entry)
	if err != nil {
		return err
	}

	err = doesExist(db, entry)
	if err != nil {
		return err
	}

	if entry.RingNumber != nil {
		err = ringExists(db, entry)
		if err != nil {
			return err
		}
	}

	return nil

}

func validateUpdateBirth(db *sqlx.DB, entry *BirthEntrySave) *apiError.APIError {

	err := birthExist(db, entry)
	if err != nil {
		return err
	}

	err = hasValidInterval(db, entry)
	if err != nil {
		return err
	}

	return nil

}

func birthExist(db *sqlx.DB, entry *BirthEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals
			where mother_id = $1
				and birth_date = $2
				and id <> $3
				and user_id = $4
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.MotherId, entry.BirthDate, entry.Id, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil

}

func doesExist(db *sqlx.DB, entry *BirthEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals
			where mother_id = $1
				and birth_date = $2
				and user_id = $3
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.MotherId, entry.BirthDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.ConflictAPIWarning("Este nascimento já existe. Deseja substitui-lo?")
	}

	return nil

}

func ringExists(db *sqlx.DB, entry *BirthEntrySave) *apiError.APIError {

	query := `
		select exists (
			select 1
			from animals
			where ring_number = $1 
				and user_id = $2
				and deleted_at is null
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, query, &exists, entry.RingNumber, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if exists {
		return apiError.NewAPIWarning(
			"Brinco já existe!",
			"Este brinco já existe. Deseja adicionar mesmo assim?",
			"RingWarning",
		)
	}

	return nil

}

func hasValidInterval(db *sqlx.DB, entry *BirthEntrySave) *apiError.APIError {

	const MIN_INTERVAL = 240

	beforeQuery := `
		select birth_date
		from animals
		where mother_id = $1 
			and birth_date < $2
			and user_id = $3
			and deleted_at is null
		order by birth_date desc
		limit 1
	`

	var beforeBirthDate sql.NullTime
	err := repositoriesUtil.GetPrimitive(db, beforeQuery, &beforeBirthDate, entry.MotherId, entry.BirthDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if beforeBirthDate.Valid {
		interval := entry.BirthDate.Sub(beforeBirthDate.Time)
		if (interval.Hours() / 24) <= MIN_INTERVAL {
			return apiError.IncorrectEntityAPIError(
				"O intervalo em relação ao nascimento anterior é muito pequeno. O intervalo deve ser maior que 240 dias!",
			)
		}
	}

	afterQuery := `
		select birth_date
		from animals
		where mother_id = $1 
			and birth_date > $2
			and user_id = $3
			and deleted_at is null
		order by birth_date
		limit 1
	`

	var afterBirthDate sql.NullTime
	err = repositoriesUtil.GetPrimitive(db, afterQuery, &afterBirthDate, entry.MotherId, entry.BirthDate, entry.UserId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if afterBirthDate.Valid {
		interval := afterBirthDate.Time.Sub(entry.BirthDate)
		if (interval.Hours() / 24) <= MIN_INTERVAL {
			return apiError.IncorrectEntityAPIError(
				"O intervalo em relação ao nascimento posterior é muito pequeno. O intervalo deve ser maior que 240 dias!",
			)
		}
	}

	return nil

}
