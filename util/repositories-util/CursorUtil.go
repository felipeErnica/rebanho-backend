package repositoriesUtil

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func CreateCursorKey[E any](sort string, list []E) (cursor string, err error) {
	listSize := len(list)
	if listSize == 0 {
		err = errors.New("A lista está vazia")
		return
	}

	lastEntry := list[len(list)-1]
	entryType := reflect.TypeOf(lastEntry)
	values := reflect.ValueOf(lastEntry)
	var value any

	for i := range entryType.NumField() {
		field := entryType.Field(i)
		tag := field.Tag.Get("db")
		if tag == "" {
			err = errors.New("O campo não existe")
			return
		}
		if strings.EqualFold(sort, tag) {
			value = values.Field(i).Interface()
		}
	}

	firstParam := getFirstParam(value)
	createdAt := values.FieldByName("CreatedAt").Interface()
	secondParam, ok := createdAt.(time.Time)
	if !ok {
		err = errors.New("Formato de CreatedAt não é data")
		return
	}

    data := fmt.Sprintf("%s,%s", firstParam, secondParam.Format(time.RFC3339Nano))
    fmt.Println(data)
    cursor = base64.StdEncoding.EncodeToString([]byte(data))
	return cursor, err
}

func getFirstParam(value any) string {
	switch t := value.(type) {
    case *string:
		param := "null"
		if t != nil {
			param = *t
		}
		return param
    case string:
		return t
	case *time.Time:
		param := "null"
		if t != nil {
			param = fmt.Sprintf("date%s", t.Format(time.RFC3339Nano)) 
		}
		return param
	case time.Time:
		return fmt.Sprintf("date%s", t.Format(time.RFC3339Nano))
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

func DecodeCursor(cursor string) (firstParam any, secondParam time.Time, err error) {
    if cursor == "" {
        return
    }

	byt, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return
	}

	arrKey := strings.Split(string(byt), ",")
	if len(arrKey) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

	parsedSecond, err := time.Parse(time.RFC3339Nano, arrKey[1])
	if err != nil {
		return arrKey[0], parsedSecond, err
	}

	if arrKey[0] == "null" {
		return nil, parsedSecond, err
	}

	if strings.HasPrefix(arrKey[0], "date") {
        date := strings.ReplaceAll(arrKey[0], "date", "")
		firstParam, err = time.Parse(time.RFC3339Nano, date)
		if err != nil {
			return
		}
		return firstParam, parsedSecond, err
	}

	firstParam = arrKey[0]
	return firstParam, parsedSecond, err
}
