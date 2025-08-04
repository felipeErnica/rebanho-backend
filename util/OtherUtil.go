package util

import (
	"errors"

	"github.com/felipeErnica/rebanho-backend/entity"
)

func GetResults[T any](result entity.Result, resultType T) (T, error) {
    if result.Err != nil {
        return resultType, result.Err
    }

    resultContent, ok := result.Result.(T); if !ok {
        err := errors.New("Falha ao obter um resultado assícrono: Tipo Incorreto")
        return resultType, err
    }
    return resultContent, nil
}
