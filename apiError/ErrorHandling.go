package apiError

import (
	"encoding/json"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/util"
)

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

func JsonServerError(err error, w http.ResponseWriter) {
	util.LogError("Falha ao codificar JSON!")
	util.LogError(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func DatabaseSendError(err error, w http.ResponseWriter) {
	util.LogError("Falha ao enviar dados ao banco de dados!")
	util.LogError(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func DatabaseGetError(err error, w http.ResponseWriter) {
	util.LogError("Falha ao recuperar dados do banco de dados!")
	util.LogError(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func WriteAPIError(err *APIError, w http.ResponseWriter) {
	util.LogError("Erro de API!")
	util.LogError(err.Message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
