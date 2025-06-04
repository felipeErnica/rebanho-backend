package repositoriesUtil

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type GroupByProps struct {
	Query     string
	TableName string
	GroupBy   string
	UserId    string
	Filter    any
	DB        *sqlx.DB
}

type TotalProps struct {
	Query     string
	TableName string
	UserId    string
	Filter    any
	DB        *sqlx.DB
}

/*
Constrói a condição de filtragem das consultas referentes a dashboards
*/
func getDashboardWhereStatement(filter any, tablename string, userId string) (string, []any, error) {
	whereStatement := fmt.Sprintf("WHERE %s.user_id = $1", tablename)
	args := []any{userId}
	if isFiltered(filter) {
		filterArgs := getFilterArgs(filter)
		filterStatement, _, err := buildFilterStatements(filter, tablename, 2)
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
func GetGroupByResults[E any](props GroupByProps) (*[]E, error) {
	groupStatement := fmt.Sprintf("GROUP BY %s", props.GroupBy)
	whereStatement, args, err := getDashboardWhereStatement(props.Filter, props.TableName, props.UserId)
	if err != nil {
		return nil, err
	}

	groupQuery := props.Query + "\n" + whereStatement + "\n" + groupStatement
	return GetList[E](props.DB, groupQuery, args...)
}

/*
Constrói uma consulta SQL que retorna informações quantitativas sobre os campos e retorna os resultados.
*/
func GetTotalResults[E any](props TotalProps) (*E, error) {
	whereStatement, args, err := getDashboardWhereStatement(props.Filter, props.TableName, props.UserId)
	if err != nil {
		return nil, err
	}

	query := props.Query + "\n" + whereStatement
	return GetOne[E](props.DB, query, args...)
}
