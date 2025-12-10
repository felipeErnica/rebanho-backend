package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitSlaughter(app *app.App) {
	repository := NewSlaughterRepository(app.DBconn)
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

	app.HandleFunc("POST /slaughter/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /slaughter/entries/page/foot", handler.GetEntriesPageFoot)
	app.HandleFunc("DELETE /slaughter/entries/{id}/delete", handler.Delete)
	app.HandleFunc("PUT /slaughter/entries/update", handler.Update)
	app.HandleFunc("PUT /slaughter/entries/add", handler.Add)
	app.HandleFunc("PUT /slaughter/entries/replace", handler.Replace)

	app.HandleFunc("GET /slaughter/groups", handler.FindGroups)
	app.HandleFunc("GET /slaughter/groups/{entryDate}/entries", handler.FindEntriesByDate)
	app.HandleFunc("GET /slaughter/groups/{entryDate}/entries/foot", handler.GetEntriesByDateFoot)

	util.LogDomainsInit("Entradas de Abate")
}

func InitButcher(app *app.App) {
	repository := newButcherRepository(app.DBconn)
	handler := ButcherHandler{repository}
	
	app.HandleFunc("PUT /slaughter/butchers/add", handler.Add)
	app.HandleFunc("PUT /slaughter/butchers/update", handler.Update)
	app.HandleFunc("PUT /slaughter/butchers/replace", handler.Replace)
	app.HandleFunc("GET /slaughter/butchers/search", handler.Search)
	app.HandleFunc("GET /slaughter/butchers/find-all", handler.FindAll)
	app.HandleFunc("DELETE /slaughter/butchers/{id}/delete", handler.Delete)

	util.LogDomainsInit("Frigoríficos")
}
