package repositoriesUtil

import (
	"fmt"
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
	Filter      any
	Query       string
	IsFirstPage bool
	FirstParam  any
	SecondParam *time.Time
	DbConn      *sqlx.DB
	UserId      string
}

/*Envia as inforações para obter a lista da página*/
func getPageList[E any](props SelectionProps[E]) error {
	query := strings.Join(strings.Fields(props.Query), " ")
	db := props.DbConn
	filterArgs := getFilterArgs(props.Filter)
	args := []any{props.UserId}
	fmt.Println(query)

	/*
	   Constrói a lista de paramêtros da solicitação, levando em conta
	   a paginação e os critério de filtro
	*/
	if len(filterArgs) != 0 {
		args = append(args, filterArgs...)
	}

	if props.FirstParam != nil {
		args = append(args, props.FirstParam)
	}

	if props.SecondParam != nil {
		args = append(args, props.SecondParam)
	}

	err := db.Select(props.Dest, query, args...)
	if err != nil {
		return err
	}

	return nil
}

/*Constrói e organiza os valores do filtro*/
func getFilterArgs(filter any) []any {
	values := reflect.ValueOf(filter)
	if values.Kind() == reflect.Pointer {
		values = values.Elem()
	}

	args := []any{}
	for i := 1; i < values.NumField(); i++ {
		field := values.Field(i)
		if !field.IsNil() {

			value := field.Elem().Interface()
			if field.Elem().Type().String() == "string" {
				value = "%" + value.(string) + "%"
				fmt.Println(value)
			}

			args = append(args, value)
		}
	}
	return args
}

/*
Constrói uma página
Params:
sort - Coluna de ordenação
order - Direção do ordenamento (crescente e descrescente)
filter - Critérios de filtro (valores e campo)
*/
func BuildPage[E any](props PageProps) (page *entity.Page[E], err error) {

	firstParam, secondParam, err := decodeCursor(props.Cursor)
	if err != nil {
		return
	}

	pageProps := PageQueryProps{
		Filter:     props.Filter,
		QueryBody:  props.QueryBody,
		Sort:       props.Sort,
		Order:      props.Order,
		Limit:      props.Limit,
		IsNull:     firstParam == nil,
		Cursor:     props.Cursor,
		NullFields: props.NullFields,
		TableName:  props.TableName,
	}

	query, err := buildPageQuery(pageProps)
	if err != nil {
		return page, err
	}

	//Parâmetros para obter as informações da página
	list := []E{}
	selectionProps := SelectionProps[E]{
		Dest:        &list,
		Query:       query,
		IsFirstPage: props.Cursor == "",
		FirstParam:  firstParam,
		SecondParam: secondParam,
		Filter:      props.Filter,
		DbConn:      props.DbConn,
		UserId:      props.UserId,
	}

	err = getPageList(selectionProps)
	if err != nil {
		return
	}

	nextCursor, err := createCursorKey(props.Sort, list)
	if err != nil {
		return
	}

	//Criação da página
	page = &entity.Page[E]{
		List:        &list,
		HasNextPage: len(list) >= props.Limit,
		NextCursor:  nextCursor,
	}

	return page, err
}
