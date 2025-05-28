package repositoriesUtil

import (
	"fmt"
	"slices"
	"strings"
)

type PageQueryProps struct {
	QueryBody  string
	Sort       string
	Order      string
	Limit      int
	IsNull     bool
	NullFields []string
	Cursor     string
	Filter     any
	TableName  string
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

func isNullable(sort string, nullFields []string) bool {
	return slices.Contains(nullFields, sort)
}

func buildWhereStatement(props PageQueryProps) (whereStatement string, err error) {
	signal := "<"
	if props.Order == "asc" {
		signal = ">"
	}
	numParam := 1

	whereStatement = fmt.Sprintf("WHERE %[1]s.deleted_at is null AND %[1]s.user_id = $%[2]d", props.TableName, numParam)
	numParam++
	if isFiltered(props.Filter) {
		filterStatements, filterParam, err := BuildFilterStatements(props.Filter, props.TableName, numParam)
		if err != nil {
			return whereStatement, err
		}
		numParam = filterParam
		whereStatement = whereStatement + filterStatements
	}

	if props.Cursor == "" {
		return whereStatement, err
	}

	paginationCriteria := fmt.Sprintf("(%[1]s.%[2]s, %[1]s.created_at) %[3]s ($%[4]d, $%[5]d)",
		props.TableName, props.Sort, signal, numParam, numParam + 1)
	if props.IsNull {
		paginationCriteria = fmt.Sprintf("%[1]s.%[2]s is null OR (%[1]s.%[2]s is not null and %[1]s.created_at %[3]% $%[4]d)",
			props.TableName, props.Sort, signal, numParam)
	}

	whereStatement = whereStatement + " AND " + paginationCriteria
	return whereStatement, err
}

func buildPageQuery(props PageQueryProps) (query string, err error) {
	whereStatement, err := buildWhereStatement(props)
	if err != nil {
		return query, err
	}

	sortStatement := buildSortStatement(props)
	limitStatement := fmt.Sprintf("LIMIT %d", props.Limit)
	query = props.QueryBody + "\n" + whereStatement + "\n" + sortStatement + "\n" + limitStatement
	return query, err
}
