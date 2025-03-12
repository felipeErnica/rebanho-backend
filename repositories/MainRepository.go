package repositories

import (
	"strings"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type Repository interface {
    getFields(string) (string, string)
    buildPage(query string, sort string, args... any) (entity.PageInterface, error)
}

type BaseRepository struct {
    Repository      Repository
    SelectPageBody  string
    DeletedAtField  string
    DateFields      []string
}

func (r *BaseRepository) firstPage(sort string, direction string) (page entity.PageInterface, err error) {
	firstField, secondField := r.Repository.getFields(sort)
	query := new(util.QueryBuilder).GetFirstPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort)
}

func (r *BaseRepository) nextPage(sort string, direction string, cursor string) (page entity.PageInterface, err error) {
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

func (r *BaseRepository) nextPageDate(sort string, direction string, cursor string) (page entity.PageInterface, err error) {
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

func (r *BaseRepository) firstPageDelete(sort string, direction string) (page entity.PageInterface, err error) {
	firstField, secondField := r.Repository.getFields(sort)
	query := new(util.QueryBuilder).GetFirstPage(r.SelectPageBody, firstField, secondField, r.DeletedAtField, direction)
	return r.Repository.buildPage(query, sort)
}

func (r *BaseRepository) nextPageDelete(sort string, direction string, cursor string) (page entity.PageInterface, err error) {
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

func (r *BaseRepository) isDateField(sort string) bool {
    for i:=0; i<len(r.DateFields); i++ {
        if strings.EqualFold(sort, r.DateFields[i]) {
            return true
        }
    }
    return false
}

func (r *BaseRepository) nextPageDateDelete(sort string, direction string, cursor string) (page entity.PageInterface, err error) {
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

func (r *BaseRepository) GetPage(sort string, direction string, cursor string) (page entity.PageInterface, err error) {

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

func (r *BaseRepository) GetPageDelete(sort string, direction string, cursor string) (page entity.PageInterface, err error) {

    if sort == "" {
        sort = "created_at"
    }

    if direction == "" {
        direction = "asc"
    }

	if cursor == "" {
		return r.firstPageDelete(sort, direction)
	}

    if r.isDateField(sort) {
        return r.nextPageDateDelete(sort, direction, cursor)
    }

    return r.nextPageDelete(sort, direction, cursor)
}
