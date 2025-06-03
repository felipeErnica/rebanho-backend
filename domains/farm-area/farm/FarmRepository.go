package farm

import (
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

func (r *FarmRepository) Add(farm *Farm) (*Farm, error) {
    return repositoriesUtil.Add(r.DB, r.TableName, farm)
}
