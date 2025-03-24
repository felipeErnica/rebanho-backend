package serverErrors

import (
	"errors"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/util"
)

func EmptyList() error {
    err:= errors.New("A matriz está vazia")
    return err
}

func InternalServerError(w http.ResponseWriter) {
    w.WriteHeader(http.StatusInternalServerError)
}

func NotFoundError(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNotFound)
}


func AuthenticationError(w http.ResponseWriter, err error) {
    util.LogError(err.Error())
    w.WriteHeader(http.StatusUnauthorized)
}

func InitServerError(err error) {
    util.LogError(err.Error())
    util.LogError("Erro na inicialização de server!")
    panic(err)
}
