package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitSlaughter(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := SlaughterHandler{repository}
	
	app.HandleFunc("GET /slaughter/dashboard/last-avg-weight", handler.GetLastAverageWeight)
	app.HandleFunc("GET /slaughter/dashboard/last-performance", handler.GetLastPerformance)
	app.HandleFunc("GET /slaughter/dashboard/slaughter-graph", handler.GetSlaughterGraph)
	app.HandleFunc("GET /slaughter/dashboard/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /slaughter/dashboard/best-mothers", handler.GetBestMothers)
	app.HandleFunc("GET /slaughter/dashboard/best-slaughterhouses", handler.GetBestSlaughterHouses)
	app.HandleFunc("GET /slaughter/dashboard/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /slaughter/dashboard/last-groups", handler.GetLastGroups)

	app.HandleFunc("POST /slaughter/info/entries/page", handler.FindEntriesPage)
	app.HandleFunc("GET /slaughter/info/groups", handler.FindGroups)
	app.HandleFunc("GET /slaughter/info/groups/{entryDate}/entries", handler.FindEntriesByDate)

	util.LogDomainsInit("Entradas de Abate")
}
