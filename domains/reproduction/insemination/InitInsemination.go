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
	app.HandleFunc("GET /reproduction/insemination/dashboard/future-births", handler.GetFutureBirths)
	app.HandleFunc("GET /reproduction/insemination/dashboard/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /reproduction/insemination/dashboard/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /reproduction/insemination/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /reproduction/insemination/dashboard/last-entries", handler.GetLastEntries)

	app.HandleFunc("POST /reproduction/insemination/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /reproduction/insemination/entries/page/foot", handler.GetEntriesFoot)
	app.HandleFunc("PUT /reproduction/insemination/entries/add", handler.AddInsemination)
	app.HandleFunc("PUT /reproduction/insemination/entries/replace", handler.ReplaceInsemination)
	app.HandleFunc("POST /reproduction/insemination/entries/{id}/delete", handler.Delete)
	app.HandleFunc("POST /reproduction/insemination/entries/{id}/delete-no-validation", handler.DeleteNoValidation)
	app.HandleFunc("POST /reproduction/insemination/entries/{id}/delete-change-father", handler.DeleteAndChangeFather)
	app.HandleFunc("PUT /reproduction/insemination/entries/update", handler.Update)
	app.HandleFunc("PUT /reproduction/insemination/entries/update-no-validation", handler.UpdateNoValidation)

	app.HandleFunc("GET /reproduction/insemination/groups/page", handler.FindGroups)
	app.HandleFunc("GET /reproduction/insemination/groups/{inseminationDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /reproduction/insemination/groups/{inseminationDate}/entries/foot", handler.GetEntriesByGroupFoot)

	app.HandleFunc("GET /reproduction/insemination/bulls/search", handler.SearchInseminationBulls)
	util.LogDomainsInit("Entradas de Inseminação")
}
