package tests

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitTests(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewService(repository)
	handler := TestEntryHandler{service}

	app.HandleFunc("GET /tests/stats/pregnancy-rate", handler.GetPregnancyRates)
	app.HandleFunc("GET /tests/stats/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /tests/stats/birth-rate", handler.GetBirthRates)
	app.HandleFunc("GET /tests/stats/test-hist", handler.GetTestHist)
	app.HandleFunc("GET /tests/stats/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /tests/stats/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /tests/stats/next-births", handler.GetNextBirths)
	app.HandleFunc("GET /tests/stats/ranked-results", handler.GetRankedResults)

	app.HandleFunc("GET /tests/page", handler.FindEntriesPage)
	app.HandleFunc("GET /tests/page/foot", handler.GetEntriesFoot)

	app.HandleFunc("POST /tests", handler.Add)
	app.HandleFunc("PUT /tests", handler.Update)
	app.HandleFunc("DELETE /tests/{id}", handler.Delete)

	app.HandleFunc("GET /tests/group/page", handler.FindGroups)
	app.HandleFunc("GET /tests/group/{testDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /tests/group/{testDate}/entries/foot", handler.GetEntriesByGroupFoot)

	app.HandleFunc("PUT /tests/group", handler.UpdateGroup)
	app.HandleFunc("DELETE /tests/group/{testDate}", handler.DeleteGroup)

	log.LogDomainsInit("Entradas de Toque")
}
