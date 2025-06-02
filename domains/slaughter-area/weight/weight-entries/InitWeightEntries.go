package weightEntries

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitWeightEntries(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := WeightEntryHandler{repository}

	app.HandleFunc("GET /slaughter-area/weight/entries/group/{groupId}", handler.FindByGroupId)
	app.HandleFunc("GET /slaughter-area/weight/entries/{id}", handler.FindById)
	app.HandleFunc("POST /slaughter-area/weight/entries/add", handler.Add)
	app.HandleFunc("POST /slaughter-area/weight/entries/save", handler.Update)
	app.HandleFunc("DELETE /slaughter-area/weight/entries/delete/{id}", handler.Delete)
	util.LogDomainsInit("Entradas de Peso")
}
