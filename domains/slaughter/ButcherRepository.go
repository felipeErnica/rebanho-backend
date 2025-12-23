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

func (r *ButcherRepository) FindById(id string, userId string) (*ButcherEntry, error) {

	query := `
		select 
			id,
			name,
			cnpj,
			discount * 100 as discount,
			address
		from butchers 
		where id = $1 
			and user_id = $2 
			and deleted_at is null
	`
	return repositoriesUtil.GetOne[ButcherEntry](r.DB, query, id, userId)
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

func (r *ButcherRepository) FindEntriesPage(
	sort string,
	order string,
	cursor string,
	filter SlaughterEntryFilter,
	butcherId string,
	userId string,
) (*entity.Page[SlaughterEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"entry_date":       {Field: "s.entry_date", Order: "desc"},
		"animal_order":     {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"animal_name":      {Field: "coalesce(a.name, '')", Order: "asc"},
		"birth_date":       {Field: "coalesce(a.birth_date, '-infinity')", Order: "desc"},
		"weight":           {Field: "s.weight", Order: "asc"},
		"dead_weight":      {Field: "s.dead_weight", Order: "asc"},
		"performance_rate": {Field: "coalesce(s.dead_weight / nullif(s.weight*(1 - s.discount_rate), 0) * 100, 0)", Order: "asc"},
		"id":               {Field: "s.id", Order: "asc"},
		"created_at":       {Field: "s.created_at", Order: "asc"},
	}

	query := `
		select 
			s.id,
			s.animal_id,
			s.butcher_id,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			a.name as animal_name, 
			concat_ws(
				' - ', 
				a.ring_number, 
				coalesce(a.name, a.sex),
				to_char(a.birth_date, 'DD/MM/YYYY')
			) as animal_info,
			a.birth_date,
			concat_ws(' - ', f.ring_number, f.name) father_name,
			concat_ws(' - ', m.ring_number, m.name) mother_name,
			h.name butcher,
			s.entry_date,
			s.discount_rate * 100 as discount_rate,
			s.weight,
			s.weight * (1 - s.discount_rate) discount_weight,
			s.dead_weight,
			coalesce(s.dead_weight / nullif(s.weight*(1 - s.discount_rate), 0) * 100, 0) performance_rate,
			s.created_at
		from slaughter_entries s
			join butchers h on h.id = s.butcher_id
			left join animals a on a.id = s.animal_id
			left join animals f on f.id = a.father_id
			left join animals m on m.id = a.mother_id
	`

	whereExpression := `
		where s.user_id = $1 
			and s.butcher_id = $2
			and s.deleted_at is null
		`

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "s", 3)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression += " and " + cursorExpression
	}

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}
	orderExpression := " order by " + sortExpression
	query += whereExpression + orderExpression

	args := []any{userId, butcherId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[SlaughterEntry](r.DB, query, sort, 200, args...)
}

func (r *ButcherRepository) FindEntriesPageFoot(
	filter SlaughterEntryFilter,
	butcherId string,
	userId string,
) (*SlaughterFoot, error) {

	query := `
		select 
			count(s.*) as animals_number,
			avg(s.weight) as avg_weight,
			avg(s.dead_weight) as avg_dead_weight,
			avg((s.dead_weight / nullif(weight * (1 - s.discount_rate), 0)) * 100) as avg_rate
		from slaughter_entries s
	`

	whereExpression := " where s.user_id = $1 and s.butcher_id = $2 and s.deleted_at is null"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "s", 3)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	args := []any{userId, butcherId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query += whereExpression

	return repositoriesUtil.GetOne[SlaughterFoot](r.DB, query, args...)
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
