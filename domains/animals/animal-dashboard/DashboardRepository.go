package animalDashboard

import (
	"time"

	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type DashboardRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db}
}

func (r *DashboardRepository) TotalByYear(userId string, minDate time.Time, maxDate time.Time, filter AnimalsDashboardFilter) (*TotalByYear, error) {
	query := `
        WITH date_series AS (
            SELECT generate_series($1, $2, interval '1 year') as year
        ),
        SELECT 
            date.year as year
            COUNT(animal_id) as total_animals
        FROM animal_entries as entries
            JOIN date_series as date ON entries.entry_date <= date.year + interval '1 year' - interval '1 day'
            AND (entries.exit_date IS NULL OR entries.exit_date > date.year + interval '1 year' - interval '1 day')
        GROUP BY date.year
    `
    total := &TotalByYear{}
    err := r.DB.Get(total, query)
    return  total, err
}

func (r *DashboardRepository) TotalBySex(userId string, filter AnimalsDashboardFilter) (*TotalBySex, error) {
	isActive := true
	filter.IsFiltered = true
	filter.IsActive = &isActive
	query := `
        SELECT
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

func (r *DashboardRepository) TotalByType(userId string, filter AnimalsDashboardFilter) (*AnimalByType, error) {
	query := `
        SELECT
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

func (r *DashboardRepository) GroupByAgeAndFarm(userId string, filter AnimalsDashboardFilter) (*[]AnimalsByAgeAndFarm, error) {
	isActive := true
	filter.IsFiltered = true
	filter.IsActive = &isActive
	query := `
        SELECT 
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
                WHERE age(animals.birth_date) BETWEEN interval '3 months' AND interval '8 months'
                AND animals.sex = 'M'
            ) AS baby_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '3 months' AND interval '8 months'
                AND animals.sex = 'F'
            ) AS baby_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '9 months' AND interval '12 months'
                AND animals.sex = 'M'
            ) AS child_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '9 months' AND interval '12 months'
                AND animals.sex = 'F'
            ) AS child_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '13 months' AND interval '24 months'
                AND animals.sex = 'M'
            ) AS young_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '13 months' AND interval '24 months'
                AND animals.sex = 'F'
            ) AS young_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '25 months' AND interval '36 months' 
                AND animals.sex = 'M'
            ) AS adult_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '25 months' AND interval '36 months' 
                AND animals.sex = 'F'
            ) AS adult_female,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '36 months' 
                AND animals.sex = 'M'
            ) AS old_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '36 months' 
                AND animals.sex = 'F'
            ) AS old_female,
            COUNT(animals.id) FILTER (
                WHERE animals.sex = 'M'
            ) AS total_male,
            COUNT(animals.id) FILTER (
                WHERE animals.sex = 'F'
            ) AS total_female 
        FROM animals
        LEFT JOIN farms ON farms.id = animals.farm_id
    `
	props := repositoriesUtil.GroupByProps[AnimalsByAgeAndFarm]{
		Query:     query,
		TableName: "animals",
		GroupBy:   "farms.name",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetGroupByResults(props)
}

func (r *DashboardRepository) GroupByAge(userId string, filter AnimalsDashboardFilter) (*[]AnimalsByAge, error) {
	query := `
        SELECT 
            age_category,
            COUNT(categorized_animals.id) FILTER (WHERE categorized_animals.sex = 'M') as male,
            COUNT(categorized_animals.id) FILTER (WHERE categorized_animals.sex = 'F') as female
        FROM (
            SELECT animals.*,
            CASE 
                WHEN age(animals.birth_date) < interval '3 months' THEN '0-2 meses'
                WHEN age(animals.birth_date) BETWEEN interval '3 months' AND interval '8 months' THEN '3-8 meses'
                WHEN age(animals.birth_date) BETWEEN interval '9 months' AND interval '12 months' THEN '9-12 meses'
                WHEN age(animals.birth_date) BETWEEN interval '13 months' AND interval '24 months' THEN '13-24 meses'
                WHEN age(animals.birth_date) BETWEEN interval '25 months' AND interval '36 months' THEN '25-36 meses'
                ELSE '+36 meses'
            END AS age_category
            FROM animals
        ) as categorized_animals
    `
	props := repositoriesUtil.GroupByProps[AnimalsByAge]{
		Query:     query,
		TableName: "categorized_animals",
		GroupBy:   "categorized_animals.age_category",
		OrderBy:   "categorized_animals.age_category",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetGroupByResults(props)
}
