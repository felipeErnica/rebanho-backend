package repositoriesUtil

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

func buildFilterArrays(array []string, field string, numParam int) (string, int) {
	params := fmt.Sprintf("$%d", numParam)
    numParam++
	for i := 1; i < len(array); i++ {
		params += fmt.Sprintf(", $%d", numParam)
		numParam++
	}
	statement := fmt.Sprintf("%s in (%s)", field, params)
	return statement, numParam
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
		statement, newNumParam = buildFilterArrays(t, sqlField, numParam)
	case bool:
		statement, newNumParam = buildFilterBool(sqlField, numParam)
	default:
		errMsg := fmt.Sprintf("Tipo Inválido: %s", structField)
		err = errors.New(errMsg)
		return
	}
	return statement, newNumParam, err
}

func BuildFilterExpressions(filter any, mainTable string, numParam int) (filterStatement string, newNumParam int, err error) {
	var buffer bytes.Buffer

	filterTypes := reflect.TypeOf(filter)
	filterValues := reflect.ValueOf(filter)
	if filterValues.Kind() == reflect.Pointer {
		filterValues = filterValues.Elem()
		filterTypes = filterTypes.Elem()
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
			statement, param, errFilter := buildFilterStatement(fieldValue, sqlField, structField, numParam)
			numParam = param
			buffer.WriteString(" AND " + statement)
			if errFilter != nil {
				err = errFilter
				return
			}
		}
	}

	return buffer.String(), numParam, err
}

func isFiltered(filter any) bool {
	filterValue := reflect.ValueOf(filter)
	if filterValue.Kind() == reflect.Pointer {
		filterValue = filterValue.Elem()
	}
	return filterValue.FieldByName("IsFiltered").Bool()
}
