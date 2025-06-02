package animalDashboard

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type DashboardRepository struct {
	Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *DashboardRepository {
	return &DashboardRepository{db}
}

func (r *DashboardRepository) GroupByAgeAndFarm(userId string) (*AnimalsByAgeAndFarm, error) {
	query := `
        SELECT 
            farms.name AS farm_name,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) <= interval '2 months'
            ) AS newborn_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '3 months' AND interval '8 months'
            ) AS baby_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '9 months' AND interval '12 months'
            ) AS child_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '13 months' AND interval'24 months'
            ) AS young_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) BETWEEN interval '25 months' AND interval '36 months'
            ) AS adult_male,
            COUNT(animals.id) FILTER (
                WHERE age(animals.birth_date) > interval '36 months'
            ) AS old_male 
        FROM animals
        LEFT JOIN farms ON farms.id = animals.farm_id
        WHERE user_id = $1
        GROUP BY farms.name;
    `
    return repositoriesUtil.GetOne[AnimalsByAgeAndFarm](r.Db, query, userId)
}
