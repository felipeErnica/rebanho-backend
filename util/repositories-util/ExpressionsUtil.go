package repositoriesUtil

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type SortField struct {
	Field string
	Order string
}

// Através de um mapa da campo com expressões SQL, retorna a expressão apropriada ao campo
func GetSortExpression(
	sortMap map[string]SortField,
	sort string,
	order string,
) (string, error) {
	sortFields := strings.Split(sort, ",")
	expressions := []string{}
	for i := range sortFields {
		key := strings.TrimSpace(sortFields[i])
		expression, ok := sortMap[key]
		if !ok {
			err := fmt.Errorf("A expressão de ordenamento (%s) solicitada não existe!", key)
			return "", err
		}

		ordering := expression.Order
		if i == 0 {
			ordering = order
		}

		expressions = append(expressions, fmt.Sprintf("%s %s", expression.Field, ordering))
	}

	sortExpression := strings.Join(expressions, ", ")
	return sortExpression, nil
}

/*
Através de um mapa de campos relacionados a expressões SQL,
retorna a expressão de ordenamento do cursor relacionada ao campo.
*/
func GetCursorExpression(
	sortMap map[string]SortField,
	sort string,
	order string,
	cursor string,
	numParam int,
) (string, int, error) {

	if cursor == "" {
		return "", numParam, nil
	}

	sortFields := strings.Split(sort, ",")
	conditions := []string{}

	for i := range sortFields {

		parts := []string{}

		key := strings.TrimSpace(sortFields[i])
		expression, ok := sortMap[key]
		if !ok {
			err := fmt.Errorf("A expressão de ordenamento solicitada (%s) não existe!", key)
			return "", numParam, err
		}

		ordering := expression.Order
		if i == 0 {
			ordering = order
		}

		signal := ">"
		if ordering == "desc" {
			signal = "<"
		}

		for j := range i {
			partKey := strings.TrimSpace(sortFields[j])
			partExpression := sortMap[partKey]
			parts = append(parts, fmt.Sprintf("%s = $%d", partExpression.Field, j+numParam))
		}

		parts = append(parts, fmt.Sprintf("%s %s $%d", expression.Field, signal, numParam+i))
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(parts, " and ")))
	}

	cursorExpression := strings.Join(conditions, " or ")
	nextNumParam := len(sortFields) + numParam

	return "(" + cursorExpression + ")", nextNumParam, nil
}

func AddCommonFields(sort string) string {
	arr := []string{sort, "id", "created_at"}
	return strings.Join(arr, ",")
}

func AddNewFields(sort string, fields ...string) string {
	arr := []string{sort}
	arr = append(arr, fields...)
	return strings.Join(arr, ",")
}

func GetWhereExpression(expressions ...string) string {
	validExp := []string{}
	for _, exp := range(expressions) {
		if exp != "" {
			validExp = append(validExp, exp)
		}
	}

	if len(validExp) == 0 {
		return ""
	}

	concatExp := strings.Join(validExp, " and ")
	whereExpression := " where " + concatExp
	return whereExpression
}

func GetFilterExpressions(filter any, mainTable string, numParam int) (string, int, error) {
	var buffer bytes.Buffer

	filterTypes := reflect.TypeOf(filter)
	filterValues := reflect.ValueOf(filter)
	if filterValues.Kind() == reflect.Pointer {
		filterValues = filterValues.Elem()
		filterTypes = filterTypes.Elem()
	}

	if !isFiltered(filterValues) {
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

	if len(array) == 0 {
		return "", numParam
	}

	params := ""
	for range array {
		params += fmt.Sprintf("$%d, ", numParam)
		numParam++
	}
	params = strings.TrimSuffix(params, ", ")
	statement := fmt.Sprintf("%s in (%s)", field, params)
	return statement, numParam
}
