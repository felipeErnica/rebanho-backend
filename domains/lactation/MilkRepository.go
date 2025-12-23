package lactation

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type MilkRepository struct {
	DB *sqlx.DB
}

func NewMilkRepository(db *sqlx.DB) *MilkRepository {
	return &MilkRepository{db}
}

func (r *MilkRepository) FindGroupsPage(
	filter LactationGroupFilter,
	order string,
	cursor string,
	userId string,
) (*entity.Page[LactationGroup], error) {

	sortMap := map[string]repositoriesUtil.SortField{"entry_date": {Field: "cte.entry_date", Order: "asc"}}

	query := `
		with cte as (
			select 
				entry_date,
				count(*) animals_number,
				sum(quantity) total_milk,
				avg(quantity) avg_milk
			from milk_entries
			where user_id = $1 and deleted_at is null
			group by 1
		)
		select 
			cte.*,
			coalesce(animals_number - lag(animals_number) over (order by entry_date), 0) number_difference,
			coalesce(((total_milk / lag(total_milk) over (order by entry_date)) - 1)*100, 0) total_rate,
			coalesce(((avg_milk / lag(avg_milk) over (order by entry_date)) - 1)*100, 0) avg_rate
		from cte
    `
	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, "entry_date", order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	var whereExpression string
	if filterExpression != "" {
		whereExpression = "where " + filterExpression
	}

	if cursorExpression != "" {
		if whereExpression != "" {
			whereExpression += " and " + cursorExpression
		} else {
			whereExpression = "where " + cursorExpression
		}
	}

	query += whereExpression + " order by cte.entry_date " + order
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[LactationGroup](r.DB, query, "entry_date", 100, args...)
}

func (r *MilkRepository) FindEntriesPage(
	filter MilkEntryFilter,
	sort string,
	order string,
	cursor string,
	userId string,
) (*entity.Page[MilkEntry], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"animal_name":  {Field: "a.name", Order: "asc"},
		"animal_order": {Field: "coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"entry_date":   {Field: "m.entry_date", Order: "desc"},
		"quantity":     {Field: "m.quantity", Order: "asc"},
		"id":           {Field: "m.id", Order: "asc"},
		"created_at":   {Field: "m.created_at", Order: "asc"},
	}

	query := `
		select
			m.id,
			m.animal_id,
			a.name as animal_name,
			concat_ws(' - ', a.ring_number, a.name) as animal_info,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
			coalesce(p.name, 'Sem Pasto') as pasture_name,
			m.entry_date,
			m.quantity,
			m.created_at
		from milk_entries m
			join animals a on a.id = m.animal_id
			left join pasture_entries pe on pe.animal_id = m.animal_id
				and m.entry_date >= pe.entry_date
				and m.entry_date < coalesce(pe.exit_date, now())
				and pe.deleted_at is null
			left join pastures p on p.id = pe.pasture_id
    `

	whereExpression := "where m.user_id = $1 and m.deleted_at is null"

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "m", 2)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
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

	query = query + whereExpression

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query += " order by " + sortExpression
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	return repositoriesUtil.GetPage[MilkEntry](r.DB, query, sort, 100, args...)
}

func (r *MilkRepository) GetEntriesPageFoot(filter MilkEntryFilter, userId string) (*MilkEntryFoot, error) {
	query := `
		select
			count(*) animals_number,
			sum(quantity) total_milk,
			avg(quantity) avg_milk
		from milk_entries m
    `
	whereExpression := "where m.user_id = $1 and m.deleted_at is null"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "m", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression += " and " + filterExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	query = query + whereExpression
	return repositoriesUtil.GetOne[MilkEntryFoot](r.DB, query, args...)
}

func (r *MilkRepository) GetGroupEntries(userId string, entryDate time.Time) (*[]MilkEntry, error) {

	query := `
		select
			m.id,
			m.animal_id,
			concat_ws(' - ', a.ring_number, a.name) animal_info,
			coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0) animal_order,
			p.name pasture_name,
			m.entry_date,
			m.quantity
		from milk_entries m
			join animals a on a.id = m.animal_id
			join pasture_entries pe on pe.animal_id = m.animal_id
				and pe.entry_date <= m.entry_date
				and coalesce(pe.exit_date, now()) > m.entry_date
				and pe.deleted_at is null
			join pastures p on p.id = pe.pasture_id
		where m.user_id = $1 and m.deleted_at is null and m.entry_date = $2
		order by coalesce(regexp_replace(a.ring_number, '[^0-9]', '', 'g')::int, 0)
    `
	return repositoriesUtil.GetList[MilkEntry](r.DB, query, userId, entryDate)
}

func (r *MilkRepository) GetGroupEntriesFoot(userId string, entryDate time.Time) (*MilkEntryFoot, error) {
	query := `
		select
			count(*) as animals_number,
			sum(quantity) as total_milk,
			avg(quantity) as avg_milk
		from milk_entries m
		where m.user_id = $1 and deleted_at is null and m.entry_date = $2
    `
	return repositoriesUtil.GetOne[MilkEntryFoot](r.DB, query, userId, entryDate)
}

func (r *MilkRepository) UpdateGroup(entryDate time.Time, groupEntry *LactationGroupSave) (*LactationGroup, *apiError.APIError) {
	
	validateErr := validateGroupUpdate(r.DB, *groupEntry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		update milk_entries
		set entry_date = $1
		where entry_date = $2 
			and user_id = $3
			and deleted_at is null
	`
	err := repositoriesUtil.Exec(r.DB, query, groupEntry.EntryDate, entryDate, groupEntry.UserId)
	if err != nil {
		return  nil, apiError.InternalServerAPIError(err)
	}

	returnQuery := `
		with milk_stats as (
			select 
				entry_date,
				count(*) animals_number,
				sum(quantity) total_milk,
				avg(quantity) avg_milk
			from milk_entries
			where user_id = :user_id and deleted_at is null
			group by 1
		),
		cte as (
			select 
				s.*,
				coalesce(animals_number - lag(animals_number) over (order by entry_date), 0) as number_difference,
				coalesce(((total_milk / lag(total_milk) over (order by entry_date)) - 1) * 100, 0) as total_rate,
				coalesce(((avg_milk / lag(avg_milk) over (order by entry_date)) - 1) * 100, 0) as avg_rate
			from milk_stats s
		)
		select * from cte where entry_date = :entry_date
	`
	response, err := repositoriesUtil.NamedGet(r.DB, returnQuery, LactationGroup{}, groupEntry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *MilkRepository) DeleteGroup(entryDate time.Time, userId string) *apiError.APIError {
	deleteQuery := `
		update milk_entries
		set deleted_at = now()
		where entry_date = $1 and user_id = $2
	`

	err := repositoriesUtil.Exec(r.DB, deleteQuery, entryDate, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *MilkRepository) Add(entry *MilkEntrySave) *apiError.APIError {

	apiErr := validateMilkEntry(r.DB, *entry)
	if apiErr != nil {
		return apiErr
	}

	insertQuery := `
		insert into milk_entries (animal_id, entry_date, quantity, user_id) 
		values (:animal_id, :entry_date, :quantity, :user_id)
	`

	err := repositoriesUtil.NamedExec(r.DB, insertQuery, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pastureErr := isDiferentPasture(r.DB, *entry)
	return pastureErr
}

func (r *MilkRepository) Replace(entry *MilkEntrySave, userId string) *apiError.APIError {

	query := `
		update milk_entries 
		set quantity = :quantity,
			created_at = now()
		where animal_id = :animal_id 
			and entry_date = :entry_date
			and user_id = :user_id
	`

	_, err := r.DB.NamedExec(query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pastureErr := isDiferentPasture(r.DB, *entry)
	return pastureErr
}

func (r *MilkRepository) Update(entry *MilkEntry) (*MilkEntry, *apiError.APIError) {

	apiErr := ValidateMilkEntryUpdate(r.DB, *entry, entry.UserId)
	if apiErr != nil {
		return nil, apiErr
	}

	query := `
		update milk_entries 
		set entry_date = :entry_date,
			quantity = :quantity
		where id = :id and user_id = :user_id
	`

	id, err := repositoriesUtil.NamedExecReturningId(r.DB, query, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	returnQuery := `
		select 
			m.id,
			m.animal_id,
			concat_ws(' - ', a.ring_number, a.name) as animal_info,
			coalesce(p.name, 'Sem Pasto') as pasture_name,
			m.entry_date,
			m.quantity
		from milk_entries m 
			join animals a on a.id = m.animal_id
			left join pasture_entries pe on 
				pe.animal_id = m.animal_id
				and pe.entry_date <= m.entry_date
				and m.entry_date <= coalesce(pe.exit_date, now())
			left join pastures p on p.id = pe.pasture_id
		where m.id = $1
	`

	response, err := repositoriesUtil.GetOne[MilkEntry](r.DB, returnQuery, id)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil
}

func (r *MilkRepository) Delete(id string) *apiError.APIError {
	deleteQuery := `
		update milk_entries
		set deleted_at = now()
		where id = $1
	`

	err := repositoriesUtil.Exec(r.DB, deleteQuery, id)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}
