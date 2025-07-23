package animalDashboard

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type DashboardRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db}
}

func (r *DashboardRepository) GroupByYear(
	userId string,
	filter DashboardFilter,
) (*[]TotalByYear, error) {
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
    `

	whereExpression := "where animals.deleted_at is null and animals.user_id = $1"
	orderBy := " order by entries.entry_date"
    groupBy := " group by extract(year from date_series.year)"
	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}
	if filterExpression != "" {
		whereExpression = whereExpression + " and " + filterExpression
	}
	query = query + whereExpression + groupBy + orderBy

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetList[TotalByYear](r.DB, query, args...)
}

func (r *DashboardRepository) TotalBySex(userId string, filter DashboardFilter) (*TotalBySex, error) {
	query := `
        select
            count(animals.id) as total_animals,
            count(animals.id) filter (where animals.sex = 'F') as total_females,
            count(animals.id) filter (where animals.sex = 'M') as total_males
        from animals
    `
	where := `
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
    `
	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		where = where + " and " + filterExpression
	}
	query = query + where

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[TotalBySex](r.DB, query, args...)
}

func (r *DashboardRepository) TotalByType(userId string, filter DashboardFilter) (*AnimalByType, error) {
	query := `
        select
            count(animals.id) filter (where animals.animal_type = 'BEEF_ANIMAL') as beef_cattle,
            count(animals.id) filter (where animals.animal_type = 'DAIRY_ANIMAL') as dairy_cattle, 
            count(animals.id) filter (where animals.animal_type = 'REPRODUCTION_ANIMAL') as reproduction_animals, 
            count(animals.id) filter (where animals.animal_type = 'OFFSPRING') as offspring
        from animals
    `
	where := ` 
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
    `

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		where = where + " and " + filterExpression
	}
	query = query + where

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[AnimalByType](r.DB, query, args...)
}

func (r *DashboardRepository) GroupByAgeAndFarm(userId string, filter DashboardFilter) (*[]AnimalsByAgeAndFarm, error) {
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
    `
	where := ` 
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
    `
    orderBy := " order by animals.birth_date"
    groupBy := " group by farm.name, farm.id"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		where = where + " and " + filterExpression
	}
	query = query + where + groupBy + orderBy

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetList[AnimalsByAgeAndFarm](r.DB, query, args...)
}

func (r *DashboardRepository) GroupByAgeAndPasture(userId string, filter DashboardFilter) (*[]AnimalsByAgeAndFarm, error) {
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
    `
	where := ` 
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
    `
    orderBy := " order by animals.birth_date"
    groupBy := " group by pastures.name, pastures.id"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		where = where + " and " + filterExpression
	}
	query = query + where + groupBy + orderBy

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetList[AnimalsByAgeAndFarm](r.DB, query, args...)
}

func (r *DashboardRepository) GroupByAge(userId string, filter DashboardFilter) (*[]AnimalsByAge, error) {
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
    `
	where := ` 
        where animals.user_id = $1
            and animals.deleted_at is null
            and animals.animal_type not in ('DEAD_ANIMAL', 'SLAUGTHERED_ANIMAL', 'OUTSIDE_ANIMAL')
    `
    orderBy := " order by animals.birth_date"
    groupBy := " group by age_category"

	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		where = where + " and " + filterExpression
	}
	query = query + where + groupBy + orderBy

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetList[AnimalsByAge](r.DB, query, args...)
}
