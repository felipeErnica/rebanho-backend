package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitSlaughter(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := SlaughterHandler{repository}
	
	app.HandleFunc("GET /slaughter/dashboard/last-weight", handler.GetLastAverageWeight)
	app.HandleFunc("GET /slaughter/dashboard/last-dead-weight", handler.GetLastDeadWeight)
	app.HandleFunc("GET /slaughter/dashboard/last-performance", handler.GetLastPerformance)
	app.HandleFunc("GET /slaughter/dashboard/weight-hist", handler.GetWeightHist)
	app.HandleFunc("GET /slaughter/dashboard/rate-hist", handler.GetRateHist)
	app.HandleFunc("GET /slaughter/dashboard/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /slaughter/dashboard/best-mothers", handler.GetBestMothers)
	app.HandleFunc("GET /slaughter/dashboard/best-slaughterhouses", handler.GetBestSlaughterHouses)
	app.HandleFunc("GET /slaughter/dashboard/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /slaughter/dashboard/last-groups", handler.GetLastGroups)

	app.HandleFunc("POST /slaughter/info/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /slaughter/info/entries/page/foot", handler.GetEntriesPageFoot)
	app.HandleFunc("GET /slaughter/info/groups", handler.FindGroups)
	app.HandleFunc("GET /slaughter/info/groups/{entryDate}/entries", handler.FindEntriesByDate)

	util.LogDomainsInit("Entradas de Abate")
}
