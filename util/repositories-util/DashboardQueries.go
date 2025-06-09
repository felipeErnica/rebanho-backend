package repositoriesUtil

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type GroupByProps[E any] struct {
	Query     string
	TableName string
	GroupBy   string
	OrderBy   string
	UserId    string
	Filter    any
	NumParam  int
    OtherArgs []any
	DB        *sqlx.DB
}

type TotalProps[E any] struct {
	Query     string
	TableName string
	UserId    string
	Filter    any
	NumParam  int
    OtherArgs []any
	DB        *sqlx.DB
}

/*
Constrói a condição de filtragem das consultas referentes a dashboards
*/
func getDashboardWhereStatement(filter any, tablename string, userId string, numParam int, otherArgs ...any) (string, []any, error) {
	whereStatement := fmt.Sprintf("WHERE %s.user_id = $%d", tablename, numParam)
	args := otherArgs
    args = append(args, userId)
	if isFiltered(filter) {
		filterArgs := getFilterArgs(filter)
		filterStatement, _, err := buildFilterStatements(filter, tablename, numParam + 1)
		if err != nil {
			return "", []any{}, err
		}
		whereStatement = whereStatement + filterStatement
		args = append(args, filterArgs...)
	}
	return whereStatement, args, nil
}

/*
Constrói uma consulta SQL de GROUP BY e retorna os resultados.
*/
func GetGroupByResults[E any](props GroupByProps[E]) (*[]E, error) {

	numParam := 1
	if props.NumParam != 0 {
		numParam = props.NumParam
	}

	groupStatement := fmt.Sprintf("GROUP BY %s", props.GroupBy)
	whereStatement, args, err := getDashboardWhereStatement(props.Filter, props.TableName, props.UserId, numParam, props.OtherArgs...)
	if err != nil {
		return nil, err
	}

	groupQuery := props.Query + "\n" + whereStatement + "\n" + groupStatement
	if props.OrderBy != "" {
		orderBy := fmt.Sprintf("ORDER BY %s", props.OrderBy)
		groupQuery = groupQuery + "\n" + orderBy
	}

	return GetList[E](props.DB, groupQuery, args...)
}

/*
Constrói uma consulta SQL que retorna informações quantitativas sobre os campos e retorna os resultados.
*/
func GetTotalResults[E any](props TotalProps[E]) (*E, error) {

	numParam := 1
	if props.NumParam != 0 {
		numParam = props.NumParam
	}

	whereStatement, args, err := getDashboardWhereStatement(props.Filter, props.TableName, props.UserId, numParam)
	if err != nil {
		return nil, err
	}

	query := props.Query + "\n" + whereStatement
	return GetOne[E](props.DB, query, args...)
}
