package repositories

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PageRepository[E entity.IEntity] interface {
    getFields(sort string) (firstField string, secondField string)
    createKey(sort string, lastEntry *E) (key string)
    Repository[E]
}

type PageRepositoryImpl[E entity.IEntity] struct {
    Base            *RepositoryImpl[E]
    PageRepository  PageRepository[E]
    DateFields      []string
}

func (p *PageRepositoryImpl[E]) getNullOrdering(direction string) string {
	nullDirection := " DESC NULLS LAST"
	if direction == "asc" {
		nullDirection = " ASC NULLS FIRST"
	}
	return nullDirection
}

func (q *PageRepositoryImpl[E]) getSignal(direction string) string {
	signal := ">"
	if direction == "desc" {
		signal = "<"
	}
	return signal
}

func (q *PageRepositoryImpl[E]) verifyQueryParams(sort string, direction string) (string, string) {
    if sort == "" {
        sort = "created_at"
    }

    if direction == "" {
        direction = "asc"
    }
    return sort, direction
}

func (r *PageRepositoryImpl[E]) buildPage(query string, sort string, args... any) (page *entity.Page[E], err error) {
    sqlStatement, err:= selectQueryList(query, args...)
    if err != nil {
        return
    }

    arr, err:=r.PageRepository.buildListEntity(sqlStatement)    
    if err != nil {
        return
    }

    nextCursor, err:= r.CreateNextCursor(sort, *arr)
    if err !=  nil {
        return
    }
    
    page = &entity.Page[E]{
        HasNextPage: len(*arr) == PAGE_LIMIT,
        NextCursor: nextCursor,
        List: arr,
    }

    return page, err
}

func (r *PageRepositoryImpl[E]) nextPage(cursor string, query *util.QueryConstructor, 
    sort string, direction string, firstField string, secondField string) (page *entity.Page[E], err error) {
    firstParam, secondParam, err:= decodeCursor(cursor) 
    if err != nil {
        return 
    }

    signal:= r.getSignal(direction)
    if firstParam == "null" {
        query.Where(firstField + " IS NULL").And(fmt.Sprintf("%s %s $1", secondField, signal))
        return r.buildPage(query.Build(), sort, secondParam)
    } else {
        query.Where(fmt.Sprintf("(%s,%s) %s ($1, $2)", firstField, secondField, signal))
        return r.buildPage(query.Build(), sort, firstParam, secondParam)
    }
}

func (r *PageRepositoryImpl[E]) nextPageDate(cursor string, query *util.QueryConstructor, 
    sort string, direction string, firstField string, secondField string) (page *entity.Page[E], err error) {
    firstParam, secondParam, err:= decodeCursorTime(cursor) 
    if err != nil {
        return 
    }

    signal:= r.getSignal(direction)
    if firstParam == nil {
        query.Where(firstField + " IS NULL").And(fmt.Sprintf("%s %s $1", secondField, signal))
        return r.buildPage(query.Build(), sort, secondParam)
    } else {
        query.Where(fmt.Sprintf("(%s,%s) %s ($1, $2)", firstField, secondField, signal))
        return r.buildPage(query.Build(), sort, firstParam, secondParam)
    }
}

func (r *PageRepositoryImpl[E]) nextRandomPage(cursor string, query *util.QueryConstructor, 
    sort string, direction string, firstField string, secondField string, args... any) (page *entity.Page[E], err error) {
    firstParam, secondParam, err:= decodeCursor(cursor) 
    if err != nil {
        return 
    }

    signal:= r.getSignal(direction)
    paramNumber:= len(args) + 1
    if firstParam == "null" {
        query.And(firstField + " IS NULL").And(fmt.Sprintf("%s %s $%d", secondField, signal, paramNumber))
        args = append(args, secondField)
        return r.buildPage(query.Build(), sort, args...)
    } else {
        query.And(fmt.Sprintf("(%s,%s) %s ($%d, $%d)", firstField, secondField, signal, paramNumber, paramNumber + 1))
        args = append(args, firstParam, secondParam)
        return r.buildPage(query.Build(), sort, args...)
    }
}

func (r *PageRepositoryImpl[E]) nextRandomPageDate(cursor string, query *util.QueryConstructor, 
    sort string, direction string, firstField string, secondField string, args... any) (page *entity.Page[E], err error){
    firstParam, secondParam, err:= decodeCursorTime(cursor) 
    if err != nil {
        return 
    }

    signal:= r.getSignal(direction)
    paramNumber:= len(args) + 1
    if firstParam == nil {
        query.And(firstField + " IS NULL").And(fmt.Sprintf("%s %s $%d", secondField, signal, paramNumber))
        args = append(args, secondParam)
        return r.buildPage(query.Build(), sort, args...)
    } else {
        query.And(fmt.Sprintf("(%s,%s) %s ($%d, $%d)", firstField, secondField, signal, paramNumber, paramNumber + 1))
        args = append(args, firstParam, secondParam)
        return r.buildPage(query.Build(), sort, args...)
    }
}

func (r *PageRepositoryImpl[E]) isDateField(sort string) bool {
    for i:=range r.DateFields {
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
    cursor = base64.StdEncoding.EncodeToString([]byte(key))
    return cursor, err
}

func (r *PageRepositoryImpl[E]) FindPage(sort string, direction string, cursor string) (page *entity.Page[E], err error) {
    sort, direction = r.verifyQueryParams(sort, direction)
    firstField, secondField:= r.PageRepository.getFields(sort)
    nullOrdering:= r.getNullOrdering(direction)

    query:= new(util.QueryConstructor).FromQuery(r.Base.SelectQueryBody)
    query.Order(firstField + nullOrdering).AndOrder(secondField + " " + direction)
    query.Limit(PAGE_LIMIT)

	if cursor == "" {
		return r.buildPage(query.Build(), sort)
    }

    if r.isDateField(sort) {
        return r.nextPageDate(cursor, query, sort, direction, firstField, secondField)
    }
		
    return r.nextPage(cursor, query, sort, direction, firstField, secondField)
}

func (r *PageRepositoryImpl[E]) FindRandomQueryPage(query *util.QueryConstructor, sort string, 
    direction string, cursor string, args... any) (page *entity.Page[E], err error) {
    sort, direction = r.verifyQueryParams(sort, direction)
    firstField, secondField:= r.PageRepository.getFields(sort)
    nullOrderinrg:= r.getNullOrdering(direction)
    query.Order(firstField + nullOrderinrg).AndOrder(secondField + nullOrderinrg)
    query.Limit(PAGE_LIMIT)

	if cursor == "" {
		return r.buildPage(query.Build(), sort, args...)
    }

    if r.isDateField(sort) {
        return r.nextRandomPageDate(cursor, query, sort, direction, firstField, secondField, args...)
    }
	return r.nextRandomPage(cursor, query, sort, direction, firstField, secondField, args...)
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

func (r *PageRepositoryImpl[E]) FindListByQuery(query *util.QueryConstructor, args... any) (list *[]E, err error) {
    return r.Base.FindListByQuery(query.Build(), args...)
}

func (r *PageRepositoryImpl[E]) Add(model E) (*E, error) {
    return r.Base.Add(model)
}

func (r *PageRepositoryImpl[E]) Save(model *E) error {
    return r.Base.Save(model)
}

func (r *PageRepositoryImpl[E]) SoftDelete(id string) error {
    return r.Base.SoftDelete(id)
}

func (r *PageRepositoryImpl[E]) HardDelete(id string) error {
    return r.Base.HardDelete(id)
}
