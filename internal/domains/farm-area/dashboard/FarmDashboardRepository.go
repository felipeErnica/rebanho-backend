package dashboard

import (
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type FarmDashboardRepository struct {
	TableName string
	DB        *sqlx.DB
}

func NewRepository(db *sqlx.DB) *FarmDashboardRepository {
	return &FarmDashboardRepository{"pastures", db}
}

func (r *FarmDashboardRepository) GetFarmInfo(userId string) (*[]FarmInfo, error) {
	query := `
        WITH pasture_animal AS (
            SELECT 
                pastures.id,
                pastures.farm_id,
                COUNT(animals.id) AS animals_number
            FROM pastures
                LEFT JOIN animals ON animals.pasture_id = pastures.id
            GROUP BY pastures.id, pastures.farm_id
        )
        SELECT
            farms.id AS farm_id,
            farms.name AS farm_name,
            COUNT(pastures.id) AS pastures_number,
            SUM(pastures.animals_number) AS animals_number
        FROM farms
            LEFT JOIN pasture_animal AS pastures ON pastures.farm_id = farms.id
        WHERE farms.user_id = $1 AND farms.deleted_at IS NULL
        GROUP BY farms.id, farms.name
        ORDER BY farms.name
    `
	return util.GetList[FarmInfo](r.DB, query, userId)
}

func (r *FarmDashboardRepository) GetPastureInfo(userId string, farmId string) (*[]PastureInfo, error) {
	query := `
        SELECT
            pastures.id AS pasture_id,
            pastures.name AS pasture_name,
            bull.id AS bull_id,
            CONCAT_WS(' - ', bull.ring_number, bull.name) AS bull_name,
            COUNT(animals.id) AS animals_number
        FROM pastures
            LEFT JOIN animals ON animals.pasture_id = pastures.id
            LEFT JOIN animals AS bull ON bull.id = pastures.bull_id
        WHERE pastures.user_id = $1 AND pastures.farm_id = $2 AND pastures.deleted_at IS NULL
        GROUP BY pastures.name, pastures.id, bull.id, bull.name
        ORDER BY pastures.name
    `
	return util.GetList[PastureInfo](r.DB, query, userId, farmId)
}
