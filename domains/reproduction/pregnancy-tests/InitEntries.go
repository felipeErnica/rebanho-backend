package pregnancyTests

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitTests(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := TestEntryHandler{repository}

	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/pregnancy-rate", handler.GetPregnancyRates)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/birth-rate", handler.GetBirthRates)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/test-hist", handler.GetTestHist)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/next-births", handler.GetNextBirths)
	app.HandleFunc("GET /reproduction/pregnancy-test/dashboard/ranked-results", handler.GetRankedResults)

	app.HandleFunc("POST /reproduction/pregnancy-test/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /reproduction/pregnancy-test/entries/page/foot", handler.GetEntriesFoot)
	app.HandleFunc("PUT /reproduction/pregnancy-test/entries/add", handler.Add)
	app.HandleFunc("PUT /reproduction/pregnancy-test/entries/replace", handler.Replace)
	app.HandleFunc("PUT /reproduction/pregnancy-test/entries/update", handler.Update)
	app.HandleFunc("DELETE /reproduction/pregnancy-test/entries/{id}/delete", handler.Delete)

	app.HandleFunc("GET /reproduction/pregnancy-test/group/page", handler.FindGroups)
	app.HandleFunc("GET /reproduction/pregnancy-test/group/{testDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /reproduction/pregnancy-test/group/{testDate}/entries/foot", handler.GetEntriesByGroupFoot)
	app.HandleFunc("PUT /reproduction/pregnancy-test/group/{testDate}/update", handler.UpdateBatch)
	app.HandleFunc("DELETE /reproduction/pregnancy-test/group/{testDate}/delete", handler.DeleteBatch)

	util.LogDomainsInit("Entradas de Toque")
}
