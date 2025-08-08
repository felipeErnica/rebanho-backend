package repositoriesUtil

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

// Através de um mapa da campo com expressões SQL, retorna a expressão apropriada ao campo
func GetSortExpression(expressionMap map[string]string, sort string, order string) (string, error) {
	sortFields := strings.Split(sort, ",")
	sortExpression := ""
	for _, field := range sortFields {
		expression, ok := expressionMap[field]
		if !ok {
			err := errors.New("A expressão de ordenamento solicitada não existe!")
			return "", err
		}
		sortExpression = sortExpression + expression + " " + order + ", "
	}
	sortExpression = strings.TrimSuffix(sortExpression, ", ")
	return sortExpression, nil
}

/*
Através de um mapa de campos relacionados a expressões SQL, retorna a expressão de ordenamento do
cursor relacionada ao campo
*/
func GetCursorExpression(
	sortMap map[string]string,
	sort string,
	order string,
	tableName string,
	cursorArgs []any,
	numParam int,
) (string, int, error) {

	if len(cursorArgs) == 0 {
		return "", numParam, nil
	}

	sortFields := strings.Split(sort, ",")

	sortExpression := ""
	for _, field := range sortFields {
		expression, ok := sortMap[field]
		if !ok {
			err := errors.New("A expressão de cursor para ordenamento solicitada não existe!")
			return "", 0, err
		}
		sortExpression = sortExpression + expression + ", "
	}
	sortExpression = strings.TrimSuffix(sortExpression, ", ")

	commonFields := fmt.Sprintf("%[1]s.created_at, %[1]s.id", tableName)
	fieldsExpression := sortExpression + ", " + commonFields

	signal := ">"
	if order == "desc" {
		signal = "<"
	}

	paramExpression := ""
	nextNumParam := numParam
	for i := range cursorArgs {
		nextNumParam := numParam + i
		paramExpression = paramExpression + fmt.Sprintf("$%d, ", nextNumParam)
	}
	paramExpression = strings.TrimSuffix(paramExpression, ", ")

	cursorExpression := fmt.Sprintf("(%s) %s (%s)", fieldsExpression, signal, paramExpression)
	return cursorExpression, nextNumParam + 1, nil
}

func GetFilterExpressions(filter any, mainTable string, numParam int) (string, int, error) {
	var buffer bytes.Buffer

	filterTypes := reflect.TypeOf(filter)
	filterValues := reflect.ValueOf(filter)
	if filterValues.Kind() == reflect.Pointer {
		filterValues = filterValues.Elem()
		filterTypes = filterTypes.Elem()
	}

    if (!isFiltered(filterValues)) {
        return "", numParam, nil
    }

	for i := 1; i < filterTypes.NumField(); i++ {
		field := filterTypes.Field(i)
		structField := field.Name

		tableName, ok := field.Tag.Lookup("table")
		if !ok {
			tableName = mainTable
		}

		sqlField := fmt.Sprintf("%s.%s", tableName, field.Tag.Get("db"))
		value := filterValues.Field(i)
		if !value.IsNil() {
			fieldValue := value.Elem().Interface()
			statement, param, err := buildFilterStatement(fieldValue, sqlField, structField, numParam)
			numParam = param
			buffer.WriteString(statement + " and ")
			if err != nil {
				return "", 0, err
			}
		}
	}
    filterExpression := buffer.String()
    filterExpression = strings.TrimSuffix(filterExpression, " and ")
	return filterExpression, numParam, nil
}

func isFiltered(filter reflect.Value) bool {
	return filter.FieldByName("IsFiltered").Bool()
}

func buildFilterStatement(value any, sqlField string, structField string, numParam int) (statement string, newNumParam int, err error) {
	switch t := value.(type) {
	case string:
		statement, newNumParam = buildFilterString(sqlField, numParam)
	case time.Time:
		if strings.HasPrefix(structField, "Max") {
			statement, newNumParam = buildFilterMaxNumberAndDate(sqlField, numParam)
			break
		}
		statement, newNumParam = buildFilterMinNumberAndDate(sqlField, numParam)
	case int:
		if strings.HasPrefix(structField, "Max") {
			statement, newNumParam = buildFilterMaxNumberAndDate(sqlField, numParam)
			break
		}
		statement, newNumParam = buildFilterMinNumberAndDate(sqlField, numParam)
	case float64:
		if strings.HasPrefix(structField, "Max") {
			statement, newNumParam = buildFilterMaxNumberAndDate(sqlField, numParam)
			break
		}
		statement, newNumParam = buildFilterMinNumberAndDate(sqlField, numParam)
	case []string:
		statement, newNumParam = GetSliceExpressions(t, sqlField, numParam)
	case bool:
		statement, newNumParam = buildFilterBool(sqlField, numParam)
	default:
		errMsg := fmt.Sprintf("Tipo Inválido: %s", structField)
		err = errors.New(errMsg)
		return
	}
	return statement, newNumParam, err
}

func buildFilterMinNumberAndDate(field string, numParam int) (string, int) {
	return fmt.Sprintf("%s >= $%d", field, numParam), numParam + 1
}

func buildFilterMaxNumberAndDate(field string, numParam int) (string, int) {
	return fmt.Sprintf("%s <= $%d", field, numParam), numParam + 1
}

func buildFilterString(field string, numParam int) (string, int) {
	statement := fmt.Sprintf("%s ilike $%d", field, numParam)
	if strings.HasSuffix(field, "id") {
		statement = fmt.Sprintf("%s = $%d", field, numParam)
	}
	return statement, numParam + 1
}

func buildFilterBool(field string, numParam int) (string, int) {
	return fmt.Sprintf("%s = $%d", field, numParam), numParam + 1
}

func GetSliceExpressions(array []string, field string, numParam int) (string, int) {
    params := ""
	for range array {
		params += fmt.Sprintf("$%d, ", numParam)
		numParam++
	}
    params = strings.TrimSuffix(params, ", ")
	statement := fmt.Sprintf("%s in (%s)", field, params)
	return statement, numParam
}
