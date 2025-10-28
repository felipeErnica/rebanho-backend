package naturalMating

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitMating(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := MatingHandler{repository}

	app.HandleFunc("GET /reproduction/mating/dashboard/pregnancy-rate", handler.GetPregnancyRateStats)
	app.HandleFunc("GET /reproduction/mating/dashboard/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /reproduction/mating/dashboard/insemination-hist", handler.GetInseminationHist)
	app.HandleFunc("GET /reproduction/mating/dashboard/future-births", handler.GetFutureBirths)
	app.HandleFunc("GET /reproduction/mating/dashboard/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /reproduction/mating/dashboard/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /reproduction/mating/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /reproduction/mating/dashboard/last-entries", handler.GetLastEntries)

	app.HandleFunc("POST /reproduction/mating/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /reproduction/mating/entries/page/foot", handler.GetEntriesFoot)

	app.HandleFunc("GET /reproduction/mating/groups/page", handler.FindGroups)
	app.HandleFunc("GET /reproduction/mating/groups/{inseminationDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /reproduction/mating/groups/{inseminationDate}/entries/foot", handler.GetEntriesByGroupFoot)

	app.HandleFunc("GET /reproduction/mating/bulls/search", handler.SearchInseminationBulls)
	util.LogDomainsInit("Monta Natural")
}
