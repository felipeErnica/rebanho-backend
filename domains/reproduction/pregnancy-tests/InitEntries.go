package pregnancyTests

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitTests(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := TestEntryHandler{repository}

	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/pregnancy-rate", handler.GetPregnancyRates)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/birth-rate", handler.GetBirthRates)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/test-hist", handler.GetTestHist)
	util.LogDomainsInit("Entradas de Toque")
}
