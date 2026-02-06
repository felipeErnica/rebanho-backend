package weight

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitWeight(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := WeightHandler{repository}

	app.HandleFunc("GET /weight/dashboard/gain-hist", handler.GetWeightGainHist)
	app.HandleFunc("GET /weight/dashboard/weight-hist", handler.GetWeightHist)
	app.HandleFunc("GET /weight/dashboard/last-gain", handler.GetLastWeightGain)
	app.HandleFunc("GET /weight/dashboard/last-weight", handler.GetLastWeight)
	app.HandleFunc("GET /weight/dashboard/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /weight/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /weight/dashboard/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /weight/dashboard/best-mothers", handler.GetBestMothers)

	app.HandleFunc("GET /weight/entries/page", handler.FindEntriesPage)
	app.HandleFunc("GET /weight/entries/page/foot", handler.GetEntriesFoot)
	app.HandleFunc("GET /weight/groups/page", handler.FindGroups)
	app.HandleFunc("GET /weight/groups/{entryDate}/entries", handler.FindEntriesByDate)
	app.HandleFunc("GET /weight/groups/{entryDate}/entries/foot", handler.GetEntriesByDateFoot)

	app.HandleFunc("DELETE /weight/entries/{id}/delete", handler.Delete)
	app.HandleFunc("PUT /weight/entries/update", handler.Update)
	log.LogDomainsInit("Entradas de Peso")
}
