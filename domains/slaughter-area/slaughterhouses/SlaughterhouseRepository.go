package slaughterhouses

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type SlaughterhouseRepository struct {
	SelectQuery string
    TableName string
    Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *SlaughterhouseRepository {
    selectQuery := "SELECT slaughterhouses.* FROM slaughterhouses"
    return &SlaughterhouseRepository{selectQuery, "slaughterhouses", db}
}

func (r *SlaughterhouseRepository) FindAll() (*[]Slaughterhouse, error) {
    query := r.SelectQuery + " WHERE user_id = $1 AND deleted_at is null"
	return repositoriesUtil.GetList[Slaughterhouse](r.Db, query)
}

func (r *SlaughterhouseRepository) FindById(id string) (*Slaughterhouse, error) {
    query := r.SelectQuery + " WHERE id = $1 AND deleted_at is null"
	return repositoriesUtil.GetOne[Slaughterhouse](r.Db, query)
}

func (r *SlaughterhouseRepository) Add(newModel *Slaughterhouse) (*Slaughterhouse, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *SlaughterhouseRepository) Update(model *Slaughterhouse) error {
	return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *SlaughterhouseRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
