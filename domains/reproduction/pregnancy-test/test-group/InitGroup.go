package testGroup

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitGroup(app *app.App) {
    repository:= NewRepository(app.DBconn, app.GetUserId())
    handler:= TestGroupHandler{repository}

    app.HandleFunc("GET reproduction/pregnancy-group",  handler.FindAll)
    app.HandleFunc("GET reproduction/pregnancy-group/{id}",  handler.FindById)
    app.HandleFunc("POST reproduction/pregnancy-group/add",  handler.Add)
    app.HandleFunc("POST reproduction/pregnancy-group/save",  handler.Update)
    app.HandleFunc("DELETE reproduction/pregnancy-group/delete/{id}",  handler.Delete)
    util.LogDomainsInit("Grupo de Toque")
}
