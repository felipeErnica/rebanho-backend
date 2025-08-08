package insemination

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitInsemination(app *app.App) {
	repository := NewEntryRepository(app.DBconn)
    handler := InseminationHandler{repository}
	app.HandleFunc("GET /reproduction/insemination/dashboard/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /reproduction/insemination/dashboard/insemination-hist", handler.GetInseminationHist)
	app.HandleFunc("GET /reproduction/insemination/dashboard/best-bull", handler.GetBestBull)
    util.LogDomainsInit("Entradas de Inseminação")
}
