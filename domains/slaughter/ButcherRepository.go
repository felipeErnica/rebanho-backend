package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type ButcherRepository struct {
	DB *sqlx.DB
}

func newButcherRepository(db *sqlx.DB) *ButcherRepository {
	return &ButcherRepository{db}
}

func (r *ButcherRepository) FindAll(userId string) (*[]ButcherEntry, error) {

	query := `
		with entries_stats as (
			select
				butcher_id,
				count(*) as animals_number,
				avg(dead_weight) as avg_weight,
				avg(coalesce(dead_weight / nullif(weight * (1 - discount_rate), 0) * 100, 0)) as avg_rate
			from slaughter_entries
			where deleted_at is null and user_id = $1
			group by butcher_id
		)
		select 
			b.id,
			b.name,
			b.cnpj,
			b.discount * 100 as discount,
			b.address,
			s.animals_number,
			s.avg_weight,
			s.avg_rate
		from butchers b
			join entries_stats s on s.butcher_id = b.id
		where b.user_id = $1 and b.deleted_at is null
		order by name
	`
	return repositoriesUtil.GetList[ButcherEntry](r.DB, query, userId)
}

func (r *ButcherRepository) Search(userId string) (*[]entity.SearchEntity, error) {

	query := `
		select
			s.id,
			s.name as label
		from butchers s
		where s.user_id = $1 and s.deleted_at is null
	`
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *ButcherRepository) Add(entry *ButcherSave) *apiError.APIError {

	validateErr := validateButcherAdd(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		insert into butchers (name, cnpj, address, discount, user_id)
		values (:name, :cnpj, :address, cast(:discount as float) / 100, :user_id)
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *ButcherRepository) Replace(entry *ButcherSave) *apiError.APIError {

	query := `
		update butchers 
		set name = :name, 
			cnpj = :cnpj, 
			discount = cast(:discount as float) / 100, 
			address = :address
		where name = :name
			or cnpj = :cnpj 
			and cnpj is not null
			and user_id = :user_id
			and deleted_at is null
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *ButcherRepository) Update(entry *ButcherSave) (*ButcherEntry, *apiError.APIError) {

	validateErr := validateButcherUpdate(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		update butchers 
		set name = :name, 
			cnpj = :cnpj, 
			discount = cast(:discount as float) / 100, 
			address = :address
		where id = :id and user_id = :user_id
	`

	err = repositoriesUtil.NamedExecTx(tx, query, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		with entries_stats as (
			select
				butcher_id,
				count(*) as animals_number,
				avg(dead_weight) as avg_weight,
				avg(coalesce(dead_weight / nullif(weight * (1 - discount_rate), 0) * 100, 0)) as avg_rate
			from slaughter_entries
			group by butcher_id
		)
		select 
			b.id,
			b.name,
			b.cnpj,
			b.discount * 100 as discount,
			b.address,
			coalesce(s.animals_number, 0) as animals_number,
			coalesce(s.avg_weight, 0) as avg_weight,
			coalesce(s.avg_rate, 0) as avg_rate
		from butchers b 
			left join entries_stats s on s.butcher_id = b.id
		where b.id = :id and b.user_id = :user_id
	`

	response, err := repositoriesUtil.NamedGetTx(tx, selectQuery, ButcherEntry{}, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *ButcherRepository) Delete(id string, userId string) *apiError.APIError {
	
	validateErr := validateButcherDelete(r.DB, id, userId)
	if validateErr != nil {
		return validateErr
	}

	query := `
		update butchers
		set deleted_at = now()
		where id = $1 and user_id = $2
	`

	err := repositoriesUtil.Exec(r.DB, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

