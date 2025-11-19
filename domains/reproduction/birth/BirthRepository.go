package birth

import (
	"database/sql"

	"github.com/felipeErnica/rebanho-backend/apiError"
	pastureEntries "github.com/felipeErnica/rebanho-backend/domains/farm-area/pasture-entries"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type BirthRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *BirthRepository {
	return &BirthRepository{db}
}

func (r *BirthRepository) GetBestIntervals(userId string) (*[]IntervalAnimal, error) {
	query := `
        with interval_list as (
			select 
				a.mother_id,
				extract(days from a.birth_date - lag(a.birth_date) over win) as birth_interval
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.birth_date is not null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
			window win as (partition by a.mother_id order by a.birth_date)
		),
		average_list as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) birth_numbers,
                avg(l.birth_interval) interval_average
            from interval_list l
                join animals a on 
					a.id = l.mother_id
					and a.death_date is null
            group by animal_name
        ),
        birth_stats as (
            select 
				avg(interval_average) gn_interval_average,
				stddev(interval_average) dev_interval
            from average_list
			where birth_numbers >= 3
        ),
        scores as (
            select 
                a.animal_name,
                a.interval_average,
                a.birth_numbers,
                (a.interval_average - b.gn_interval_average) / nullif(b.dev_interval, 0) as reproductive_score,
                ((a.interval_average / nullif(b.gn_interval_average, 0)) - 1) * 100 average_rate
            from average_list a, birth_stats b
			where a.birth_numbers >= 3
        )
        select *
        from scores
		where -reproductive_score > 0
        order by -reproductive_score desc
		limit 10
    `
	return repositoriesUtil.GetList[IntervalAnimal](r.DB, query, userId)
}

func (r *BirthRepository) GetWorstIntervals(userId string) (*[]IntervalAnimal, error) {
	query := `
        with interval_list as (
			select 
				a.mother_id,
				extract(days from a.birth_date - lag(a.birth_date) over win) as birth_interval
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.birth_date is not null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
			window win as (partition by a.mother_id order by a.birth_date)
		),
		average_list as (
            select
                concat_ws(' - ', a.ring_number, a.name) animal_name,
                count(l.*) birth_numbers,
                avg(l.birth_interval) interval_average
            from interval_list l
                join animals a on 
					a.id = l.mother_id
					and a.death_date is null
            group by animal_name
        ),
        birth_stats as (
            select 
				avg(interval_average) gn_interval_average,
				stddev(interval_average) dev_interval
            from average_list
			where birth_numbers >= 3
        ),
        scores as (
            select 
                a.animal_name,
                a.interval_average,
                a.birth_numbers,
                (a.interval_average - b.gn_interval_average) / nullif(b.dev_interval, 0) as reproductive_score,
                ((a.interval_average / nullif(b.gn_interval_average, 0)) - 1) * 100 average_rate
            from average_list a, birth_stats b
			where a.birth_numbers >= 3
        )
        select *
        from scores
		where reproductive_score > 0
        order by reproductive_score desc
		limit 10
    `
	return repositoriesUtil.GetList[IntervalAnimal](r.DB, query, userId)
}

func (r *BirthRepository) GetBirthIntervalHistory(userId string) (*IntervalStats, error) {
	query := `
		with year_birth_interval as (
			select 
				date_trunc('year', a.birth_date) birth_date,
				extract(days from a.birth_date - lag(a.birth_date) over win) birth_interval
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.birth_date is not null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
			window win as (partition by a.mother_id order by a.birth_date)
		),
		cte as (
			select 
				birth_date,
				avg(birth_interval) interval_average
			from year_birth_interval
			where birth_interval <> 0
			group by 1
			order by birth_date desc
			limit 10
		)
		select * from cte order by birth_date
    `
	results, err := repositoriesUtil.GetList[BirthIntervalHist](r.DB, query, userId)

	if err != nil {
		return nil, err
	}

	intervalHist := *results
	var currentInterval, previousInterval, intervalTrend float64

	switch lenght := len(intervalHist); lenght {
	case 0:
		currentInterval = 0
		previousInterval = 0
		intervalTrend = 0
	case 1:
		currentInterval = intervalHist[lenght-1].IntervalAverage
		previousInterval = 0
		intervalTrend = 0
	default:
		currentInterval = intervalHist[lenght-1].IntervalAverage
		previousInterval = intervalHist[lenght-2].IntervalAverage
		intervalTrend = ((currentInterval / previousInterval) - 1) * 100
	}

	intervalResult := IntervalStats{
		IntervalTrend:     intervalTrend,
		BirthIntervalHist: intervalHist,
		CurrentInterval:   currentInterval,
	}

	return &intervalResult, nil
}

func (r *BirthRepository) GetLastBirthsNumber(userId string) (*CurrentStats, error) {
	query := `
		with cte as (
			select
				date_trunc('month', a.birth_date) entry_date,
				count(a.*) birth_total
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.birth_date is not null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
			group by 1
			order by 1 desc
			limit 10
		)
		select * from cte order by 1
    `
	results, err := repositoriesUtil.GetList[BirthsNumberEntry](r.DB, query, userId)

	if err != nil {
		return nil, err
	}

	hist := *results
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[lenght-1].BirthTotal
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].BirthTotal
		previous = hist[lenght-2].BirthTotal
		trend = current - previous
	}

	birthStats := CurrentStats{
		Hist:    hist,
		Trend:   trend,
		Current: current,
	}

	return &birthStats, nil
}

func (r *BirthRepository) GetYearBirthsNumber(userId string) (*CurrentStats, error) {
	query := `
		with cte as (
			select
				date_trunc('year', a.birth_date) entry_date,
				count(a.*) birth_total
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.birth_date is not null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
			group by 1
			order by 1 desc
			limit 20
		)
		select * from cte order by entry_date
    `
	results, err := repositoriesUtil.GetList[BirthsNumberEntry](r.DB, query, userId)

	if err != nil {
		return nil, err
	}

	hist := *results
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[lenght-1].BirthTotal
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].BirthTotal
		previous = hist[lenght-2].BirthTotal
		trend = util.CalculatePercentageTrend(current, previous)
	}

	birthStats := CurrentStats{
		Hist:    hist,
		Trend:   trend,
		Current: current,
	}

	return &birthStats, nil
}

func (r *BirthRepository) GetYearDeathsNumber(userId string) (*CurrentStats, error) {
	query := `
		with cte as (
			select
				date_trunc('year', a.death_date) entry_date,
				count(a.*) deaths_total
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
                and age(a.death_date, a.birth_date) < interval '1 year'
			group by 1
			order by 1 desc
			limit 20
		)
		select * from cte order by entry_date
    `
	results, err := repositoriesUtil.GetList[DeathsNumberEntry](r.DB, query, userId)

	if err != nil {
		return nil, err
	}

	hist := *results
	var current, previous, trend float64

	switch lenght := len(hist); lenght {
	case 0:
		current = 0
		previous = 0
		trend = 0
	case 1:
		current = hist[lenght-1].DeathsTotal
		previous = 0
		trend = 0
	default:
		current = hist[lenght-1].DeathsTotal
		previous = hist[lenght-2].DeathsTotal
		trend = util.CalculatePercentageTrend(current, previous)
	}

	deathStats := CurrentStats{
		Hist:    hist,
		Trend:   trend,
		Current: current,
	}

	return &deathStats, nil
}

func (r *BirthRepository) GetDeathIndex(userId string) (*DeathStats, error) {
	query := `
        with death_tbl as (
            select 
                date_trunc('year', a.death_date) date,
                count(a.*) deaths
            from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
            where
                a.user_id = $1
                and a.death_date is not null
                and age(a.death_date, a.birth_date) < interval '1 year'
                and a.deleted_at is null
            group by 1
        ),
        birth_tbl as (
            select
                date_trunc('year', a.birth_date) date,
                count(a.*) births
            from animals a
            where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.birth_date is not null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
            group by 1
        ),
        cte as (
			select
				date,
				coalesce((deaths::float / nullif(births, 0)::float)*100, 0) death_index
			from birth_tbl full join death_tbl using(date)
			order by 1 desc
			limit 10
		)
		select * from cte order by date
    `
	results, err := repositoriesUtil.GetList[DeathIndexHist](r.DB, query, userId)

	if err != nil {
		return nil, err
	}

	indexHist := *results
	var currentIndex, previousIndex, indexTrend float64

	switch lenght := len(indexHist); lenght {
	case 0:
		currentIndex = 0
		previousIndex = 0
		indexTrend = 0
	case 1:
		currentIndex = indexHist[lenght-1].DeathIndex
		previousIndex = 0
		indexTrend = 0
	default:
		currentIndex = indexHist[lenght-1].DeathIndex
		previousIndex = indexHist[lenght-2].DeathIndex
		indexTrend = ((currentIndex / previousIndex) - 1) * 100
	}

	deathStats := DeathStats{
		DeathIndexHist:  indexHist,
		DeathIndexTrend: indexTrend,
		DeathIndex:      currentIndex,
	}

	return &deathStats, nil
}

func (r *BirthRepository) GetBirthHistory(userId string) (*[]BirthsByDate, error) {
	query := `
        with death_data as ( 
            select 
                date_trunc('month', death_date) date,
                count(*) death_total 
			from animals  
            where 
                user_id = $1
                and deleted_at is null 
                and death_date is not null
                and age(death_date, birth_date) < interval '1 year'
            group by 1
        ), 
        birth_data as (
            select 
                date_trunc('month', a.birth_date) date,
                count(a.*) as birth_total
            from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
            where 
                a.user_id = $1
                and a.deleted_at is null 
				and a.animal_type <> 'OUTSIDE_ANIMAL'
				and a.birth_date is not null
            group by 1
        ),
        cte as (
			select
				date,
				coalesce(birth_data.birth_total,0) birth_total,
				coalesce(death_data.death_total, 0) death_total
			from birth_data full join death_data using(date)
			order by date desc
			limit 60
		) 
		select * from cte order by date
    `
	return repositoriesUtil.GetList[BirthsByDate](r.DB, query, userId)
}

func (r *BirthRepository) TotalBySex(userId string) (*[]TotalBirthsBySex, error) {
	query := `
        with cte as (
			select 
				date_trunc('month', a.birth_date) birth_month,
				count(a.*) filter (where a.sex = 'M') males,
				count(a.*) filter (where a.sex = 'F') females
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
				and a.birth_date is not null
			group by birth_month
			order by birth_month desc
			limit 24
		)
		select * from cte order by birth_month
    `
	return repositoriesUtil.GetList[TotalBirthsBySex](r.DB, query, userId)
}

func (r *BirthRepository) GetYearBySex(userId string) (*[]TotalBirthsBySex, error) {
	query := `
		select 
			date_trunc('year', a.birth_date) birth_month,
			count(a.*) filter (where a.sex = 'M') males,
			count(a.*) filter (where a.sex = 'F') females
		from animals a
			join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
		where 
			a.user_id = $1 
			and a.deleted_at is null
			and a.birth_date is not null
			and a.animal_type <> 'OUTSIDE_ANIMAL'
		group by birth_month
		order by birth_month desc
		limit 10
    `
	return repositoriesUtil.GetList[TotalBirthsBySex](r.DB, query, userId)
}

func (r *BirthRepository) GetLastBirths(userId string) (*[]BirthEntry, error) {
	query := `
        select 
            concat_ws(' - ', m.ring_number, m.name) mother_info,
            a.birth_date as calf_birth_date,
            a.sex as calf_sex,
            concat_ws(' - ', f.ring_number, f.name) as calf_father,
			extract(days from a.birth_date - lag(a.birth_date) over win) as birth_interval
        from animals a
            join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
            left join animals f on f.id = a.father_id
		where 
			a.user_id = $1 
			and a.deleted_at is null
			and a.birth_date is not null
			and a.animal_type <> 'OUTSIDE_ANIMAL'
		window win as (partition by a.mother_id order by a.birth_date)
		order by a.birth_date desc, coalesce(regexp_replace(m.ring_number, '[^0-9]', '', 'g')::int, 0)
		limit 15
    `
	return repositoriesUtil.GetList[BirthEntry](r.DB, query, userId)
}

func (r *BirthRepository) FindPage(
	userId string,
	sort string,
	order string,
	filter BirthEntryFilter,
	cursor string,
) (*entity.Page[BirthEntry], error) {

	sort = repositoriesUtil.AddNewFields(sort, "id")
	sortMap := map[string]repositoriesUtil.SortField{
		"calf_birth_date": {Field: "cte.calf_birth_date", Order: "asc"},
		"mother_order":    {Field: "cte.mother_order", Order: "asc"},
		"mother_name":     {Field: "cte.mother_name", Order: "asc"},
		"birth_interval":  {Field: "coalesce(cte.birth_interval, 0)", Order: "asc"},
		"id":              {Field: "cte.id", Order: "asc"},
	}

	query := `
		with cte as (
			select 
				a.id,
				a.mother_id,
				m.name as mother_name,
				concat_ws(' - ', m.ring_number, m.name) as mother_info,
				coalesce(regexp_replace(m.ring_number, '[^0-9]', '', 'g')::int, 0) as mother_order,
				a.birth_date as calf_birth_date,
				a.sex as calf_sex,
				case 
					when a.name is null then ''
					else concat_ws(' - ', a.ring_number, a.name)
				end as calf_name,
				a.father_id as calf_father_id,
				concat_ws(' - ', f.ring_number, f.name) calf_father,
				extract(days from a.birth_date - lag(a.birth_date) over win) as birth_interval
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
				left join animals f on f.id = a.father_id
			where 
				a.user_id = $1 
				and a.deleted_at is null 
				and a.animal_type <> 'OUTSIDE_ANIMAL'
				and a.birth_date is not null
				and a.mother_id is not null
			window win as (partition by a.mother_id order by a.birth_date)
		)
		select * from cte
    `
	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	cursorExpression, nextParam, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	var whereExpression string

	if filterExpression != "" {
		whereExpression = " where " + filterExpression
	}

	if cursorExpression != "" {
		if whereExpression == "" {
			whereExpression = " where " + cursorExpression
		} else {
			whereExpression += " and " + cursorExpression
		}
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	orderExpression := " order by " + sortExpression
	query += whereExpression + orderExpression
	return repositoriesUtil.GetPage[BirthEntry](r.DB, query, sort, 100, args...)
}

func (r *BirthRepository) FindPageFooter(userId string, filter BirthEntryFilter) (*BirthFooter, error) {
	query := `
		with animal_cte as (
			select
				a.mother_id,
				a.father_id as calf_father_id,
				a.birth_date as calf_birth_date,
				a.sex as calf_sex,
				extract(days from a.birth_date - lag(a.birth_date) over win) as birth_interval
			from animals a
				join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			where 
				a.user_id = $1 
				and a.deleted_at is null
				and a.birth_date is not null
				and a.animal_type <> 'OUTSIDE_ANIMAL'
			window win as (partition by a.mother_id order by a.birth_date)
		)
		select
			count(a.*) as total,
			avg(a.birth_interval) as interval_average
		from animal_cte a
    `
	whereExpression := ""

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "a", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression = " where " + filterExpression
	}

	query += whereExpression
	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	return repositoriesUtil.GetOne[BirthFooter](r.DB, query, args...)
}

func (r *BirthRepository) UpdateBirth(entry *BirthEntrySave) (*BirthEntry, *apiError.APIError) {

	validateErr := validateUpdateBirth(r.DB, entry)
	if validateErr != nil {
		return nil, validateErr
	}

	query := `
		update animals
		set birth_date = :birth_date,
			sex = :sex,
			father_id = :father_id,
			observation = :observation
		where id = :id
	`
	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
		select 
			a.id,
			a.mother_id,
			m.name as mother_name,
			concat_ws(' - ', m.ring_number, m.name) as mother_info,
			a.birth_date as calf_birth_date,
			a.sex as calf_sex,
			case 
				when a.name is null then ''
				else concat_ws(' - ', a.ring_number, a.name)
			end as calf_name,
			a.father_id as calf_father_id,
			concat_ws(' - ', f.ring_number, f.name) calf_father,
			extract(days from a.birth_date - lag(a.birth_date) over win) as birth_interval
		from animals a
			join animals m on m.id = a.mother_id and m.animal_type <> 'OUTSIDE_ANIMAL'
			left join animals f on f.id = a.father_id
		where a.id = $1
		window win as (partition by a.mother_id order by a.birth_date)
	`

	result, err := repositoriesUtil.GetOne[BirthEntry](r.DB, selectQuery, entry.Id)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return result, nil
}

func (r *BirthRepository) ReplaceBirth(entry *BirthEntrySave) *apiError.APIError {

	query := `
		update animals
		set sex = :sex,
			father_id = :father_id,
			observation = :observation
		where mother_id = :mother_id 
			and birth_date = :birth_date
			and user_id = :user_id
	`
	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *BirthRepository) GetFather(entry *BirthEntrySave) (*BirthEntrySave, *apiError.APIError) {

	fatherId, err := getFatherId(r.DB, entry)
	if err != nil {
		return nil, err
	}

	if fatherId != "" {
		entry.FatherId = &fatherId
	}

	return entry, nil
}

func (r *BirthRepository) AddBirth(entry *BirthEntrySave) *apiError.APIError {

	validateErr := validateAddBirth(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	defer tx.Rollback()


	birthQuery := `
		insert into animals (ring_number, sex, birth_date, father_id, mother_id, animal_type, observation, user_id)
		values (:ring_number, :sex, :birth_date, :father_id, :mother_id, 'OFFSPRING', :observation, :user_id)
	`
	if entry.RingNumber == nil {
		birthQuery = `
			insert into animals (ring_number, sex, birth_date, father_id, mother_id, animal_type, observation, user_id)
			values (
				(select ring_number from animals where id = :mother_id and user_id = :user_id), 
				:sex, 
				:birth_date, 
				:father_id, 
				:mother_id, 
				'OFFSPRING', 
				:observation, 
				:user_id
			)
		`
	}

	newId, err := repositoriesUtil.NamedExecReturningIdTx(tx, birthQuery, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pastureQuery := `
		select pasture_id
		from pasture_entries
		where animal_id = $1
			and exit_date is null
			and deleted_at is null
		order by entry_date desc
		limit 1
	`

	var pastureId sql.NullString
	err = repositoriesUtil.GetPrimitiveTx(tx, pastureQuery, &pastureId, entry.MotherId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if !pastureId.Valid {
		tx.Commit()
		return nil
	}

	pastureEntry := &pastureEntries.PastureEntry{
		AnimalId:  newId,
		PastureId: pastureId.String,
		EntryDate: entry.BirthDate,
		UserId:    entry.UserId,
	}

	pastureEntryQuery := `
		insert into pasture_entries (animal_id, pasture_id, entry_date, user_id)
		values (
			:animal_id,
			:pasture_id,
			:entry_date,
			:user_id
		)
	`
	err = repositoriesUtil.NamedExecTx(tx, pastureEntryQuery, pastureEntry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	tx.Commit()
	return nil
}

func (r *BirthRepository) AddBirthNoValidation(entry *BirthEntrySave) *apiError.APIError {

	tx, err := r.DB.Beginx()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	defer tx.Rollback()


	birthQuery := `
		insert into animals (ring_number, sex, birth_date, father_id, mother_id, animal_type, observation, user_id)
		values (:ring_number, :sex, :birth_date, :father_id, :mother_id, 'OFFSPRING', :observation, :user_id)
	`
	newId, err := repositoriesUtil.NamedExecReturningIdTx(tx, birthQuery, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pastureQuery := `
		select pasture_id
		from pasture_entries
		where animal_id = $1
			and exit_date is null
			and deleted_at is null
		order by entry_date desc
		limit 1
	`

	var pastureId sql.NullString
	err = repositoriesUtil.GetPrimitiveTx(tx, pastureQuery, &pastureId, entry.MotherId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	if !pastureId.Valid {
		tx.Commit()
		return nil
	}

	pastureEntry := &pastureEntries.PastureEntry{
		AnimalId:  newId,
		PastureId: pastureId.String,
		EntryDate: entry.BirthDate,
		UserId:    entry.UserId,
	}

	pastureEntryQuery := `
		insert into pasture_entries (animal_id, pasture_id, entry_date, user_id)
		values (
			:animal_id,
			:pasture_id,
			:entry_date,
			:user_id
		)
	`
	err = repositoriesUtil.NamedExecTx(tx, pastureEntryQuery, pastureEntry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	tx.Commit()
	return nil
}
