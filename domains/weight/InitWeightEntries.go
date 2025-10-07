package weight

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitWeight(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := WeightHandler{repository}

	app.HandleFunc("GET /weight/dashboard/year-gain", handler.GetYearWeightGain)
	app.HandleFunc("GET /weight/dashboard/year-weight", handler.GetYearWeight)
	app.HandleFunc("GET /weight/dashboard/last-gain", handler.GetLastWeightGain)
	app.HandleFunc("GET /weight/dashboard/last-weight", handler.GetLastWeight)
	app.HandleFunc("GET /weight/dashboard/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /weight/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /weight/dashboard/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /weight/dashboard/best-mothers", handler.GetBestMothers)

	app.HandleFunc("POST /weight/info/entries/page", handler.FindEntriesPage)
	app.HandleFunc("GET /weight/info/groups", handler.FindGroups)
	app.HandleFunc("GET /weight/info/groups/{entryDate}/entries", handler.FindEntriesByDate)
	util.LogDomainsInit("Entradas de Peso")
}
