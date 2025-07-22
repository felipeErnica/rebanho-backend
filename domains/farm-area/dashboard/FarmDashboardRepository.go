package dashboard

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
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
        with pasture_animal as (
            select 
                pastures.id,
                pastures.farm_id,
                count(animals.id) as animals_number
            from pastures
                left join animals on animals.pasture_id = pastures.id
            group by pastures.id, pastures.farm_id
        )
        select
            farms.id as farm_id,
            farms.name as farm_name,
            count(pastures.id) as pastures_number,
            sum(pastures.animals_number) as animals_number
        from farms
            left join pasture_animal as pastures on pastures.farm_id = farms.id
        where farms.user_id = $1 and farms.deleted_at is null
        group by farms.id, farms.name
        order by farms.name
    `
    return repositoriesUtil.GetList[FarmInfo](r.DB, query, userId)
}

func (r *FarmDashboardRepository) GetPastureInfo(userId string, farmId string) (*[]PastureInfo, error) {
	query := `
        select
            pastures.id as pasture_id,
            pastures.name as pasture_name,
            bull.id as bull_id,
            concat_ws(' - ', bull.ring_number, bull.name) as bull_name,
            count(animals.id) as animals_number
        from pastures
            left join animals on animals.pasture_id = pastures.id
            left join animals as bull on bull.id = pastures.bull_id
        where pastures.user_id = $1 and pastures.farm_id = $2 and pastures.deleted_at is null
        group by pastures.name, pastures.id, bull.id, bull.name
        order by pastures.name
    `
    return repositoriesUtil.GetList[PastureInfo](r.DB, query, userId, farmId)
}
