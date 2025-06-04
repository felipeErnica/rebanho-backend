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
	props := repositoriesUtil.TotalProps{
		Query:     query,
		TableName: "animals",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetTotalResults[TotalBySex](props)
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
	props := repositoriesUtil.GroupByProps{
		Query:     query,
		TableName: "animals",
		GroupBy:   "farms.name",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
	return repositoriesUtil.GetGroupByResults[AnimalsByAgeAndFarm](props)
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
	props := repositoriesUtil.GroupByProps{
		Query:     query,
		TableName: "categorized_animals",
		GroupBy:   "categorized_animals.age_category",
		UserId:    userId,
		Filter:    filter,
		DB:        r.DB,
	}
    return repositoriesUtil.GetGroupByResults[AnimalsByAge](props)
}
