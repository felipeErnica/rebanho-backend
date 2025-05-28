package repositoriesUtil

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/jmoiron/sqlx"
)

type PageProps struct {
	QueryBody  string
	Sort       string
	Order      string
	Limit      int
	NullFields []string
	Cursor     string
	Filter     any
	TableName  string
	DbConn     *sqlx.DB
	UserId     string
}

type SelectionProps[E any] struct {
	Dest        *[]E
	Query       string
	IsFiltered  bool
	IsFirstPage bool
	FirstParam  any
	SecondParam time.Time
	FilterArgs  []any
	DbConn      *sqlx.DB
	UserId      string
}

func BuildPage[E any](props PageProps) (page *entity.Page[E], err error) {

	firstParam, secondParam, err := DecodeCursor(props.Cursor)
	if err != nil {
		return
	}

	pageProps := PageQueryProps{
		QueryBody:  props.QueryBody,
		Sort:       props.Sort,
		Order:      props.Order,
		Limit:      props.Limit,
		IsNull:     firstParam == nil,
		Cursor:     props.Cursor,
		NullFields: props.NullFields,
		TableName:  props.TableName,
	}

	query := BuildPageQuery(pageProps)
	list := []E{}

	filterValue := reflect.ValueOf(props.Filter)
	if filterValue.Kind() == reflect.Pointer {
		filterValue = filterValue.Elem()
	}
	isFiltered := filterValue.FieldByName("IsFiltered").Bool()

	selectionProps := SelectionProps[E]{
		Dest:        &list,
		Query:       query,
		IsFirstPage: props.Cursor == "",
		IsFiltered:  isFiltered,
		FirstParam:  firstParam,
		SecondParam: secondParam,
		DbConn:      props.DbConn,
		UserId:      props.UserId,
	}

	err = SelectPage(selectionProps)
	if err != nil {
		return
	}

	nextCursor, err := CreateCursorKey(props.Sort, list)
	if err != nil {
		return
	}

	page = &entity.Page[E]{
		List:        &list,
		HasNextPage: len(list) >= props.Limit,
		NextCursor:  nextCursor,
	}

	return page, err
}

func SelectPage[E any](props SelectionProps[E]) error {
	query := strings.Join(strings.Fields(props.Query), " ")
	db := props.DbConn

	if props.IsFiltered && props.IsFirstPage {
		err := db.Select(props.Dest, query, props.UserId, props.FilterArgs)
		return err
	}

	if !props.IsFiltered && props.IsFirstPage {
		err := db.Select(props.Dest, query, props.UserId)
		return err
	}

	if !props.IsFiltered && !props.IsFirstPage {
		err := db.Select(props.Dest, query, props.UserId, props.FirstParam, props.SecondParam)
		return err
	}

	if props.IsFiltered && !props.IsFirstPage {
		err := db.Select(props.Dest, query, props.UserId, props.FirstParam, props.SecondParam, props.FilterArgs)
		return err
	}

	return errors.New("Formato Inválido")
}
