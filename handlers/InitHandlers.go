package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitHandlers(app *app.App) {
    InitCorsOptions(app)
    InitAnimal(app)
    InitLactation(app)
    InitMilk(app)
    InitPastureEntryHandler(app)
    InitPastureHandler(app)
    InitPregnancyLossHandler(app)
    InitWeightEntries(app)
    InitWeightGroups(app)
    InitSlaugherhouses(app)
    InitSlaughterGroup(app)
    InitSlaughterEntry(app)
    InitPregnancyTestGroup(app)
    InitPregnancyTestEntry(app)
    InitInseminationGroup(app)
    InitInseminationEntry(app)
    InitUserAuthentication(app)
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
