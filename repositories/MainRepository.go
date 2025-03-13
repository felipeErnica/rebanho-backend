package repositories

import (
	"database/sql"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/google/uuid"
)

type Repository[E entity.IEntity] interface {
    setNewId(model *E, id string)
    getFields(string) (string, string)
    buildPage(query string, sort string, args... any) (*entity.Page[E], error)
    buildEntity(sqlRow *sql.Row) (*E, error)
    buildListEntity(sqlRows *sql.Rows) (*[]E, error)
    saveOrUpdateScan(query string, base *E) error
}

type BaseRepository[E entity.IEntity] struct {
    DeleteQuery     string
    Repository      Repository[E]
    SimpleQueryBody string
    SelectPageBody  string
    DeletedAtField  string
    InsertQuery     string
    UpdateQuery     string
    DateFields      []string
}

func (r *BaseRepository[E]) firstPage(sort string, direction string) (page *entity.Page[E], err error) {
	firstField, secondField := r.Repository.getFields(sort)
	query := new(util.QueryBuilder).GetFirstPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort)
}

func (r *BaseRepository[E]) nextPage(sort string, direction string, cursor string) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursor(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.Repository.getFields(sort)
	if firstParam != "null" {
		query := new(util.QueryBuilder).GetNextPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
		return r.Repository.buildPage(query, sort, firstParam, secondParam)
	}
	query := new(util.QueryBuilder).GetNextPageNull(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort, secondParam)
}

func (r *BaseRepository[E]) nextPageDate(sort string, direction string, cursor string) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursorTime(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.Repository.getFields(sort)
	if firstParam != nil {
		query := new(util.QueryBuilder).GetNextPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
		return r.Repository.buildPage(query, sort, firstParam, secondParam)
	}

	query := new(util.QueryBuilder).GetNextPageNull(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort, secondParam)
}

func (r *BaseRepository[E]) firstPageConditional(sort string, direction string, 
    conditionalField string, conditionalValue any) (page *entity.Page[E], err error) {
	firstField, secondField := r.Repository.getFields(sort)
	query := new(util.QueryBuilder).GetFirstPageConditional(r.SelectPageBody, conditionalField, 
        firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort, conditionalValue)
}

func (r *BaseRepository[E]) nextPageConditional(sort string, direction string, cursor string, 
    conditionalField string, conditionalValue any) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursor(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.Repository.getFields(sort)
	if firstParam != "null" {
		query := new(util.QueryBuilder).GetNextPageConditional(r.SelectPageBody, conditionalField, firstField, secondField, 
            r.DeletedAtField, direction)
		return r.Repository.buildPage(query, sort, conditionalValue, firstParam, secondParam)
	}
	query := new(util.QueryBuilder).GetNextPageNullConditional(r.SelectPageBody, conditionalField,
        firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort, conditionalValue, secondParam)
}

func (r *BaseRepository[E]) nextPageDateConditional(sort string, direction string, cursor string, 
    conditionalField string, conditionalValue any) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursorTime(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.Repository.getFields(sort)
	if firstParam != nil {
		query := new(util.QueryBuilder).GetNextPageConditional(r.SelectPageBody, conditionalField, firstField, secondField, 
            r.DeletedAtField, direction)
		return r.Repository.buildPage(query, sort, conditionalValue, firstParam, secondParam)
	}

	query := new(util.QueryBuilder).GetNextPageNullConditional(r.SelectPageBody, conditionalField, firstField, secondField,
        r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort, conditionalValue, secondParam)
}

func (r *BaseRepository[E]) firstPageDelete(sort string, direction string) (page *entity.Page[E], err error) {
	firstField, secondField := r.Repository.getFields(sort)
	query := new(util.QueryBuilder).GetDeletedFirstPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort)
}

func (r *BaseRepository[E]) nextPageDelete(sort string, direction string, cursor string) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursor(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.Repository.getFields(sort)
	if firstParam != "null" {
		query := new(util.QueryBuilder).GetDeletedNextPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
		return r.Repository.buildPage(query, sort, firstParam, secondParam)
	}
	query := new(util.QueryBuilder).GetDeletedNextPageNull(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort, secondParam)
}

func (r *BaseRepository[E]) isDateField(sort string) bool {
    for i:=0; i<len(r.DateFields); i++ {
        if strings.EqualFold(sort, r.DateFields[i]) {
            return true
        }
    }
    return false
}

func (r *BaseRepository[E]) nextPageDateDelete(sort string, direction string, cursor string) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursorTime(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.Repository.getFields(sort)
	if firstParam != nil {
		query := new(util.QueryBuilder).GetDeletedNextPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
		return r.Repository.buildPage(query, sort, firstParam, secondParam)
	}

	query := new(util.QueryBuilder).GetDeletedNextPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort, secondParam)
}

func (r *BaseRepository[E]) FindPage(sort string, direction string, cursor string) (page *entity.Page[E], err error) {

    if sort == "" {
        sort = "created_at"
    }

    if direction == "" {
        direction = "asc"
    }

	if cursor == "" {
		return r.firstPage(sort, direction)
    }

    if r.isDateField(sort) {
        return r.nextPageDate(sort, direction, cursor)
    }
		
    return r.nextPage(sort, direction, cursor)
}

func (r *BaseRepository[E]) FindPageCondional(sort string, direction string, cursor string, 
    conditionalField string, conditionalValue any) (page *entity.Page[E], err error) {

    if sort == "" {
        sort = "created_at"
    }

    if direction == "" {
        direction = "asc"
    }

	if cursor == "" {
		return r.firstPageConditional(sort, direction, conditionalField, conditionalValue)
    }

    if r.isDateField(sort) {
        return r.nextPageDateConditional(sort, direction, cursor, conditionalField, conditionalValue)
    }
		
    return r.nextPageConditional(sort, direction, cursor, conditionalField, conditionalValue)
}

func (r *BaseRepository[E]) FindDeletedPage(sort string, direction string, cursor string) (page *entity.Page[E], err error) {

    if sort == "" {
        sort = "deleted_at"
    }

    if direction == "" {
        direction = "desc"
    }

	if cursor == "" {
		return r.firstPageDelete(sort, direction)
	}

    if r.isDateField(sort) {
        return r.nextPageDateDelete(sort, direction, cursor)
    }

    return r.nextPageDelete(sort, direction, cursor)
}

func (r *BaseRepository[E]) FindById(id string) (*E, error) {
    query:=r.SimpleQueryBody + "\n" + "WHERE id = $1"
    sqlRow:=selectQueryOne(query, id)
    entity, err:= r.Repository.buildEntity(sqlRow)
    return entity, err
}

func (r *BaseRepository[E]) FindByQuery(query string, args... any) (*E, error) {
    query = r.SimpleQueryBody + "\n" + query
    sqlRow:=selectQueryOne(query, args...)
    entity, err:= r.Repository.buildEntity(sqlRow)
    return entity, err
}

func (r *BaseRepository[E]) FindListByQuery(query string, args... any) (list *[]E, err error) {
    query = r.SimpleQueryBody + "\n" + query
    sqlRow, err:=selectQueryList(query, args...)
    if err != nil {
        return
    }
    entity, err:= r.Repository.buildListEntity(sqlRow)
    return entity, err
}

func (r *BaseRepository[E]) Add(model E) (*E, error) {
    id:= uuid.New().String()
    r.Repository.setNewId(&model, id)
    err:= r.Repository.saveOrUpdateScan(r.InsertQuery, &model)
    return &model, err
}

func (r *BaseRepository[E]) Save(model *E) error {
    err:= r.Repository.saveOrUpdateScan(r.UpdateQuery, model)
    return err
}

func (r *BaseRepository[E]) Delete(id string) error {
    deleteTime:=time.Now()
    return execQuery(r.DeleteQuery, deleteTime, id)
}
