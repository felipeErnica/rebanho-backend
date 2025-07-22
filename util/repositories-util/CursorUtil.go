package repositoriesUtil

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

/*
Cria um novo cursor com a informação do último objeto da lista.
Os parâmetros são selecionados com base na coluna de ordenamento e na data de criação
*/
func CreateCursorKey[E any](sort string, list []E) (cursor string, err error) {
	listSize := len(list)
	if listSize == 0 {
		err = errors.New("A lista está vazia")
		return
	}

	sortFields := strings.Split(sort, ",")
	lastEntry := list[len(list)-1]
	entryType := reflect.TypeOf(lastEntry)
	values := reflect.ValueOf(lastEntry)
	cursorArgs := []string{}

	for _, sortField := range sortFields {
		arg, err := getValueFromSortField(sortField, entryType, values)
		if err != nil {
			return "", err
		}
		cursorArgs = append(cursorArgs, arg)
	}

	createdAt := values.FieldByName("CreatedAt").Interface()
	id := values.FieldByName("Id").Interface()
	castedCreatedAt, ok := createdAt.(time.Time)
	if !ok {
		err = errors.New("Formato de CreatedAt não é data")
		return
	}
	castedId, ok := id.(uuid.UUID)
	if !ok {
		err = errors.New("Formato de Id não é uuid")
		return
	}

	args := []string{castedCreatedAt.Format(time.RFC3339Nano), castedId.String()}
	args = append(args, cursorArgs...)

	data := ""
	for _, arg := range args {
		data = data + arg + ","
	}
	data = strings.TrimSuffix(data, ",")

	cursor = base64.StdEncoding.EncodeToString([]byte(data))
	return cursor, err
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
		err := errors.New(fmt.Sprintf("O campo não existe: %s", sortField))
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

/*Decodificação do Cursor, para obter as informações necessárias para a próxima página*/
func DecodeCursor(cursor string) (any, *time.Time, uuid.UUID, error) {
	if cursor == "" {
		return nil, nil, uuid.Nil, nil
	}

	byt, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, nil, uuid.Nil, err
	}

	arrKey := strings.Split(string(byt), ",")
	if len(arrKey) != 3 {
		err = errors.New("cursor is invalid")
		return nil, nil, uuid.Nil, err
	}

	parsedCreatedAt, err := time.Parse(time.RFC3339Nano, arrKey[1])
	if err != nil {
		return nil, nil, uuid.Nil, err
	}

	parsedId, err := uuid.Parse(arrKey[2])
	if err != nil {
		return nil, nil, uuid.Nil, err
	}

	if arrKey[0] == "null" {
		return nil, &parsedCreatedAt, parsedId, nil
	}

	//Tratamento de valores de data, verificando e apagando o prefixo
	if strings.HasPrefix(arrKey[0], "{date}") {
		date := strings.ReplaceAll(arrKey[0], "{date}", "")
		firstParam, err := time.Parse(time.RFC3339Nano, date)
		if err != nil {
			return nil, nil, uuid.Nil, err
		}
		return firstParam, &parsedCreatedAt, parsedId, err
	}
	return arrKey[0], &parsedCreatedAt, parsedId, err
}
