package repositoriesUtil

import (
	"reflect"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PageProps struct {
	Sort   string
	Order  string
	Cursor string
	Filter any
	UserId string
}

type PageBuilderProps struct {
	QueryBody       string
	CountQuery      string
	Limit           *int
	TableName       string
	DbConn          *sqlx.DB
	SortExpressions []SortExpression
	PageProps
}

type PageSelectionProps[E any] struct {
	Dest        *[]E
	Filter      any
	Query       string
	IsFirstPage bool
	FirstParam  any
	CreatedAt   *time.Time
	Id          uuid.UUID
	DbConn      *sqlx.DB
	UserId      string
}

/*Envia as inforações para obter a lista da página*/
func getPageList[E any](props PageSelectionProps[E]) error {
	query := strings.Join(strings.Fields(props.Query), " ")
	db := props.DbConn
	filterArgs := GetFilterArgs(props.Filter)
	args := []any{props.UserId}
	util.LogInfo(query, true)

	//Constrói a lista de paramêtros da solicitação, levando em conta
	//a paginação e os critério de filtro.
	if len(filterArgs) != 0 {
		args = append(args, filterArgs...)
	}

	if props.FirstParam != nil {
		args = append(args, props.FirstParam)
	}

	if props.CreatedAt != nil {
		args = append(args, props.CreatedAt)
	}

	if props.Id != uuid.Nil {
		args = append(args, props.Id)
	}

	err := db.Select(props.Dest, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func getSliceArgs(slice []string, args []any) []any {
	for _, arg := range slice {
		args = append(args, arg)
	}
	return args
}

/*Constrói e organiza os valores do filtro*/
func GetFilterArgs(filter any) []any {
	values := reflect.ValueOf(filter)
	fields := reflect.TypeOf(filter)
	if values.Kind() == reflect.Pointer {
		values = values.Elem()
		fields = fields.Elem()
	}

	args := []any{}
	for i := 1; i < values.NumField(); i++ {
		fieldValue := values.Field(i)
		fieldType := fields.Field(i)
		if !fieldValue.IsNil() {
			value := fieldValue.Elem().Interface()
			if fieldValue.Elem().Type().String() == "string" && !strings.HasSuffix(fieldType.Name, "Id") {
				value = "%" + value.(string) + "%"
			}
			if fieldValue.Elem().Kind() == reflect.Slice {
				slice := fieldValue.Elem().Interface().([]string)
				args = getSliceArgs(slice, args)
			} else {
				args = append(args, value)
			}
		}
	}
	return args
}

func getTotalCount(query string, userId string, filter any, db *sqlx.DB) (int, error) {
	query = strings.Join(strings.Fields(query), " ")
	args := []any{userId}
	filterArgs := GetFilterArgs(filter)
	args = append(args, filterArgs...)
    util.LogInfo(query, true)
	result := db.QueryRow(query, args...)
    if result.Err() != nil {
        return 0, result.Err()
    }

    var total int
    result.Scan(&total)
	return total, nil
}

/*
Constrói uma página
Params:
sort - Coluna de ordenação
order - Direção do ordenamento (crescente e descrescente)
filter - Critérios de filtro (valores e campo)
*/
func BuildPage[E any](props PageBuilderProps) (page *entity.Page[E], err error) {

	firstParam, createdAt, id, err := DecodeCursor(props.Cursor)
	if err != nil {
		return
	}

	const PAGE_LIMIT = 200

	limit := PAGE_LIMIT
	if props.Limit != nil {
		limit = *props.Limit
	}

	pageProps := PageQueryProps{
		Filter:          props.Filter,
		CountQuery:      props.CountQuery,
		QueryBody:       props.QueryBody,
		Sort:            props.Sort,
		Order:           props.Order,
		Limit:           limit,
		Cursor:          props.Cursor,
		TableName:       props.TableName,
		SortExpressions: props.SortExpressions,
	}

	query, err := buildPageQuery(pageProps)
	if err != nil {
		return page, err
	}

	countQuery, err := buildCountQuery(pageProps)
	if err != nil {
		return page, err
	}
	total, err := getTotalCount(countQuery, props.UserId, props.Filter, props.DbConn)
	if err != nil {
		return page, err
	}

	//Parâmetros para obter as informações da página
	list := []E{}
	selectionProps := PageSelectionProps[E]{
		Dest:        &list,
		Query:       query,
		IsFirstPage: props.Cursor == "",
		FirstParam:  firstParam,
		CreatedAt:   createdAt,
		Id:          id,
		Filter:      props.Filter,
		DbConn:      props.DbConn,
		UserId:      props.UserId,
	}

	err = getPageList(selectionProps)
	if err != nil {
		return
	}

	nextCursor, err := CreateCursorKey(props.Sort, list)
	if err != nil {
		return
	}

	//Criação da página
	page = &entity.Page[E]{
		List:        &list,
		Total:       total,
		HasNextPage: len(list) >= limit,
		NextCursor:  nextCursor,
	}

	return page, err
}
