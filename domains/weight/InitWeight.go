package weight

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
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

	app.HandleFunc("POST /weight/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /weight/entries/page/foot", handler.GetEntriesFoot)
	app.HandleFunc("GET /weight/groups/page", handler.FindGroups)
	app.HandleFunc("GET /weight/groups/{entryDate}/entries", handler.FindEntriesByDate)
	app.HandleFunc("GET /weight/groups/{entryDate}/entries/foot", handler.GetEntriesByDateFoot)

	app.HandleFunc("DELETE /weight/entries/{id}/delete", handler.Delete)
	app.HandleFunc("PUT /weight/entries/update", handler.Update)
	util.LogDomainsInit("Entradas de Peso")
}
