package animals

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type AnimalRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *AnimalRepository {
	return &AnimalRepository{db}
}

func (r *AnimalRepository) GetBirthHist(userId string) (*CardEntry, error) {

	query := `
		with cte as (
			select
				date_trunc('month', birth_date) as entry_date,
				count(*) as animals_number
			from animals
			where user_id = $1 
				and deleted_at is null
				and birth_date is not null
			group by 1
			order by 1 desc
			limit 12
		)
		select * from cte order by entry_date
	`

	hist, err := repositoriesUtil.GetList[AnimalsNumberHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	var current, past, trend int

	histEntries := *hist
	switch lenght := len(histEntries); lenght {
	case 0:
		current = 0
		past = 0
		trend = 0
	case 1:
		current = histEntries[0].AnimalsNumber
		past = 0
		trend = 0
	default:
		current = histEntries[lenght-1].AnimalsNumber
		past = histEntries[lenght-2].AnimalsNumber
		trend = current - past
	}

	cardEntry := &CardEntry{
		Current: current,
		Trend:   float64(trend),
		Hist:    histEntries,
	}

	return cardEntry, nil
}

func (r *AnimalRepository) GetDairyHist(userId string) (*CardEntry, error) {

	query := `
		with calendar as (
			select generate_series(
				date_trunc('month', max(end_date) - interval '12 months'),
				date_trunc('month', max(end_date)),
				interval '1 month'
			) as entry_date
			from lactations
		)
		select
			c.entry_date,
			count(*) as animals_number
		from lactations 
			join calendar c on start_date <= c.entry_date
				and c.entry_date <= coalesce(end_date, now())
		where user_id = $1 and deleted_at is null
		group by 1
		order by 1
	`

	hist, err := repositoriesUtil.GetList[AnimalsNumberHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	var current, past int
	var trend float64

	histEntries := *hist
	switch lenght := len(histEntries); lenght {
	case 0:
		current = 0
		past = 0
		trend = 0
	case 1:
		current = histEntries[0].AnimalsNumber
		past = 0
		trend = 0
	default:
		current = histEntries[lenght-1].AnimalsNumber
		past = histEntries[lenght-2].AnimalsNumber
		trend = util.CalculatePercentageTrend(float64(current), float64(past))
	}

	cardEntry := &CardEntry{
		Current: current,
		Trend:   trend,
		Hist:    histEntries,
	}

	return cardEntry, nil
}

func (r *AnimalRepository) GetDeathHist(userId string) (*CardEntry, error) {

	query := `
		with cte as (
			select
				date_trunc('month', death_date) as entry_date,
				count(*) as animals_number
			from animals a
			where user_id = $1 
				and deleted_at is null
				and death_date is not null
				and not exists (
					select 1
					from slaughter_entries s
					where s.animal_id = a.id
						and a.death_date = s.entry_date
						and s.user_id = $1
				)
			group by 1
			order by 1 desc
			limit 12
		)
		select *
		from cte 
		order by entry_date
	`

	hist, err := repositoriesUtil.GetList[AnimalsNumberHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	var current, past, trend int

	histEntries := *hist
	switch lenght := len(histEntries); lenght {
	case 0:
		current = 0
		past = 0
		trend = 0
	case 1:
		current = histEntries[0].AnimalsNumber
		past = 0
		trend = 0
	default:
		current = histEntries[lenght-1].AnimalsNumber
		past = histEntries[lenght-2].AnimalsNumber
		trend = current - past
	}

	cardEntry := &CardEntry{
		Current: current,
		Trend:   float64(trend),
		Hist:    histEntries,
	}

	return cardEntry, nil
}

func (r *AnimalRepository) GetSlaughterHist(userId string) (*CardEntry, error) {

	query := `
		with cte as (
			select
				entry_date,
				count(*) as animals_number
			from slaughter_entries
			where user_id = $1 and deleted_at is null
			group by 1
			order by 1 desc
			limit 12
		)
		select *
		from cte
		order by entry_date
	`

	hist, err := repositoriesUtil.GetList[AnimalsNumberHist](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	var current, past int
	var trend float64

	histEntries := *hist
	switch lenght := len(histEntries); lenght {
	case 0:
		current = 0
		past = 0
		trend = 0
	case 1:
		current = histEntries[0].AnimalsNumber
		past = 0
		trend = 0
	default:
		current = histEntries[lenght-1].AnimalsNumber
		past = histEntries[lenght-2].AnimalsNumber
		trend = util.CalculatePercentageTrend(float64(current), float64(past))
	}

	cardEntry := &CardEntry{
		Current: current,
		Trend:   trend,
		Hist:    histEntries,
	}

	return cardEntry, nil
}

func (r *AnimalRepository) GetAnimalTypes(userId string) (*AnimalByType, error) {

	query := `
		select
			count(*) filter (where animal_type = 'DAIRY_ANIMAL') as dairy_animals,
			count(*) filter (where animal_type = 'OFFSPRING') as offspring,
			count(*) filter (where animal_type = 'BEEF_ANIMAL') as beef_animals,
			count(*) filter (where animal_type = 'REPRODUCTION_ANIMAL') as reproduction_animals
		from animals a
		where user_id = $1 
			and deleted_at is null
			and is_outside_animal = false
			and death_date is null
	`
	result, err := repositoriesUtil.GetOne[AnimalByType](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AnimalRepository) GetLastDeaths(userId string) (*[]Animal, error) {

	query := `
		select
			id,
			concat_ws(
				' - ', 
				ring_number, 
				coalesce(name, sex),
				to_char(birth_date, 'DD/MM/YYYY')
			) as name,
			sex,
			animal_type,
			death_date,
			observation
		from animals a
		where user_id = $1 
			and deleted_at is null
			and is_outside_animal = false
			and death_date is not null
			and not exists (
				select 1
				from slaughter_entries s
				where s.animal_id = a.id
					and s.user_id = $1
					and s.deleted_at is null
			)
		order by death_date desc
		limit 20
	`
	result, err := repositoriesUtil.GetList[Animal](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AnimalRepository) GetAnimalsByAge(userId string) (*AnimalByType, error) {

	query := `
		select
			count(*) filter (where age(birth_date) < interval '2 months' and sex = 'M') as calf_male,
			count(*) filter (where age(birth_date) < interval '2 months' and sex = 'F') as calf_female,
			count(*) filter (where age(birth_date) between interval '2 months' and interval interval '8 months' and sex = 'M') as young_male,
			count(*) filter (where age(birth_date) between interval '2 months' and interval interval '8 months' and sex = 'F') as young_female,
		from animals a
		where user_id = $1 
			and deleted_at is null
			and is_outside_animal = false
			and death_date is null
	`
	result, err := repositoriesUtil.GetOne[AnimalByType](r.DB, query, userId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *AnimalRepository) GroupByYear(userId string) (*[]TotalByYear, error) {
	query := `
        with min_max as (
            select 
                min(make_date(extract(year from entry_date)::int, 12, 31)) as min_date, 
                max(make_date(extract(year from entry_date)::int, 12, 31)) as max_date 
            from animal_entries
        ),
        date_series as (select generate_series(min_date, max_date, interval '1 year') as year from min_max)
        select 
            extract(year from date_series.year) as year,
            count(animal_id) as total_animals
        from animal_entries as entries
            join date_series on entries.entry_date <= date_series.year
            and (entries.exit_date is null or entries.exit_date > date_series.year)
		where animals.deleted_at is null and animals.user_id = $1
		order by entries.entry_date
		group by extract(year from date_series.year)
    `
	return repositoriesUtil.GetList[TotalByYear](r.DB, query, userId)
}

func (r *AnimalRepository) TotalBySex(userId string) (*TotalBySex, error) {
	query := `
        select
            count(animals.id) as total_animals,
            count(animals.id) filter (where animals.sex = 'F') as total_females,
            count(animals.id) filter (where animals.sex = 'M') as total_males
        from animals
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
    `
	return repositoriesUtil.GetOne[TotalBySex](r.DB, query, userId)
}

func (r *AnimalRepository) TotalByType(userId string) (*AnimalByType, error) {
	query := `
        select
            count(animals.id) filter (where animals.animal_type = 'BEEF_ANIMAL') as beef_cattle,
            count(animals.id) filter (where animals.animal_type = 'DAIRY_ANIMAL') as dairy_cattle, 
            count(animals.id) filter (where animals.animal_type = 'REPRODUCTION_ANIMAL') as reproduction_animals, 
            count(animals.id) filter (where animals.animal_type = 'OFFSPRING') as offspring
        from animals
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
    `
	return repositoriesUtil.GetOne[AnimalByType](r.DB, query, userId)
}

func (r *AnimalRepository) GetAgeAndSex(userId string) (*[]AnimalsByAge, error) {
	query := `
        select 
            case 
                when age(birth_date) < interval '2 months' then '0-2 meses'
                when age(birth_date) between interval '2 months' and interval '8 months' then '2-8 meses'
                when age(birth_date) between interval '8 months' and interval '12 months' then '8-12 meses'
                when age(birth_date) between interval '12 months' and interval '24 months' then '12-24 meses'
                when age(birth_date) between interval '24 months' and interval '36 months' then '24-36 meses'
                when age(birth_date) > interval '36 months' then '+36 meses'
				else 'Sem Data'
            end as category,
			count(*) filter (where sex = 'M') as male,
			count(*) filter (where sex = 'F') as female
        from animals
        where user_id = $1
            and deleted_at is null
			and is_outside_animal = false
			and birth_date is not null
			and death_date is not null
		group by category
		order by min(birth_date) desc
    `
	return repositoriesUtil.GetList[AnimalsByAge](r.DB, query, userId)
}

func (r *AnimalRepository) FindPage(
	userId string,
	cursor string,
	sort string,
	order string,
	filter AnimalFilter,
) (page *entity.Page[Animal], err error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"name":                   {Field: "coalesce(a.name, '')", Order: "asc"},
		"average_birth_interval": {Field: "coalesce(b.average_birth_interval, 0)", Order: "asc"},
		"average_lac_interval":   {Field: "coalesce(l.average_lac_interval, 0)", Order: "asc"},
		"average_prod":           {Field: "coalesce(milk.average_prod, 0)", Order: "asc"},
		"average_peak":           {Field: "coalesce(l.average_peak, 0)", Order: "asc"},
		"death_date":             {Field: "coalesce(a.death_date, '-infinity')", Order: "asc"},
		"weaning_date":           {Field: "coalesce(a.weaning_date, '-infinity')", Order: "asc"},
		"birth_date":             {Field: "coalesce(a.birth_date, '-infinity')", Order: "desc"},
		"animal_order":           {Field: "coalesce(nullif(regexp_replace(a.ring_number, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"created_at":             {Field: "a.created_at", Order: "asc"},
		"id":                     {Field: "a.id", Order: "asc"},
	}

	query := `
		with peak_stats as (
			select
				l.id as lactation_id,
				max(m.quantity) as peak
			from milk_entries m
			join lactations l
				on l.animal_id = m.animal_id
			   and l.start_date <= m.entry_date
			   and coalesce(l.end_date, now()) >= m.entry_date
			   and l.deleted_at is null
			where m.deleted_at is null
			group by l.id
		),

		lac_interval_cte as (
			select
				l.animal_id,
				l.start_date,
				l.end_date,
				extract(days from l.start_date - lag(l.start_date) over (partition by l.animal_id order by l.start_date)) as lac_interval,
				p.peak
			from lactations l
			join peak_stats p on p.lactation_id = l.id
			where l.user_id = $1 and l.deleted_at is null
		),

		lac_stats as (
			select
				animal_id,
				avg(lac_interval) as average_lac_interval,
				avg(peak) as average_peak
			from lac_interval_cte
			group by animal_id
		),

		birth_interval_cte as (
			select
				mother_id,
				extract(days from birth_date - lag(birth_date) over (partition by mother_id order by birth_date)) as birth_interval
			from animals
			where user_id = $1
			  and deleted_at is null
			  and mother_id is not null
		),

		birth_stats as (
			select
				mother_id,
				avg(birth_interval) as average_birth_interval
			from birth_interval_cte
			group by mother_id
		)

		select
			a.id,
			a.father_id,
			a.mother_id,
			coalesce(nullif(regexp_replace(a.ring_number, '[^0-9]', '', 'g'), '')::int, 0) animal_order,
			a.ring_number,
			a.name,
			a.sex,
			a.birth_date,
			concat_ws(' - ', f.ring_number, f.name) as father_name,
			concat_ws(' - ', m.ring_number, m.name) as mother_name,
			a.weight_birth,
			a.weaning_date,
			a.death_date,
			a.animal_type,
			milk.average_prod,
			l.average_lac_interval,
			l.average_peak,
			b.average_birth_interval,
			a.observation,
			a.created_at
		from animals a
			left join animals m on m.id = a.mother_id
			left join animals f on f.id = a.father_id
			left join (
				select
					animal_id,
					avg(quantity) as average_prod
				from milk_entries
				where deleted_at is null
				group by animal_id
			) milk on milk.animal_id = a.id
			left join birth_stats b on b.mother_id = a.id
			left join lac_stats l on l.animal_id = a.id
	`

	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "a", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	mainExpression := "a.is_outside_animal = false and a.deleted_at is null and a.user_id = $1"
	whereExpression := repositoriesUtil.GetWhereExpression(mainExpression, filterExpression, cursorExpression)

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)
	orderExpression := " order by " + sortExpression
	query += whereExpression + orderExpression
	return repositoriesUtil.GetPage[Animal](r.DB, query, sort, 200, args...)
}

func (r *AnimalRepository) FindPageFooter(userId string, filter AnimalFilter) (*AnimalFooter, error) {

	query := `
		with lac_stats as (
			select
				l.animal_id,
				avg(l.start_date - lag(l.start_date) over (order by l.start_date partition by l.animal_id)) as average_lac_interval,
				avg(m.peak) as average_peak
			from lactations l
				join lateral (
					select max(quantity) peak
					from milk_entries 
					where animal_id = l.animal_id
						and entry_date >= l.start_date
						and entry_date <= coalesce(l.entry_date, now())
						and user_id = $1
						and deleted_at is null
				) m on true
			where l.user_id = $1 and l.deleted_at is null
			group by 1
		),
		cte as (
			select 
				a.id,
				a.name,
				milk.average_prod,
				l.average_lac_interval,
				l.average_peak,
				b.average_birth_interval,
				a.observation
			from animals a
				left join animals m on m.id = a.mother_id
				left join animals f on f.id = a.father_id
				left join (
					select
						animal_id,
						avg(quantity) average_prod
					from milk_entries
					where user_id = $1 and deleted_at is null
					group by 1
				) milk on milk.animal_id = a.id
				left join (
					select
						mother_id,
						avg(birth_date - lag(birth_date) over (order by birth_date partition by mother_id)) average_birth_interval
					from animals
					where user_id = $1 and deleted_at is null
					group by 1
				) b on b.mother_id = a.mother_id
				left join lac_stats l on l.animal_id = a.id
			where a.user_id = $1 and a.deleted_at is null
		)
		select 
			count(cte.id) as total,
			avg(cte.average_prod) as average_prod,
			avg(cte.average_lac_interval) as average_lac_interval,
			avg(cte.average_birth_interval) as average_birth_interval,
			avg(cte.average_peak) as average_peak
		from cte
	`

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := repositoriesUtil.GetWhereExpression(filterExpression)
	query += whereExpression

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[AnimalFooter](r.DB, query, args...)
}

func (r *AnimalRepository) FindById(id string, userId string) (*Animal, error) {
	query := `
        select 
			a.id,
			a.ring_number,
			a.name,
			a.sex,
			a.father_id,
			a.mother_id,
            concat_ws(' - ', f.ring_number, f.name) as father_name, 
            concat_ws(' - ', m.ring_number, m.name) as mother_name,
			a.animal_type,
			a.birth_date,
			a.death_date,
			a.weaning_date,
			a.observation,
			a.is_embryo_donor,
			a.is_transfer_bull,
			a.is_breeding_bull,
			a.is_insemination_bull,
			a.is_outside_animal,
            format('%s (%s)', p.name, fa.name) as pasture_name
        from animals a
            left join animals f ON f.id = a.father_id
            left join animals m ON m.id = a.mother_id
			left join lateral (
				select pe.pasture_id
				from pasture_entries pe
				where pe.user_id = $2 and pe.animal_id = $1
				order by entry_date desc
				limit 1
			) pe on true
            left join pastures p ON p.id = pe.pasture_id
            left join farms fa ON fa.id = p.farm_id
		where a.id = $1 and a.user_id = $2
	`
	return repositoriesUtil.GetOne[Animal](r.DB, query, id, userId)
}

func (r *AnimalRepository) SearchFather(userId string) (*[]entity.SearchEntity, error) {
	queryInput := `
        select id, concat_ws(' - ', ring_number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'M' 
            and animal_type = 'REPRODUCTION_ANIMAL'
            and (name is not null)
            and deleted_at is null
        order by coalesce(nullif(regexp_replace(ring_number, '[^0-9]', '', 'g'), '')::int, 0), coalesce(name, '')
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, queryInput, userId)
}

func (r *AnimalRepository) FindMaleOffspring(id string, userId string) (*[]entity.SearchEntity, error) {
	queryInput := `
        select 
			id, 
			concat_ws(
				' - ', 
				sex, 
				coalesce(to_char(birth_date, 'DD/MM/YYYY'), 'Desconhecido'),
				case 
					when death_date is not null then 'Morto'
				end
			) as label 
		from animals 
        where mother_id = $1
            and sex = 'M' 
            and animal_type = 'OFFSPRING'
			and death_date is null
			and user_id = $2
            and deleted_at is null
        order by birth_date desc
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, queryInput, id, userId)
}

func (r *AnimalRepository) FindFemaleOffspring(id string, userId string) (*[]entity.SearchEntity, error) {
	queryInput := `
        select 
			id, 
			concat_ws(
				' - ', 
				sex, 
				coalesce(to_char(birth_date, 'DD/MM/YYYY'), 'Desconhecido'),
				case 
					when death_date is not null then 'Morto'
				end
			) as label 
		from animals 
        where mother_id = $1
            and sex = 'F' 
            and animal_type = 'OFFSPRING'
			and user_id = $2
            and deleted_at is null
        order by birth_date desc
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, queryInput, id, userId)
}

func (r *AnimalRepository) SearchAnimals(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select 
			id, 
			concat_ws(
				' - ', 
				ring_number, 
				name, 
				to_char(birth_date, 'DD/MM/YYYY'),
				case 
					when death_date is not null then 'Morto'
				end
			) as label 
		from animals 
        where user_id = $1 
            and is_outside_animal = false
            and deleted_at is null
        order by 
			coalesce(nullif(regexp_replace(ring_number, '[^0-9]', '', 'g'), '')::int, 0), 
			coalesce(name, ''), 
			coalesce(birth_date, '-infinity')
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) SearchAllMothers(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select 
			id, 
			concat_ws(
				' - ', 
				ring_number, 
				name,
				case 	
					when death_date is not null then 'Morto'
				end
			) as label 
            from animals 
        where user_id = $1 
            and sex = 'F' 
            and animal_type in ('REPRODUCTION_ANIMAL', 'DAIRY_ANIMAL') 
            and name is not null
            and deleted_at is null
        order by coalesce(nullif(regexp_replace(ring_number, '[^0-9]', '', 'g'), '')::int, 0), label
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) SearchMother(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
		from animals 
        where user_id = $1 
            and sex = 'F' 
            and animal_type in ('REPRODUCTION_ANIMAL', 'DAIRY_ANIMAL') 
			and name is not null 
			and is_outside_animal = false
			and deleted_at is null
        order by coalesce(nullif(regexp_replace(ring_number, '[^0-9]', '', 'g'), '')::int, 0), label
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) SearchDairyAnimals(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
		from animals 
        where user_id = $1 
            and sex = 'F' 
            and animal_type = 'DAIRY_ANIMAL'
            and deleted_at is null
        order by coalesce(nullif(regexp_replace(ring_number, '[^0-9]', '', 'g'), '')::int, 0), name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) SearchBull(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
		from animals 
        where user_id = $1 
            and sex = 'M' 
			and animal_type = 'REPRODUCTION_ANIMAL'
            and name is not null
            and deleted_at is null
        order by 
			coalesce(nullif(regexp_replace(ring_number, '[^0-9]' ,'', 'g'), '')::int, 0), 
			name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) Delete(id string, userId string) *apiError.APIError {

	validateErr := validDelete(r.DB, id, userId)
	if validateErr != nil {
		return validateErr
	}

	tx, err := r.DB.Beginx()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		update animals
		set deleted_at = now()
		where id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	lacQuery := `
		update lactations
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, lacQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	milkQuery := `
		update milk_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, milkQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	embryoQuery := `
		update embryo_transfer
		set deleted_at = now()
		where receiver_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, embryoQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	inseminationQuery := `
		update insemination_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, inseminationQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	breedingQuery := `
		update natural_breedings
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, breedingQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pastureQuery := `
		update pasture_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, pastureQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pregnancyQuery := `
		update pregnancy_tests
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, pregnancyQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	slaughterQuery := `
		update slaughter_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, slaughterQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	weightQuery := `
		update weight_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, weightQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *AnimalRepository) DeleteNoValidation(id string, userId string) *apiError.APIError {

	tx, err := r.DB.Beginx()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	defer tx.Rollback()

	query := `
		update animals
		set deleted_at = now()
		where id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, query, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	lacQuery := `
		update lactations
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, lacQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	milkQuery := `
		update milk_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, milkQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	embryoQuery := `
		update embryo_transfer
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, embryoQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	inseminationQuery := `
		update embryo_transfer
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, inseminationQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	breedingQuery := `
		update natural_breedings
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, breedingQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pastureQuery := `
		update natural_breedings
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, pastureQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	pregnancyQuery := `
		update pregnancy_tests
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, pregnancyQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	slaughterQuery := `
		update slaughter_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, slaughterQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	weightQuery := `
		update weight_entries
		set deleted_at = now()
		where animal_id = $1 and user_id = $2
	`

	err = repositoriesUtil.ExecTx(tx, weightQuery, id, userId)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	err = tx.Commit()
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *AnimalRepository) Update(newEntry *AnimalSave) (*Animal, *apiError.APIError) {

	validateErr := validateUpdate(r.DB, newEntry)
	if validateErr != nil {
		return nil, validateErr
	}

	updateQuery := `
		update animals
		set ring_number = :ring_number,
			name = :name,
			death_date = :death_date,
			birth_date = :birth_date,
			sex = :sex,
			father_id = :father_id,
			mother_id = :mother_id,
			animal_type = :animal_type,
			weaning_date = :weaning_date,
			is_outside_animal = :is_outside_animal,
			is_insemination_bull = :is_insemination_bull,
			is_transfer_bull = :is_transfer_bull,
			is_breeding_bull = :is_breeding_bull,
			is_embryo_donor = :is_embryo_donor
		where id = :id and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, updateQuery, newEntry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
        select 
			a.id,
			a.ring_number,
			a.name,
			a.sex,
			a.father_id,
			a.mother_id,
            concat_ws(' - ', f.ring_number, f.name) as father_name, 
            concat_ws(' - ', m.ring_number, m.name) as mother_name,
			a.animal_type,
			a.birth_date,
			a.death_date,
			a.weaning_date,
			a.observation,
            format('%s (%s)', p.name, fa.name) as pasture_name
        from animals a
            left join animals f ON f.id = a.father_id
            left join animals m ON m.id = a.mother_id
			left join lateral (
				select pe.pasture_id
				from pasture_entries pe
				where pe.user_id = $2 and pe.animal_id = $1
				order by entry_date desc
				limit 1
			) pe on true
            left join pastures p ON p.id = pe.pasture_id
            left join farms fa ON fa.id = p.farm_id
		where a.id = $1 and a.user_id = $2
    `

	response, err := repositoriesUtil.GetOne[Animal](r.DB, selectQuery, newEntry.Id, newEntry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil

}

func (r *AnimalRepository) UpdateNoValidation(newEntry *AnimalSave) (*Animal, *apiError.APIError) {

	updateQuery := `
		update animals
		set ring_number = :ring_number,
			name = :name,
			death_date = :death_date,
			birth_date = :birth_date,
			sex = :sex,
			father_id = :father_id,
			mother_id = :mother_id,
			animal_type = :animal_type,
			weaning_date = :weaning_date,
			is_outside_animal = :is_outside_animal,
			is_insemination_bull = :is_insemination_bull,
			is_transfer_bull = :is_transfer_bull,
			is_breeding_bull = :is_breeding_bull,
			is_embryo_donor = :is_embryo_donor
		where id = :id and user_id = :user_id
	`

	err := repositoriesUtil.NamedExec(r.DB, updateQuery, newEntry)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	selectQuery := `
        select 
			a.id,
			a.ring_number,
			a.name,
			a.sex,
			a.father_id,
			a.mother_id,
            concat_ws(' - ', f.ring_number, f.name) as father_name, 
            concat_ws(' - ', m.ring_number, m.name) as mother_name,
			a.animal_type,
			a.birth_date,
			a.death_date,
			a.weaning_date,
			a.observation,
            format('%s (%s)', p.name, fa.name) as pasture_name
        from animals a
            left join animals f ON f.id = a.father_id
            left join animals m ON m.id = a.mother_id
			left join lateral (
				select pe.pasture_id
				from pasture_entries pe
				where pe.user_id = $2 and pe.animal_id = $1
				order by entry_date desc
				limit 1
			) pe on true
            left join pastures p ON p.id = pe.pasture_id
            left join farms fa ON fa.id = p.farm_id
		where a.id = $1 and a.user_id = $2
    `

	response, err := repositoriesUtil.GetOne[Animal](r.DB, selectQuery, newEntry.Id, newEntry.UserId)
	if err != nil {
		return nil, apiError.InternalServerAPIError(err)
	}

	return response, nil

}

func (r *AnimalRepository) Add(entry *AnimalSave) *apiError.APIError {

	validateErr := validateAdd(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		insert into animals (
			ring_number, 
			name, 
			sex, 
			father_id, 
			mother_id, 
			animal_type, 
			birth_date, 
			death_date, 
			weaning_date, 
			is_insemination_bull,
			is_transfer_bull,
			is_breeding_bull,
			is_embryo_donor,
			is_outside_animal,
			observation,
			user_id
		)
		values (
			:ring_number, 
			:name, 
			:sex, 
			:father_id, 
			:mother_id, 
			:animal_type, 
			:birth_date, 
			:death_date, 
			:weaning_date, 
			:is_insemination_bull,
			:is_transfer_bull,
			:is_breeding_bull,
			:is_embryo_donor,
			:is_outside_animal,
			:observation,
			:user_id
		)
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *AnimalRepository) AddNoValidation(entry *AnimalSave) *apiError.APIError {

	query := `
		insert into animals (
			ring_number, 
			name, 
			sex, 
			father_id, 
			mother_id, 
			animal_type, 
			birth_date, 
			death_date, 
			weaning_date, 
			is_insemination_bull,
			is_transfer_bull,
			is_breeding_bull,
			is_embryo_donor,
			is_outside_animal,
			observation,
			user_id
		)
		values (
			:ring_number, 
			:name, 
			:sex, 
			:father_id, 
			:mother_id, 
			:animal_type, 
			:birth_date, 
			:death_date, 
			:weaning_date, 
			:is_insemination_bull,
			:is_transfer_bull,
			:is_breeding_bull,
			:is_embryo_donor,
			:is_outside_animal,
			:observation,
			:user_id
		)
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}

func (r *AnimalRepository) Replace(entry *AnimalSave) *apiError.APIError {

	validateErr := validateReplace(r.DB, entry)
	if validateErr != nil {
		return validateErr
	}

	query := `
		insert into animals (
			ring_number, 
			name, 
			sex, 
			father_id, 
			mother_id, 
			animal_type, 
			birth_date, 
			death_date, 
			weaning_date, 
			is_insemination_bull,
			is_transfer_bull,
			is_breeding_bull,
			is_embryo_donor,
			is_outside_animal,
			observation,
			user_id
		)
		values (
			:ring_number, 
			:name, 
			:sex, 
			:father_id, 
			:mother_id, 
			:animal_type, 
			:birth_date, 
			:death_date, 
			:weaning_date, 
			:is_insemination_bull,
			:is_transfer_bull,
			:is_breeding_bull,
			:is_embryo_donor,
			:is_outside_animal,
			:observation,
			:user_id
		)
	`

	err := repositoriesUtil.NamedExec(r.DB, query, entry)
	if err != nil {
		return apiError.InternalServerAPIError(err)
	}

	return nil
}
