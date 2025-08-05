package util

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/felipeErnica/rebanho-backend/entity"
)

func GetResults(result entity.Result, resultVar any) (error) {
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
