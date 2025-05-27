package util

import (
	"fmt"
	"strings"
)

type PageQueryProps struct {
	QueryBody        string
	Sort             string
	Order            string
	Args             []any
	Limit            int
	IsNull           bool
	Cursor           string
	FilterStatements string
}

func buildSortStatement(sort string, order string, isNullable bool) string {
	orderField := order + " " + strings.ToUpper(sort)
	sortStatement := fmt.Sprintf("ORDER BY %s, created_at DESC", orderField)
	if isNullable {
		nullStatement := "NULLS FIRST"
		if order == "asc" {
			nullStatement = "NULLS LAST"
		}
		orderField := order + " " + strings.ToUpper(sort) + " " + nullStatement
		sortStatement = fmt.Sprintf("ORDER BY %s, created_at DESC", orderField)
	}
	return sortStatement
}

func buildWhereStatement(sort string, order string, isNull bool, filterStatements string) string {
	signal := ">"
	if sort == "asc" {
		signal = "<"
	}
	orderCriteria := fmt.Sprintf("(%s, created_at) %s ($1, $2) and deleted_at is null", order, signal)
	if isNull {
		orderCriteria = fmt.Sprintf("%[1]s is null OR (%[1]s is not null and created_at %[2]% $1)", order, signal)
	}
	whereStatement := "WHERE " + filterStatements + " AND " + orderCriteria
	return whereStatement
}

func BuildPageQuery(props PageQueryProps) string {
	whereStatement := buildWhereStatement(props.Sort, props.Order, props.IsNull, props.FilterStatements)
	sortStatement := buildSortStatement(props.Sort, props.Order, props.IsNull)
	limitStatement := fmt.Sprintf("LIMIT %d", props.Limit)
	query := props.QueryBody + "\n" + whereStatement + "\n" + sortStatement + "\n" + limitStatement
	if props.Cursor == "" {
		query = props.QueryBody + "\n" + sortStatement + "\n" + limitStatement
	}
	return query
}
