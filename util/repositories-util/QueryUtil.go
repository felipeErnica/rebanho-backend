package repositoriesUtil

import (
	"fmt"
)

type SortExpression struct {
	Column     string
	Expression string
}

func NewSort(column string, statement string) *SortExpression {
	return &SortExpression{
		Column:     column,
		Expression: statement,
	}
}

type PageQueryProps struct {
	CountQuery      string
	QueryBody       string
	Sort            string
	Order           string
	Limit           int
	Cursor          string
	Filter          any
	TableName       string
	SortExpressions []SortExpression
}

func buildSortStatement(props PageQueryProps) string {
	for _, condition := range props.SortExpressions {
		if condition.Column == props.Sort {
			return condition.Expression
		}
	}
	return fmt.Sprintf("%s.%s", props.TableName, props.Sort)
}

func buildOrderStatement(props PageQueryProps) string {
	paramField := buildSortStatement(props) + " " + props.Order
	createdAtField := fmt.Sprintf("%s.created_at", props.TableName)
	idField := fmt.Sprintf("%s.id", props.TableName)
	sortStatement := fmt.Sprintf("order by %s, %s, %s", paramField, createdAtField, idField)
	return sortStatement
}

func buildWhereStatement(props PageQueryProps) (whereStatement string, err error) {
	signal := "<"
	if props.Order == "asc" {
		signal = ">"
	}
	numParam := 1

	whereStatement = fmt.Sprintf("where %[1]s.deleted_at is null and %[1]s.user_id = $%[2]d", props.TableName, numParam)
	numParam++
	if isFiltered(props.Filter) {
		filterStatements, filterParam, err := BuildFilterExpressions(props.Filter, props.TableName, numParam)
		if err != nil {
			return whereStatement, err
		}
		numParam = filterParam
		whereStatement = whereStatement + filterStatements
	}

	if props.Cursor == "" {
		return whereStatement, err
	}

	paramField := buildSortStatement(props)
	createdAtField := fmt.Sprintf("%s.created_at", props.TableName)
	idField := fmt.Sprintf("%s.id", props.TableName)
	fieldsExpression := fmt.Sprintf("%s, %s, %s", paramField, createdAtField, idField)
	paramsExpression := fmt.Sprintf("$%d, $%d, $%d", numParam, numParam+1, numParam+2)
	paginationCriteria := fmt.Sprintf("(%s) %s (%s)", fieldsExpression, signal, paramsExpression)
	whereStatement = whereStatement + " and " + paginationCriteria
	return whereStatement, err
}

func buildPageQuery(props PageQueryProps) (query string, err error) {
    whereStatement, err := buildWhereStatement(props)
	if err != nil {
		return query, err
	}
	sortStatement := buildOrderStatement(props)
	limitStatement := fmt.Sprintf("limit %d", props.Limit)
	query = props.QueryBody + "\n" + whereStatement + "\n" + sortStatement + "\n" + limitStatement
	return query, err
}

func buildCountQuery(props PageQueryProps) (string, error) {
    whereStatement := fmt.Sprintf("where %[1]s.deleted_at is null and %[1]s.user_id = $1", props.TableName)
	filterStatement, _, err := BuildFilterExpressions(props.Filter, props.TableName, 2)
	if err != nil {
		return "", err
	}
    whereStatement = whereStatement + filterStatement
    query := props.CountQuery + "\n" + whereStatement
    return query, err
}
