package repositoriesUtil

import (
	"encoding/base64"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
)

/*Decodificação do Cursor, para obter as informações necessárias para a próxima página*/
func GetCursorArgs(cursor string) ([]any, error) {
	args := []any{}
	if cursor == "" {
		return args, nil
	}

	byt, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return args, err
	}

	arrKey := strings.Split(string(byt), ",")
	if len(arrKey) < 3 {
		err = errors.New("cursor is invalid")
		return args, err
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, arrKey[0])
	if err != nil {
		return args, err
	}

	parsedId, err := uuid.Parse(arrKey[1])
	if err != nil {
		return args, err
	}

	for i := 2; i < len(arrKey); i++ {
		if arrKey[i] == "null" {
			args = append(args, nil)
		}
		arg, err := verifyDate(arrKey[i])
		if err != nil {
			return args, err
		}
		args = append(args, arg)
	}
	args = append(args, &parsedCreatedAt, parsedId)

	return args, nil
}

// Tratamento de valores de data, verificando e apagando o prefixo de verificação
func verifyDate(arg string) (any, error) {
	if strings.HasPrefix(arg, "{date}") {
		date := strings.ReplaceAll(arg, "{date}", "")
		dateParam, err := time.Parse(time.RFC3339Nano, date)
		if err != nil {
			return "", err
		}
		return dateParam, nil
	}
	return arg, nil
}

/*Constrói e organiza os valores do filtro*/
func GetFilterArgs(filter any) []any {
	values := reflect.ValueOf(filter)
	fields := reflect.TypeOf(filter)
	if values.Kind() == reflect.Pointer {
		values = values.Elem()
		fields = fields.Elem()
	}

	args := []any{}
	for i := 1; i < values.NumField(); i++ {
		fieldValue := values.Field(i)
		fieldType := fields.Field(i)
		if !fieldValue.IsNil() {
			value := fieldValue.Elem().Interface()
			if fieldValue.Elem().Type().String() == "string" && !strings.HasSuffix(fieldType.Name, "Id") {
				value = "%" + value.(string) + "%"
			}
			if fieldValue.Elem().Kind() == reflect.Slice {
				slice := fieldValue.Elem().Interface().([]string)
				args = GetSliceArgs(slice, args)
			} else {
				args = append(args, value)
			}
		}
	}
	return args
}

func GetSliceArgs(slice []string, args []any) []any {
	for _, arg := range slice {
		args = append(args, arg)
	}
	return args
}
