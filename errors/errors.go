package server_errors

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/util"
)

func InternalServerError(w http.ResponseWriter) {
    w.WriteHeader(http.StatusInternalServerError)
}

func NotFoundError(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNotFound)
}

func InitServerError(err error) {
    util.LogError(err.Error())
    util.LogError("Erro na inicialização de server!")
    panic(err)
}
