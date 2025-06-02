package weightGroup

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type WeightGroupRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *WeightGroupRepository {
	selectQuery := "SELECT weight_groups.* FROM weight_groups"
	return &WeightGroupRepository{selectQuery, "weight_groups", db}
}

func (r *WeightGroupRepository) FindAll(userId string) (*[]WeightGroup, error) {
	query := r.SelectQuery + " WHERE user_id = $1 AND deleted_at is null"
	return repositoriesUtil.GetList[WeightGroup](r.Db, query)
}

func (r *WeightGroupRepository) FindById(id string) (*WeightGroup, error) {
	query := r.SelectQuery + " WHERE id = $1 AND deleted_at is null"
	return repositoriesUtil.GetOne[WeightGroup](r.Db, query, id)
}

func (r *WeightGroupRepository) Add(newModel *WeightGroup) (*WeightGroup, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *WeightGroupRepository) Update(model *WeightGroup) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *WeightGroupRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
