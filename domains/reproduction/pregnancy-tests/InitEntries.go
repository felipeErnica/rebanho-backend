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
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/next-births", handler.GetNextBirths)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/best-results", handler.GetBestResults)
	util.LogDomainsInit("Entradas de Toque")
}
