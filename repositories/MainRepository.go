package repositories

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
)

type Repository[E entity.IEntity] interface {
	setNewEntity(model *E)
	buildEntity(row *sql.Row) (model *E, err error)
	buildListEntity(rows *sql.Rows) (arr *[]E, err error)
	saveOrUpdateScan(query string, model *E) error
}

type RepositoryImpl[E entity.IEntity] struct {
    TableName       string
	Repository      Repository[E]
	SelectQueryBody string
	InsertQuery     string
	UpdateQuery     string
}

func (r *RepositoryImpl[E]) FindAll() (*[]E, error) {
    sqlRows, err:=selectQueryList(r.SelectQueryBody)
    list, err:= r.Repository.buildListEntity(sqlRows)
    return list, err
}

func (r *RepositoryImpl[E]) FindById(id string) (*E, error) {
    query:=fmt.Sprintf("%s WHERE %s.id = $1", r.SelectQueryBody, r.TableName)
    sqlRow:=selectQueryOne(query, id)
    entity, err:= r.Repository.buildEntity(sqlRow)
    return entity, err
}

func (r *RepositoryImpl[E]) FindByQuery(query string, args... any) (*E, error) {
    query = r.SelectQueryBody + "\n" + query
    sqlRow:=selectQueryOne(query, args...)
    entity, err:= r.Repository.buildEntity(sqlRow)
    return entity, err
}

func (r *RepositoryImpl[E]) FindListByQuery(query string, args... any) (list *[]E, err error) {
    query = r.SelectQueryBody + "\n" + query
    sqlRow, err:=selectQueryList(query, args...)
    if err != nil {
        return
    }
    entity, err:= r.Repository.buildListEntity(sqlRow)
    return entity, err
}

func (r *RepositoryImpl[E]) Add(model E) (*E, error) {
    r.Repository.setNewEntity(&model)
    err:= r.Repository.saveOrUpdateScan(r.InsertQuery, &model)
    return &model, err
}

func (r *RepositoryImpl[E]) Save(model *E) error {
    err:= r.Repository.saveOrUpdateScan(r.UpdateQuery, model)
    return err
}

func (r *RepositoryImpl[E]) Delete(id string) error {
    timeDeletion:=time.Now()
    query:=fmt.Sprintf("UPDATE %s SET deleted_at = $1 WHERE id = $2", r.TableName)
    return execQuery(query, timeDeletion, id)
}
