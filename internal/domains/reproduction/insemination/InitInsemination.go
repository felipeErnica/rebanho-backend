package insemination

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitInsemination(app *app.App) {
	repository := NewEntryRepository(app.DBconn)
	service := NewService(repository)
	handler := InseminationHandler{service}

	app.HandleFunc("GET /insemination/stats/pregnancy-rate", handler.GetPregnancyRateStats)
	app.HandleFunc("GET /insemination/stats/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /insemination/stats/insemination-hist", handler.GetInseminationHist)
	app.HandleFunc("GET /insemination/stats/future-births", handler.GetFutureBirths)
	app.HandleFunc("GET /insemination/stats/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /insemination/stats/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /insemination/stats/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /insemination/stats/last-entries", handler.GetLastEntries)

	app.HandleFunc("GET /insemination/page", handler.FindEntriesPage)
	app.HandleFunc("GET /insemination/page/foot", handler.GetEntriesFoot)

	app.HandleFunc("POST /insemination", handler.Add)
	app.HandleFunc("PUT /insemination", handler.Update)
	app.HandleFunc("DELETE /insemination", handler.Delete)

	app.HandleFunc("GET /insemination/groups/page", handler.FindGroups)
	app.HandleFunc("GET /insemination/groups/{inseminationDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /insemination/groups/{inseminationDate}/foot", handler.GetEntriesByGroupFoot)
	app.HandleFunc("PUT /insemination/groups", handler.UpdateGroup)
	app.HandleFunc("DELETE /insemination/groups", handler.DeleteGroup)

	log.LogDomainsInit("Entradas de Inseminação")
}
