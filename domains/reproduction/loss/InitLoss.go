package loss

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitLoss(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := PregnancyLossHandler{repository}

	app.HandleFunc("GET /reproduction/losses/page", handler.FindPage)
	app.HandleFunc("GET /reproduction/losses/animal/{animalId}", handler.FindByAnimalId)
	app.HandleFunc("GET /reproduction/losses/{id}", handler.FindById)
	app.HandleFunc("POST /reproduction/losses/add", handler.Add)
	app.HandleFunc("POST /reproduction/losses/save", handler.Update)
	app.HandleFunc("DELETE /reproduction/losses/delete/{id}", handler.Delete)
	util.LogDomainsInit("Perdas de Parição")
}
