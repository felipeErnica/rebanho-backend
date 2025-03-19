package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/util"
)

func InitHandlers(mux *http.ServeMux) {
    InitAnimal(mux)
    InitLactation(mux)
    InitMilk(mux)
    InitPastureEntryHandler(mux)
    InitPastureHandler(mux)
    InitPregnancyLossHandler(mux)
    InitWeightEntries(mux)
    InitWeightGroups(mux)
    InitSlaugherhouses(mux)
}

func LogControllersInit(name string) {
    util.LogInfo("Requisições de " + name + " iniciadas com sucesso!")
}

func JsonServerError(err error, w http.ResponseWriter) {
    util.LogError("Falha ao decodificar JSON!")
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

func writeResponse(w http.ResponseWriter, response []byte) {
    w.Header().Set("Content-Type","application/json")
    w.Write(response)
}
