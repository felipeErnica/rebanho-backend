package insemination

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitInsemination(app *app.App) {
	repository := NewEntryRepository(app.DBconn)
    handler := InseminationHandler{repository}

	app.HandleFunc("GET /reproduction/insemination/dashboard/pregnancy-rate", handler.GetPregnancyRateStats)
	app.HandleFunc("GET /reproduction/insemination/dashboard/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /reproduction/insemination/dashboard/insemination-hist", handler.GetInseminationHist)
	app.HandleFunc("GET /reproduction/insemination/dashboard/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /reproduction/insemination/dashboard/pregnants-number", handler.GetPregnantsNumber)
	app.HandleFunc("GET /reproduction/insemination/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /reproduction/insemination/dashboard/last-entries", handler.GetLastEntries)

	app.HandleFunc("POST /reproduction/insemination/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /reproduction/insemination/entries/page/foot", handler.GetEntriesFoot)

	app.HandleFunc("GET /reproduction/insemination/groups/page", handler.FindGroups)
	app.HandleFunc("GET /reproduction/insemination/groups/page/foot", handler.GetGroupsFoot)
	app.HandleFunc("GET /reproduction/insemination/groups/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /reproduction/insemination/groups/entries/foot", handler.GetEntriesByGroupFoot)

	app.HandleFunc("GET /reproduction/insemination/bulls/search", handler.SearchInseminationBulls)
    util.LogDomainsInit("Entradas de Inseminação")
}
