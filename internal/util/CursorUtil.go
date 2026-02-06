package util

import (
	"encoding/base64"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

/*
Cria um novo cursor com a informação do último objeto da lista.
Os parâmetros são selecionados com base no parâmetro de ordenamento.
*/
func CreateCursorKey[E any](sort string, list []E) string {
	listSize := len(list)
	if listSize == 0 {
		return ""
	}

	sortFields := strings.Split(sort, ",")
	lastEntry := list[listSize-1]
	entryType := reflect.TypeOf(lastEntry)
	values := reflect.ValueOf(lastEntry)
	cursorArgs := []string{}

	for _, sortField := range sortFields {
		sortField = strings.TrimSpace(sortField)
		arg, err := getValueFromSortField(sortField, entryType, values)
		if err != nil {
			return ""
		}
		cursorArgs = append(cursorArgs, arg)
	}

	data := ""
	for _, arg := range cursorArgs {
		data = data + arg + ","
	}
	data = strings.TrimSuffix(data, ",")

	cursor := base64.StdEncoding.EncodeToString([]byte(data))
	return cursor
}

func getValueFromSortField(sortField string, entryType reflect.Type, values reflect.Value) (string, error) {
	var value any
	isFound := false
	for i := range entryType.NumField() {
		field := entryType.Field(i)
		tag := field.Tag.Get("db")
		if sortField == tag {
			isFound = true
			value = values.Field(i).Interface()
		}
	}

	if !isFound {
		err := fmt.Errorf("O campo não existe: %s", sortField)
		return "", err
	}
	return getParamValue(value), nil
}

/*
Converte o valor do parâmetros para texto, se o valor for nulo,
retorna "null", se for data, insere o prefixo {date} para tratamento
apropiado posteriormente
*/
func getParamValue(value any) string {
	switch t := value.(type) {
	case *string:
		param := ""
		if t != nil {
			param = *t
		}
		return param
	case *time.Time:
		param := "-infinity"
		if t != nil {
			param = fmt.Sprintf("{date}%s", t.Format(time.RFC3339Nano))
		}
		return param
	case *int:
		param := "0"
		if t != nil {
			param = fmt.Sprintf("%d", *t)
		}
		return param
	case *float64:
		param := "0"
		if t != nil {
			param = fmt.Sprintf("%f", *t)
		}
		return param
	case string:
		return t
	case time.Time:
		return fmt.Sprintf("{date}%s", t.Format(time.RFC3339Nano))
	case float64:
		return fmt.Sprintf("%f", t)
	case int:
		return fmt.Sprintf("%d", t)
	case bool:
		return strconv.FormatBool(t)
	default:
		return "null"
	}
}
