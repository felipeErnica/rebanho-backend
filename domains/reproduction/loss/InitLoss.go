package loss

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitLoss(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := LossHandler{repository}

	app.HandleFunc("GET /reproduction/losses/dashboard/loss-rate", handler.GetLossRate)
	app.HandleFunc("GET /reproduction/losses/dashboard/losses-hist", handler.GetLossesHist)
	app.HandleFunc("GET /reproduction/losses/dashboard/losses-animals", handler.GetMostLossesAnimals)

	app.HandleFunc("POST /reproduction/losses/entries/page", handler.FindPage)
	app.HandleFunc("POST /reproduction/losses/entries/page/foot", handler.GetPageFoot)

	util.LogDomainsInit("Perdas de Parição")
}
