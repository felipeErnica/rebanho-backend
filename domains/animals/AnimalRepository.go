package animals

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type AnimalRepository struct {
	SelectQuery string
	TableName   string
	DB          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *AnimalRepository {
	selectQuery := `
        select animals.*, 
            coalesce(nullif(regexp_replace(animals.ring_number, '[^0-9]', '', 'g'), '')::int, 0) as animal_order,
            concat_ws(' - ', father.ring_number, father.name) as father_name, 
            concat_ws(' - ', mother.ring_number, mother.name) as mother_name,
            pastures.name as pasture_name,
            farms.id as farm_id, farms.name as farm_name
        from animals
            left join animals as father ON father.id = animals.father_id
            left join animals as mother ON mother.id = animals.mother_id
            left join pastures ON pastures.id = animals.pasture_id
            left join farms ON farms.id = pastures.farm_id
    `
	return &AnimalRepository{selectQuery, "animals", db}
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

func (r *AnimalRepository) GroupByAgeAndFarm(userId string) (*[]AnimalsByAgeAndFarm, error) {
	query := `
        select 
            farms.id as farm_id,
            farms.name as farm_name,
            count(animals.id) filter (
                where age(animals.birth_date) < interval '3 months'
                and animals.sex = 'M'
            ) as newborn_male,
            count(animals.id) filter (
                where age(animals.birth_date) < interval '3 months'
                and animals.sex = 'F'
            ) as newborn_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '3 months' and interval '9 months'
                and animals.sex = 'M'
            ) as baby_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '3 months' and interval '9 months'
                and animals.sex = 'F'
            ) as baby_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '9 months' and interval '13 months'
                and animals.sex = 'M'
            ) as child_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '9 months' and interval '13 months'
                and animals.sex = 'F'
            ) as child_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '13 months' and interval '25 months'
                and animals.sex = 'M'
            ) as young_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '13 months' and interval '25 months'
                and animals.sex = 'F'
            ) as young_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '25 months' and interval '37 months'
                and animals.sex = 'M'
            ) as adult_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '25 months' and interval '37 months'
                and animals.sex = 'F'
            ) as adult_female,
            count(animals.id) filter (
                where age(animals.birth_date) > interval '37 months' 
                and animals.sex = 'M'
            ) as old_male,
            count(animals.id) filter (
                where age(animals.birth_date) > interval '37 months' 
                and animals.sex = 'F'
            ) as old_female,
            count(animals.id) filter (
                where animals.sex = 'M'
            ) as total_male,
            count(animals.id) filter (
                where animals.sex = 'F'
            ) as total_female,
            count(animals.id) as total 
        from animals
			left join pastures on pastures.id = animals.pasture_id
			left join farms on farms.id = pastures.farm_id
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
		order by animals.birth_date
		group by farm.name, farm.id
    `
	return repositoriesUtil.GetList[AnimalsByAgeAndFarm](r.DB, query, userId)
}

func (r *AnimalRepository) GroupByAgeAndPasture(userId string) (*[]AnimalsByAgeAndFarm, error) {
	query := `
        select 
            pastures.id as farm_id,
            pastures.name as farm_name,
            count(animals.id) filter (
                where age(animals.birth_date) < interval '3 months'
                and animals.sex = 'M'
            ) as newborn_male,
            count(animals.id) filter (
                where age(animals.birth_date) < interval '3 months'
                and animals.sex = 'F'
            ) as newborn_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '3 months' and interval '9 months'
                and animals.sex = 'M'
            ) as baby_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '3 months' and interval '9 months'
                and animals.sex = 'F'
            ) as baby_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '9 months' and interval '13 months'
                and animals.sex = 'M'
            ) as child_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '9 months' and interval '13 months'
                and animals.sex = 'F'
            ) as child_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '13 months' and interval '25 months'
                and animals.sex = 'M'
            ) as young_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '13 months' and interval '25 months'
                and animals.sex = 'F'
            ) as young_female,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '25 months' and interval '37 months'
                and animals.sex = 'M'
            ) as adult_male,
            count(animals.id) filter (
                where age(animals.birth_date) between interval '25 months' and interval '37 months'
                and animals.sex = 'F'
            ) as adult_female,
            count(animals.id) filter (
                where age(animals.birth_date) > interval '37 months' 
                and animals.sex = 'M'
            ) as old_male,
            count(animals.id) filter (
                where age(animals.birth_date) > interval '37 months' 
                and animals.sex = 'F'
            ) as old_female,
            count(animals.id) filter (
                where animals.sex = 'M'
            ) as total_male,
            count(animals.id) filter (
                where animals.sex = 'F'
            ) as total_female,
            count(animals.id) as total 
        from animals
        left join pastures on pastures.id = animals.pasture_id
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
		order by animals.birth_date
		group by pastures.name, pastures.id
    `
	return repositoriesUtil.GetList[AnimalsByAgeAndFarm](r.DB, query, userId)
}

func (r *AnimalRepository) GroupByAge(userId string) (*[]AnimalsByAge, error) {
	query := `
        select 
            case 
                when age(animals.birth_date) < interval '3 months' then '0-2 meses'
                when age(animals.birth_date) between interval '3 months' and interval '9 months' then '3-8 meses'
                when age(animals.birth_date) between interval '9 months' and interval '13 months' then '9-12 meses'
                when age(animals.birth_date) between interval '13 months' and interval '25 months' then '13-24 meses'
                when age(animals.birth_date) between interval '25 months' and interval '37 months' then '25-36 meses'
                when age(animals.birth_date) > interval '37 months' then '+36 meses'
                else 'Desconhecido'
            end as age_category,
            max(animals.birth_date) as max_birth_date,
            min(animals.birth_date) as min_birth_date,
            count(animals.id) filter (where animals.sex = 'M') as male,
            count(animals.id) filter (where animals.sex = 'F') as female
        from animals
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
		order by animals.birth_date
		group by age_category
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
		"name":                   {Field: "coalesce(animals.name, '')", Order: "asc"},
		"isr":                    {Field: "coalesce(animals.isr, 0)", Order: "asc"},
		"average_birth_interval": {Field: "coalesce(animals.average_birth_interval, 0)", Order: "asc"},
		"average_prod_interval":  {Field: "coalesce(animals.average_prod_interval, 0)", Order: "asc"},
		"average_prod":           {Field: "coalesce(animals.average_prod, 0)", Order: "asc"},
		"average_peak":           {Field: "coalesce(animals.average_peak, 0)", Order: "asc"},
		"death_date":             {Field: "coalesce(animals.death_date, '-infinity')", Order: "asc"},
		"weaning_date":           {Field: "coalesce(animals.weaning_date, '-infinity')", Order: "asc"},
		"birth_date":             {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"animal_order":           {Field: "coalesce(nullif(regexp_replace(animals.ring_number, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
	}

	whereExpression := "where animals.deleted_at is null and animals.user_id = $1"
	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}
	sortExpression = " order by " + sortExpression

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression = whereExpression + " and " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression = whereExpression + " and " + cursorExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	query := r.SelectQuery + whereExpression + sortExpression
	return repositoriesUtil.GetPage[Animal](r.DB, query, sort, 200, args...)
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

func (r *AnimalRepository) FindByFatherId(fatherId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.father_id = $1"
	return repositoriesUtil.GetList[Animal](r.DB, query, fatherId)
}

func (r *AnimalRepository) FindByMotherId(motherId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.mother_id = $1"
	return repositoriesUtil.GetList[Animal](r.DB, query, motherId)
}

func (r *AnimalRepository) FindByName(name string, userId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.name = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[Animal](r.DB, query, name, userId)
}

func (r *AnimalRepository) FindByNumber(ringNumber string, userId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.ring_number = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[Animal](r.DB, query, ringNumber, userId)
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
