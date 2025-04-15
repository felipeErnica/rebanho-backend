package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/google/uuid"
)

type Repository[E entity.IEntity] interface {
	setNewEntity(model *E, id string, createdAt time.Time)
	buildEntity(row *sql.Row) (model *E, err error)
	buildListEntity(rows *sql.Rows) (arr *[]E, err error)
	saveOrUpdateScan(query string, model *E) error
}

type RepositoryImpl[E entity.IEntity] struct {
	TableName       string
	Repository      Repository[E]
	SelectQueryBody util.SelectConstructor
	InsertQuery     util.InsertQuery
	UpdateQuery     util.UpdateQuery
}

func (r *RepositoryImpl[E]) FindAll() (arr *[]E, err error) {
	query := r.SelectQueryBody.Where(fmt.Sprintf("%s.user_id = $1", r.TableName)).Build()
	sqlRows, err := selectQueryList(query)
	if err != nil {
		return
	}
	list, err := r.Repository.buildListEntity(sqlRows)
	sqlRows.Close()
	return list, err
}

func (r *RepositoryImpl[E]) FindById(id string) (*E, error) {
	query := r.SelectQueryBody.Where(r.TableName + ".id = $1")
	sqlRow := selectQueryOne(query.Build(), id)
	entity, err := r.Repository.buildEntity(sqlRow)
	return entity, err
}

func (r *RepositoryImpl[E]) FindByQuery(query *util.SelectConstructor, args ...any) (*E, error) {
	sqlRow := selectQueryOne(query.Build(), args...)
	entity, err := r.Repository.buildEntity(sqlRow)
	return entity, err
}

func (r *RepositoryImpl[E]) FindListByQuery(query *util.SelectConstructor, args ...any) (list *[]E, err error) {
	sqlRow, err := selectQueryList(query.Build(), args...)
	if err != nil {
		return
	}
	entity, err := r.Repository.buildListEntity(sqlRow)
	sqlRow.Close()
	return entity, err
}

func (r *RepositoryImpl[E]) Add(model *E) (*E, error) {
	id := uuid.NewString()
	createdAt := time.Now()
	r.Repository.setNewEntity(model, id, createdAt)
	err := r.Repository.saveOrUpdateScan(r.InsertQuery.Build(), model)
	return model, err
}

func (r *RepositoryImpl[E]) Save(model *E) error {
	err := r.Repository.saveOrUpdateScan(r.UpdateQuery.Build(), model)
	return err
}

func (r *RepositoryImpl[E]) Delete(id string) error {
	timeDeletion := time.Now()
	query := util.NewDeleteQuery(r.TableName)
	return execQuery(query, timeDeletion, id)
}
