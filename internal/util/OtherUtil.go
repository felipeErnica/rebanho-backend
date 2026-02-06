package util

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
)

func GetResults(result entity.Result, resultVar any) error {
	if result.Err != nil {
		return result.Err
	}

	t := reflect.TypeOf(resultVar)
	if t.Kind() != reflect.Pointer {
		return errors.New("A variável deve ser um ponteiro")
	}

	v := reflect.ValueOf(resultVar).Elem()
	resultValue := reflect.ValueOf(result.Result)
	if resultValue.Kind() != v.Kind() {
		return fmt.Errorf("Tipo incorreto: O tipo do resultado é: %T", result.Result)
	}

	v.Set(resultValue)
	return nil
}

func ParseBool(str string) (bool, error) {
	if str == "" {
		return false, nil
	}
	bln, err := strconv.ParseBool(str)
	return bln, err
}

func FormatWarningMessage(messages ...string) string {
	msgBody := FormatMessageBody(messages...)
	return "As seguintes ocorrências foram detectadas:" + msgBody
}

func FormatMessageBody(messages ...string) string {
	formatedMsg := []string{}
	for i, message := range messages {
		msg := fmt.Sprintf("\n%d - %s", i+1, message)
		formatedMsg = append(formatedMsg, msg)
	}
	resultMsg := strings.Join(formatedMsg, "")
	return resultMsg
}
