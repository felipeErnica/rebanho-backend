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
	Where     string
	DB        *sqlx.DB
}

type TotalProps[E any] struct {
	Query     string
	TableName string
	UserId    string
	Filter    any
	NumParam  int
	OtherArgs []any
	Where     string
	DB        *sqlx.DB
}

type WhereProps struct {
	TableName string
	UserId    string
	Filter    any
	NumParam  int
	OtherArgs []any
	Where     string
}

/*
Constrói a condição de filtragem das consultas referentes a dashboards
*/
func getDashboardWhereStatement(props WhereProps) (string, []any, error) {

	whereStatement := fmt.Sprintf("where %s.user_id = $%d and deleted_at is null", props.TableName, props.NumParam)
	if props.Where != "" {
		whereStatement = props.Where + fmt.Sprintf(" and %s.user_id = $%d and deleted_at is null", props.TableName, props.NumParam)
	}

	args := props.OtherArgs
	args = append(args, props.UserId)
	if isFiltered(props.Filter) {
		filterArgs := GetFilterArgs(props.Filter)
		filterStatement, _, err := BuildFilterExpressions(props.Filter, props.TableName, props.NumParam+1)
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
	whereStatement, args, err := getDashboardWhereStatement(WhereProps{
		NumParam:  numParam,
		TableName: props.TableName,
		Filter:    props.Filter,
		UserId:    props.UserId,
		OtherArgs: props.OtherArgs,
		Where:     props.Where,
	})

	if err != nil {
		return nil, err
	}

	groupQuery := props.Query + "\n" + whereStatement + "\n" + groupStatement
	if props.OrderBy != "" {
		orderBy := fmt.Sprintf("order by %s", props.OrderBy)
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

	whereStatement, args, err := getDashboardWhereStatement(WhereProps{
		NumParam:  numParam,
		TableName: props.TableName,
		Filter:    props.Filter,
		UserId:    props.UserId,
		OtherArgs: props.OtherArgs,
		Where:     props.Where,
	})
	if err != nil {
		return nil, err
	}

	query := props.Query + "\n" + whereStatement
	return GetOne[E](props.DB, query, args...)
}
