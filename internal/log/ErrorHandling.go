package log

import (
	"encoding/json"
	"net/http"
)

func WriteInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
}

func NotFoundError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
}

func AuthenticationError(w http.ResponseWriter, err error) {
	LogError(err.Error())
	w.WriteHeader(http.StatusUnauthorized)
}

func InitServerError(err error) {
	LogError(err.Error())
	LogError("Erro na inicialização de server!")
	panic(err)
}

func JsonServerError(err error, w http.ResponseWriter) {
	LogError("Falha ao codificar JSON!")
	LogError(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func DatabaseSendError(err error, w http.ResponseWriter) {
	LogError("Falha ao enviar dados ao banco de dados!")
	LogError(err.Error())
	w.WriteHeader(http.StatusInternalServerError)
}

func WriteError(w http.ResponseWriter, err error) {
	apiErr := InternalServerAPIError(err)
	WriteAPIError(apiErr, w)
}

func WriteAPIError(err *APIError, w http.ResponseWriter) {
	LogError("Erro de API!")
	LogError(err.Message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	json.NewEncoder(w).Encode(err)
}
