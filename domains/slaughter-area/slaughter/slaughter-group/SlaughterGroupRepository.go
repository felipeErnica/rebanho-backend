package slaughterGroup

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type SlaughterGroupRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *SlaughterGroupRepository {
	selectQuery := `
    SELECT slaughter_groups.*, slaughterhouses.name
    FROM slaughter_groups
        LEFT JOIN slaughterhouses ON slaughterhouses.id = slaughter_groups.slaughterhouse_id
    `
	return &SlaughterGroupRepository{selectQuery, "slaughter_groups", db}
}

func (r *SlaughterGroupRepository) FindAll(userId string) (*[]SlaughterGroup, error) {
	query := r.SelectQuery + " WHERE slaughter_groups.user_id = $1 AND slaughter_groups.deleted_at is null"
	return repositoriesUtil.GetList[SlaughterGroup](r.Db, query, userId)
}

func (r *SlaughterGroupRepository) FindById(id string) (*SlaughterGroup, error) {
	query := r.SelectQuery + " WHERE slaughter_groups.id = $1 AND slaughter_groups.deleted_at is null"
	return repositoriesUtil.GetOne[SlaughterGroup](r.Db, query, id)
}

func (r *SlaughterGroupRepository) Add(newModel *SlaughterGroupSave) (*SlaughterGroupSave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *SlaughterGroupRepository) Update(model *SlaughterGroupSave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *SlaughterGroupRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
