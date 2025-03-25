package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/google/uuid"
)

type Repository[E entity.IEntity] interface {
	SetNewEntity(model *E, id string, createdAt time.Time)
	BuildEntity(row *sql.Row) (model *E, err error)
	BuildListEntity(rows *sql.Rows) (arr *[]E, err error)
	SaveOrUpdateScan(query string, model *E) error
    Delete(id string) error
}

type RepositoryImpl[E entity.IEntity] struct {
    TableName       string
	Repository      Repository[E]
	SelectQueryBody string
	InsertQuery     string
	UpdateQuery     string
}

func (r *RepositoryImpl[E]) FindAll() (arr *[]E, err error) {
    sqlRows, err:=SelectQueryList(r.SelectQueryBody)
    if err != nil {
        return
    }
    list, err:= r.Repository.BuildListEntity(sqlRows)
    sqlRows.Close()
    return list, err
}

func (r *RepositoryImpl[E]) FindById(id string) (*E, error) {
    query:=new(util.QueryConstructor).FromQuery(r.SelectQueryBody).Where(r.TableName + ".id = $1").Build()
    sqlRow:=SelectQueryOne(query, id)
    entity, err:= r.Repository.BuildEntity(sqlRow)
    return entity, err
}

func (r *RepositoryImpl[E]) FindByQuery(query string, args... any) (*E, error) {
    query = r.SelectQueryBody + "\n" + query
    sqlRow:=SelectQueryOne(query, args...)
    entity, err:= r.Repository.BuildEntity(sqlRow)
    return entity, err
}

func (r *RepositoryImpl[E]) FindListByQuery(query string, args... any) (list *[]E, err error) {
    query = r.SelectQueryBody + "\n" + query
    sqlRow, err:=SelectQueryList(query, args...)
    if err != nil {
        return
    }
    entity, err:= r.Repository.BuildListEntity(sqlRow)
    sqlRow.Close()
    return entity, err
}

func (r *RepositoryImpl[E]) Add(model E) (*E, error) {
    id:=uuid.NewString()
    createdAt:=time.Now()
    r.Repository.SetNewEntity(&model, id, createdAt)
    err:= r.Repository.SaveOrUpdateScan(r.InsertQuery, &model)
    return &model, err
}

func (r *RepositoryImpl[E]) Save(model *E) error {
    err:= r.Repository.SaveOrUpdateScan(r.UpdateQuery, model)
    return err
}

func (r *RepositoryImpl[E]) SoftDelete(id string) error {
    timeDeletion:=time.Now()
    query:=new(util.QueryConstructor).SoftDelete(r.TableName).Build()
    return ExecQuery(query, timeDeletion, id)
}

func (r *RepositoryImpl[E]) HardDelete(id string) error {
    query:=new(util.QueryConstructor).Delete(r.TableName).Build()
    return ExecQuery(query, id)
}

func (r *RepositoryImpl[E]) Delete(id string) error {
    return r.Repository.Delete(id)
}
