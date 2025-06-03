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
	query := `
        SELECT 
            farms.name AS farm_name,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) < interval '3 months'
                AND animals.sex = 'M'
            ) AS newborn_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '3 months' AND interval '8 months'
                AND animals.sex = 'M'
            ) AS baby_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '9 months' AND interval '12 months'
                AND animals.sex = 'M'
            ) AS child_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '13 months' AND interval '24 months'
                AND animals.sex = 'M'
            ) AS young_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '25 months' AND interval '36 months' 
                AND animals.sex = 'M'
            ) AS adult_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '36 months' 
                AND animals.sex = 'M'
            ) AS old_male,
            COUNT(animals.id) FILTER (
                WHERE animals.sex = 'M'
            ) AS total_male 
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
