package repositoriesUtil

import (
	"fmt"
	"strings"
	"slices"
)

type PageQueryProps struct {
	QueryBody        string
	Sort             string
	Order            string
	Args             []any
	Limit            int
	IsNull           bool
	NullFields       []string
	Cursor           string
	FilterStatements string
	TableName        string
}

func buildSortStatement(props PageQueryProps) string {
	orderField := fmt.Sprintf("%s.%s", props.TableName, props.Sort) + " " + strings.ToUpper(props.Order)
	sortStatement := fmt.Sprintf("ORDER BY %s, %s.created_at DESC", orderField, props.TableName)

	if isNullable(props.Sort, props.NullFields) {
		nullStatement := "NULLS FIRST"
		if props.Order == "asc" {
			nullStatement = "NULLS LAST"
		}
		orderField = orderField + " " + nullStatement
		sortStatement = fmt.Sprintf("ORDER BY %s, %s.created_at DESC", orderField, props.TableName)
	}
	return sortStatement
}

func  isNullable(sort string, nullFields []string) bool {
	return slices.Contains(nullFields, sort)
}

func buildWhereStatement(props PageQueryProps) string {
	signal := "<"
	if props.Order == "asc" {
		signal = ">"
	}
	commonCriteria := fmt.Sprintf("WHERE %[1]s.deleted_at is null AND %[1]s.user_id = $1", props.TableName)

	if props.Cursor == "" {
		return commonCriteria
	}

	orderCriteria := fmt.Sprintf("(%[1]s.%[2]s, %[1]s.created_at) %[3]s ($2, $3)", props.TableName, props.Sort, signal)
	if props.IsNull {
		orderCriteria = fmt.Sprintf("%[1]s.%[2]s is null OR (%[1]s.%[2]s is not null and %[1]s.created_at %[3]% $1)",
			props.TableName, props.Sort, signal)
	}

	whereStatement := commonCriteria + " AND " + orderCriteria
	if props.FilterStatements != "" {
		whereStatement = whereStatement + " AND " + props.FilterStatements
	}
	return whereStatement
}

func BuildPageQuery(props PageQueryProps) string {
	whereStatement := buildWhereStatement(props)
	sortStatement := buildSortStatement(props)
	limitStatement := fmt.Sprintf("LIMIT %d", props.Limit)
	query := props.QueryBody + "\n" + whereStatement + "\n" + sortStatement + "\n" + limitStatement
	return query
}
