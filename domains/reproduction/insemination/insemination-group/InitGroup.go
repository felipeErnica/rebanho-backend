package inseminationGroup

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitGroup(app *app.App) {
	repository := NewRepository(app.DBconn, app.GetUserId())
    handler := GroupHandler{repository}

	app.HandleFunc("GET insemination/groups", handler.FindAll)
	app.HandleFunc("GET insemination/groups/{id}", handler.FindById)
	app.HandleFunc("POST insemination/groups", handler.Add)
	app.HandleFunc("POST insemination/groups/save", handler.Save)
	app.HandleFunc("DELETE insemination/groups/{id}", handler.Delete)
	util.LogDomainsInit("Grupos de Inseminação")
}
