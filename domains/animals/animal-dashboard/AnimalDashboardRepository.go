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
        WITH min_max as (
            select 
                min(make_date(extract(year from entry_date)::int,12,31)) as min_date, 
                max(make_date(extract(year from entry_date)::int,12,31)) as max_date 
            from animal_entries
        ),
        date_series AS (select generate_series(min_date, max_date, interval '1 year') as year from min_max)
        select 
            EXTRACT(YEAR FROM date_series.year) as year,
            count(animal_id) as total_animals
        FROM animal_entries as entries
            JOIN date_series ON entries.entry_date <= date_series.year
            AND (entries.exit_date IS NULL OR entries.exit_date > date_series.year)
    `

	props := repositoriesUtil.GroupByProps[TotalByYear]{
		Query:     query,
		TableName: "entries",
		GroupBy:   "year",
		OrderBy:   "year",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}

	return repositoriesUtil.GetGroupByResults(props)
}

func (r *DashboardRepository) TotalBySex(userId string, filter DashboardFilter) (*TotalBySex, error) {
	isActive := true
	filter.IsFiltered = true
	filter.IsActive = &isActive
	query := `
        select
            COUNT(animals.id) as total_animals,
            COUNT(animals.id) FILTER (WHERE animals.sex = 'F') as total_females,
            COUNT(animals.id) FILTER (WHERE animals.sex = 'M') as total_males
        FROM animals
    `
	props := repositoriesUtil.TotalProps[TotalBySex]{
		Query:     query,
		TableName: "animals",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}

	return repositoriesUtil.GetTotalResults(props)
}

func (r *DashboardRepository) TotalByType(userId string, filter DashboardFilter) (*AnimalByType, error) {
	query := `
        select
            COUNT(animals.id) FILTER (WHERE animals.type = 'BEEF_CATTLE') as beef_cattle,
            COUNT(animals.id) FILTER (WHERE animals.type = 'DAIRY_CATTLE') as dairy_cattle, 
            COUNT(animals.id) FILTER (WHERE animals.type = 'REPRODUCTION_ANIMALS') as reproduction_animals, 
            COUNT(animals.id) FILTER (WHERE animals.type = 'OFFSPRING') as offspring
        FROM animals
    `
	props := repositoriesUtil.TotalProps[AnimalByType]{
		Query:     query,
		TableName: "animals",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetTotalResults(props)
}

func (r *DashboardRepository) GroupByAgeAndFarm(userId string, filter DashboardFilter) (*[]AnimalsByAgeAndFarm, error) {
	isActive := true
	filter.IsFiltered = true
	filter.IsActive = &isActive
	query := `
        select 
            farms.id as farm_id,
            farms.name AS farm_name,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) < interval '3 months'
                AND animals.sex = 'M'
            ) AS newborn_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) < interval '3 months'
                AND animals.sex = 'F'
            ) AS newborn_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '3 months' and interval '9 months'
                AND animals.sex = 'M'
            ) AS baby_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '3 months' and interval '9 months'
                AND animals.sex = 'F'
            ) AS baby_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '9 months' and interval '13 months'
                AND animals.sex = 'M'
            ) AS child_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '9 months' and interval '13 months'
                AND animals.sex = 'F'
            ) AS child_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '13 months' and interval '25 months'
                AND animals.sex = 'M'
            ) AS young_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '13 months' and interval '25 months'
                AND animals.sex = 'F'
            ) AS young_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '25 months' and interval '37 months'
                AND animals.sex = 'M'
            ) AS adult_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '25 months' and interval '37 months'
                AND animals.sex = 'F'
            ) AS adult_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '37 months' 
                AND animals.sex = 'M'
            ) AS old_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '37 months' 
                AND animals.sex = 'F'
            ) AS old_female,
            COUNT(animals.id) FILTER (
                WHERE animals.sex = 'M'
            ) AS total_male,
            COUNT(animals.id) FILTER (
                WHERE animals.sex = 'F'
            ) AS total_female,
            COUNT(animals.id) as total 
        FROM animals
        LEFT JOIN farms ON farms.id = animals.farm_id
    `
	props := repositoriesUtil.GroupByProps[AnimalsByAgeAndFarm]{
		Query:     query,
		TableName: "animals",
		GroupBy:   "farms.name, farms.id",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetGroupByResults(props)
}

func (r *DashboardRepository) GroupByAgeAndPasture(userId string, filter DashboardFilter) (*[]AnimalsByAgeAndFarm, error) {
	isActive := true
	filter.IsFiltered = true
	filter.IsActive = &isActive
	query := `
        select 
            pastures.id as farm_id,
            pastures.name AS farm_name,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) < interval '3 months'
                AND animals.sex = 'M'
            ) AS newborn_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) < interval '3 months'
                AND animals.sex = 'F'
            ) AS newborn_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '3 months' and interval '9 months'
                AND animals.sex = 'M'
            ) AS baby_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '3 months' and interval '9 months'
                AND animals.sex = 'F'
            ) AS baby_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '9 months' and interval '13 months'
                AND animals.sex = 'M'
            ) AS child_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '9 months' and interval '13 months'
                AND animals.sex = 'F'
            ) AS child_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '13 months' and interval '25 months'
                AND animals.sex = 'M'
            ) AS young_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '13 months' and interval '25 months'
                AND animals.sex = 'F'
            ) AS young_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '25 months' and interval '37 months'
                AND animals.sex = 'M'
            ) AS adult_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) between interval '25 months' and interval '37 months'
                AND animals.sex = 'F'
            ) AS adult_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '37 months' 
                AND animals.sex = 'M'
            ) AS old_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '37 months' 
                AND animals.sex = 'F'
            ) AS old_female,
            COUNT(animals.id) FILTER (
                WHERE animals.sex = 'M'
            ) AS total_male,
            COUNT(animals.id) FILTER (
                WHERE animals.sex = 'F'
            ) AS total_female,
            COUNT(animals.id) as total 
        FROM animals
        LEFT JOIN pastures ON pastures.id = animals.pasture_id
    `
	props := repositoriesUtil.GroupByProps[AnimalsByAgeAndFarm]{
		Query:     query,
		TableName: "animals",
		GroupBy:   "pastures.name, pastures.id",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetGroupByResults(props)
}

func (r *DashboardRepository) GroupByAge(userId string, filter DashboardFilter) (*[]AnimalsByAge, error) {
	query := `
        select 
            CASE 
                WHEN age(animals.birth_date) < interval '3 months' THEN '0-2 meses'
                WHEN age(animals.birth_date) BETWEEN interval '3 months' AND interval '9 months' THEN '3-8 meses'
                WHEN age(animals.birth_date) BETWEEN interval '9 months' AND interval '13 months' THEN '9-12 meses'
                WHEN age(animals.birth_date) BETWEEN interval '13 months' AND interval '25 months' THEN '13-24 meses'
                WHEN age(animals.birth_date) BETWEEN interval '25 months' AND interval '37 months' THEN '25-36 meses'
                WHEN age(animals.birth_date) > interval '37 months' THEN '+36 meses'
                ELSE 'Desconhecido'
            END AS age_category,
            MAX(animals.birth_date) as max_birth_date,
            MIN(animals.birth_date) as min_birth_date,
            COUNT(animals.id) FILTER (WHERE animals.sex = 'M') as male,
            COUNT(animals.id) FILTER (WHERE animals.sex = 'F') as female
        FROM animals
    `
	props := repositoriesUtil.GroupByProps[AnimalsByAge]{
		Query:     query,
		TableName: "animals",
		GroupBy:   "age_category",
		OrderBy:   "min_birth_date DESC",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetGroupByResults(props)
}
