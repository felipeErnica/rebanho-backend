package inseminationGroup

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitGroup(app *app.App) {
	repository := NewRepository(app.DBconn)
    handler := GroupHandler{repository}

	app.HandleFunc("GET /reproduction/insemination/groups", handler.FindAll)
	app.HandleFunc("GET /reproduction/insemination/groups/{id}", handler.FindById)
	app.HandleFunc("POST /reproduction/insemination/groups", handler.Add)
	app.HandleFunc("POST /reproduction/insemination/groups/save", handler.Save)
	app.HandleFunc("DELETE /reproduction/insemination/groups/{id}", handler.Delete)
	util.LogDomainsInit("Grupos de Inseminação")
}
