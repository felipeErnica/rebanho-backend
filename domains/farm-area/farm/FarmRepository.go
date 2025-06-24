package farm

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type FarmRepository struct {
	Query     string
	TableName string
	DB        *sqlx.DB
}

func NewRepository(db *sqlx.DB) *FarmRepository {
	query := "SELECT farms.* FROM farms"
	return &FarmRepository{query, "farms", db}
}

func (r *FarmRepository) SearchFarm(userId string, input string) (*[]entity.SearchEntity, error) {
	query := `
        select id, name as label 
        from farms 
        where user_id = $1 and name ilike $2 and deleted_at is null
        order by label
        `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId, input)
}

func (r *FarmRepository) Add(farm *Farm) (*Farm, error) {
	return repositoriesUtil.Add(r.DB, r.TableName, farm)
}
