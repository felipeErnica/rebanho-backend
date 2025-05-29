package inseminationGroup

import (
	"github.com/felipeErnica/rebanho-backend/repositories"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *GroupRepository {
	selectQuery := `
        SELECT group.*, bull.name as bull_name
        FROM insemination_groups as group
            LEFT JOIN animals as bull ON bull.id = group.bull_id
    `
	return &GroupRepository{selectQuery, "group", db}
}

func (r *GroupRepository) FindAll() (*[]Group, error) {
    query := r.SelectQuery + " WHERE group.user_id = $1 AND group.deleted_at is null ORDER BY group.insemination_date"
	return repositoriesUtil.GetList[Group](r.Db, query, repositories.GetUserId())
}

func (r *GroupRepository) FindById(id string) (*Group, error) {
    query := r.SelectQuery + " WHERE group.id = $1 AND group.deleted_at is null"
	return repositoriesUtil.GetOne[Group](r.Db, query, id)
}

func (r *GroupRepository) Add(newModel *GroupSave) (*GroupSave, error) {
    return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *GroupRepository) Update(model *GroupSave) error {
    return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *GroupRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
