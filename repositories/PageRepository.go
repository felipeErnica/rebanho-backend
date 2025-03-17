package repositories

import (
	"encoding/base64"
	"errors"
	"strings"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PageRepository[E entity.IEntity] interface {
    getFields(sort string) (firstField string, secondField string)
    createKey(sort string, lastEntry *E) (key string)
    buildPage(query string, sort string, args... any) (page *entity.Page[E], err error)
    Repository[E]
}

type PageRepositoryImpl[E entity.IEntity] struct {
    Base            *RepositoryImpl[E]
    PageRepository  PageRepository[E]
    DateFields      []string
}

func (r *PageRepositoryImpl[E]) firstPage(sort string, direction string) (page *entity.Page[E], err error) {
	firstField, secondField := r.PageRepository.getFields(sort)
	query := new(util.QueryBuilder).GetFirstPage(r.Base.SelectQueryBody, firstField, secondField, direction)
	return r.PageRepository.buildPage(query, sort)
}

func (r *PageRepositoryImpl[E]) nextPage(sort string, direction string, cursor string) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursor(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.PageRepository.getFields(sort)
	if firstParam != "null" {
		query := new(util.QueryBuilder).GetNextPage(r.Base.SelectQueryBody, firstField, secondField, direction)
		return r.PageRepository.buildPage(query, sort, firstParam, secondParam)
	}
	query := new(util.QueryBuilder).GetNextPageNull(r.Base.SelectQueryBody, firstField, secondField, direction)
	return r.PageRepository.buildPage(query, sort, secondParam)
}

func (r *PageRepositoryImpl[E]) nextPageDate(sort string, direction string, cursor string) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursorTime(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.PageRepository.getFields(sort)
	if firstParam != nil {
		query := new(util.QueryBuilder).GetNextPage(r.Base.SelectQueryBody, firstField, secondField, direction)
		return r.PageRepository.buildPage(query, sort, firstParam, secondParam)
	}

	query := new(util.QueryBuilder).GetNextPageNull(r.Base.SelectQueryBody, firstField, secondField, direction)
	return r.PageRepository.buildPage(query, sort, secondParam)
}

func (r *PageRepositoryImpl[E]) firstPageConditional(sort string, direction string, 
    conditionalField string, conditionalValue any) (page *entity.Page[E], err error) {
	firstField, secondField := r.PageRepository.getFields(sort)
	query := new(util.QueryBuilder).GetFirstPageConditional(r.Base.SelectQueryBody, conditionalField, 
        firstField, secondField, direction)
	return r.PageRepository.buildPage(query, sort, conditionalValue)
}

func (r *PageRepositoryImpl[E]) nextPageConditional(sort string, direction string, cursor string, 
    conditionalField string, conditionalValue any) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursor(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.PageRepository.getFields(sort)
	if firstParam != "null" {
		query := new(util.QueryBuilder).GetNextPageConditional(r.Base.SelectQueryBody, conditionalField, 
            firstField, secondField, direction)
		return r.PageRepository.buildPage(query, sort, conditionalValue, firstParam, secondParam)
	}
	query := new(util.QueryBuilder).GetNextPageNullConditional(r.Base.SelectQueryBody, conditionalField,
        firstField, secondField, direction)
	return r.PageRepository.buildPage(query, sort, conditionalValue, secondParam)
}

func (r *PageRepositoryImpl[E]) nextPageDateConditional(sort string, direction string, cursor string, 
    conditionalField string, conditionalValue any) (page *entity.Page[E], err error) {
	firstParam, secondParam, err := decodeCursorTime(cursor)
	if err != nil {
		return
	}

	firstField, secondField := r.PageRepository.getFields(sort)
	if firstParam != nil {
		query := new(util.QueryBuilder).GetNextPageConditional(r.Base.SelectQueryBody, conditionalField, 
            firstField, secondField, direction)
		return r.PageRepository.buildPage(query, sort, conditionalValue, firstParam, secondParam)
	}

	query := new(util.QueryBuilder).GetNextPageNullConditional(r.Base.SelectQueryBody, conditionalField, 
        firstField, secondField, direction)
	return r.PageRepository.buildPage(query, sort, conditionalValue, secondParam)
}

func (r *PageRepositoryImpl[E]) isDateField(sort string) bool {
    for i:=0; i<len(r.DateFields); i++ {
        if strings.EqualFold(sort, r.DateFields[i]) {
            return true
        }
    }
    return false
}

func (r *PageRepositoryImpl[E]) CreateNextCursor(sort string, array []E) (cursor string, err error) {
    if len(array) == 0 {
        err = errors.New("A matriz está vazia!")
        return
    }
    lastEntry:=array[len(array) - 1]
    key:= r.PageRepository.createKey(sort, &lastEntry)
    cursor = base64.RawStdEncoding.EncodeToString([]byte(key))
    return cursor, err
}

func (r *PageRepositoryImpl[E]) FindPage(sort string, direction string, cursor string) (page *entity.Page[E], err error) {

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

func (r *PageRepositoryImpl[E]) FindPageCondional(sort string, direction string, cursor string, 
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

func (r *PageRepositoryImpl[E]) FindAll() (*[]E, error) {
    return r.Base.FindAll()
}

func (r *PageRepositoryImpl[E]) FindById(id string) (*E, error) {
    return r.Base.FindById(id)
}

func (r *PageRepositoryImpl[E]) FindByQuery(query string, args... any) (*E, error) {
    return r.Base.FindByQuery(query, args...)
}

func (r *PageRepositoryImpl[E]) FindListByQuery(query string, args... any) (list *[]E, err error) {
    return r.Base.FindListByQuery(query, args...)
}

func (r *PageRepositoryImpl[E]) Add(model E) (*E, error) {
    return r.Base.Add(model)
}

func (r *PageRepositoryImpl[E]) Save(model *E) error {
    return r.Base.Save(model)
}

func (r *PageRepositoryImpl[E]) Delete(id string) error {
    return r.Base.Delete(id)
}
