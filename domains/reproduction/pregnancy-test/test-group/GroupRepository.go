package testGroup

import (
	"github.com/felipeErnica/rebanho-backend/repositories"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type TestGroupRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TestGroupRepository {
	selectQuery := ` SELECT pregnancy_test_groups.* FROM pregnancy_test_groups`
	return &TestGroupRepository{selectQuery, "pregnancy_test_groups", db}
}

func (r *TestGroupRepository) FindAll() (*[]TestGroup, error) {
	query := r.SelectQuery + " WHERE pregnancy_test_groups.user_id = $1 AND pregnancy_test_groups.deleted_at is null"
	return repositoriesUtil.GetList[TestGroup](r.Db, query, repositories.GetUserId())
}

func (r *TestGroupRepository) FindById(id string) (*TestGroup, error) {
	query := r.SelectQuery + " WHERE pregnancy_test_groups.id = $1 AND pregnancy_test_groups.deleted_at is null"
	return repositoriesUtil.GetOne[TestGroup](r.Db, query, id)
}

func (r *TestGroupRepository) Add(newModel *TestGroup) (*TestGroup, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *TestGroupRepository) Update(model *TestGroup) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *TestGroupRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
