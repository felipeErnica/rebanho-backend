package birth

import (
	"database/sql"

	"github.com/felipeErnica/rebanho-backend/apiError"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

func getFatherId(db *sqlx.DB, entry *BirthEntrySave) (string, *apiError.APIError) {

	fatherId, err := getInseminationFather(db, entry)
	if err != nil {
		return "", err
	}

	if fatherId != "" {
		return fatherId, nil
	}
	

	fatherId, err = getTransferFather(db, entry)
	if err != nil {
		return "", err
	}

	if fatherId != "" {
		return fatherId, nil
	}
	
	fatherId, err = getBreedingsFather(db, entry)
	if err != nil {
		return "", err
	}

	if fatherId != "" {
		return fatherId, nil
	}

	fatherId, err = getPastureFather(db, entry)
	if err != nil {
		return "", err
	}

	return fatherId, nil
	
}

func getInseminationFather(db *sqlx.DB, entry *BirthEntrySave) (string, *apiError.APIError) {

	existsQuery := `
		select exists (
			select 1
			from insemination_entries i
			where i.deleted_at is null
				and i.user_id = $1
				and i.animal_id = $2
				and i.insemination_date < $3
				and age($3, i.insemination_date) between interval '240 days' and '340 days'
				and not exists (
					select 1
					from pregnancy_tests t
					where t.deleted_at is null
						and t.pregnancy_status = 'FAILED'
						and t.test_date between i.insemination_date and $3
						and t.animal_id = $2
				)
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, existsQuery, &exists, entry.UserId, entry.MotherId, entry.BirthDate)
	if err != nil {
		return "", apiError.InternalServerAPIError(err)
	}

	if !exists {
		return "", nil
	}

	fatherQuery := `
		select bull_id
		from insemination_entries
		where deleted_at is null
			and user_id = $1 
			and animal_id = $2
			and insemination_date < $3
		order by insemination_date desc
		limit 1
	`

	var fatherId sql.NullString
	err = repositoriesUtil.GetPrimitive(db, fatherQuery, &fatherId, entry.UserId, entry.MotherId, entry.BirthDate)
	if err != nil {
		return "", apiError.InternalServerAPIError(err)
	}

	if fatherId.Valid {
		return fatherId.String, nil
	}

	return "", nil
}

func getBreedingsFather(db *sqlx.DB, entry *BirthEntrySave) (string, *apiError.APIError) {

	existsQuery := `
		select exists (
			select 1
			from breeding_entries i
			where i.deleted_at is null
				and i.user_id = $1
				and i.animal_id = $2
				and i.breeding_date < $3
				and age($3, i.breeding_date) between interval '240 days' and '340 days'
				and not exists (
					select 1
					from pregnancy_tests t
					where t.deleted_at is null
						and t.pregnancy_status = 'FAILED'
						and t.test_date between i.breeding_date and $3
						and t.animal_id = $2
				)
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, existsQuery, &exists, entry.UserId, entry.MotherId, entry.BirthDate)
	if err != nil {
		return "", apiError.InternalServerAPIError(err)
	}

	if !exists {
		return "", nil
	}

	fatherQuery := `
		select bull_id
		from breeding_entries
		where deleted_at is null
			and user_id = $1 
			and animal_id = $2
			and breeding_date < $3
		order by breeding_date desc
		limit 1
	`

	var fatherId sql.NullString
	err = repositoriesUtil.GetPrimitive(db, fatherQuery, &fatherId, entry.UserId, entry.MotherId, entry.BirthDate)
	if err != nil {
		return "", apiError.InternalServerAPIError(err)
	}

	if fatherId.Valid {
		return fatherId.String, nil
	}

	return "", nil
}

func getTransferFather(db *sqlx.DB, entry *BirthEntrySave) (string, *apiError.APIError) {

	existsQuery := `
		select exists (
			select 1
			from embryo_transfer i
			where i.deleted_at is null
				and i.user_id = $1
				and i.receiver_id = $2
				and i.transfer_date < $3
				and age($3, i.transfer_date) between interval '240 days' and '340 days'
				and not exists (
					select 1
					from pregnancy_tests t
					where t.deleted_at is null
						and t.pregnancy_status = 'FAILED'
						and t.test_date between i.transfer_date and $3
						and t.animal_id = $2
				)
		)
	`

	var exists bool
	err := repositoriesUtil.GetPrimitive(db, existsQuery, &exists, entry.UserId, entry.MotherId, entry.BirthDate)
	if err != nil {
		return "", apiError.InternalServerAPIError(err)
	}

	if !exists {
		return "", nil
	}

	fatherQuery := `
		select bull_id
		from embryo_transfer
		where deleted_at is null
			and user_id = $1 
			and receiver_id = $2
			and transfer_date < $3
		order by transfer_date desc
		limit 1
	`

	var fatherId sql.NullString
	err = repositoriesUtil.GetPrimitive(db, fatherQuery, &fatherId, entry.UserId, entry.MotherId, entry.BirthDate)
	if err != nil {
		return "", apiError.InternalServerAPIError(err)
	}

	if fatherId.Valid {
		return fatherId.String, nil
	}

	return "", nil
}

func getPastureFather(db *sqlx.DB, entry *BirthEntrySave) (string, *apiError.APIError) {

	entriesQuery := `
		with entries_query as (
			select pasture_id
			from pasture_entries
			where deleted_at is null
				and user_id = $1
				and animal_id = $2
				and ($3::timestamptz - interval '308 days') between entry_date and coalesce(exit_date, now())
		)
		select pe.animal_id
		from pasture_entries pe
			cross join entries_query e
			join animals a on a.id = pe.animal_id
				and a.animal_type = 'REPRODUCTION_ANIMAL'
				and a.sex = 'M'
		where pe.deleted_at is null
			and pe.user_id = $1
			and ($3::timestamptz - interval '308 days') between pe.entry_date and coalesce(pe.exit_date, now())
			and pe.pasture_id = e.pasture_id
		order by (coalesce(pe.exit_date, now()) - pe.entry_date) desc
		limit 1
	`

	var fatherId sql.NullString
	err := repositoriesUtil.GetPrimitive(db, entriesQuery, &fatherId, entry.UserId, entry.MotherId, entry.BirthDate)
	if err != nil {
		return "", apiError.InternalServerAPIError(err)
	}

	if fatherId.Valid {
		return fatherId.String, nil
	}

	return "", nil
}
