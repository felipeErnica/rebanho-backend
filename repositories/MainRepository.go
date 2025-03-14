package repositories

import (
	"database/sql"
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/google/uuid"
)

type Repository[E entity.IEntity] interface {
	setNewId(*E, string)
	buildEntity(*sql.Row) (*E, error)
	buildListEntity(*sql.Rows) (*[]E, error)
	saveOrUpdateScan(string, *E) error
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
    id:= uuid.New().String()
    r.Repository.setNewId(&model, id)
    err:= r.Repository.saveOrUpdateScan(r.InsertQuery, &model)
    return &model, err
}

func (r *RepositoryImpl[E]) Save(model *E) error {
    err:= r.Repository.saveOrUpdateScan(r.UpdateQuery, model)
    return err
}

func (r *RepositoryImpl[E]) Delete(id string) error {
    query:=fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.TableName)
    return execQuery(query, id)
}
