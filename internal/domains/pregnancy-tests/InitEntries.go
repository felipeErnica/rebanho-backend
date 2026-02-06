package pregnancyTests

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitTests(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewService(repository)
	handler := TestEntryHandler{service}

	app.HandleFunc("GET /pregnancy-test/stats/pregnancy-rate", handler.GetPregnancyRates)
	app.HandleFunc("GET /pregnancy-test/stats/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /pregnancy-test/stats/birth-rate", handler.GetBirthRates)
	app.HandleFunc("GET /pregnancy-test/stats/test-hist", handler.GetTestHist)
	app.HandleFunc("GET /pregnancy-test/stats/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /pregnancy-test/stats/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /pregnancy-test/stats/next-births", handler.GetNextBirths)
	app.HandleFunc("GET /pregnancy-test/stats/ranked-results", handler.GetRankedResults)

	app.HandleFunc("GET /pregnancy-test/page", handler.FindEntriesPage)
	app.HandleFunc("GET /pregnancy-test/page/foot", handler.GetEntriesFoot)

	app.HandleFunc("POST /pregnancy-test", handler.Add)
	app.HandleFunc("PUT /pregnancy-test", handler.Update)
	app.HandleFunc("DELETE /pregnancy-test/{id}", handler.Delete)

	app.HandleFunc("GET /pregnancy-test/group/page", handler.FindGroups)
	app.HandleFunc("GET /pregnancy-test/group/{testDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /pregnancy-test/group/{testDate}/entries/foot", handler.GetEntriesByGroupFoot)

	app.HandleFunc("PUT /pregnancy-test/group", handler.UpdateBatch)
	app.HandleFunc("DELETE /pregnancy-test/group/{testDate}", handler.DeleteBatch)

	log.LogDomainsInit("Entradas de Toque")
}
