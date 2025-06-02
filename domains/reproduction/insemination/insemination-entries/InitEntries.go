package inseminationEntries

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitEntries(app *app.App) {
	repository := NewEntryRepository(app.DBconn)
    handler := EntriesHandler{repository}
	app.HandleFunc("GET /reproduction/insemination/groups/{groupId}/entries", handler.FindByGroupId)
	app.HandleFunc("GET /reproduction/insemination/entries/{id}", handler.FindById)
	app.HandleFunc("POST /reproduction/insemination/entries/add", handler.Add)
	app.HandleFunc("POST /reproduction/insemination/entries/save", handler.Update)
	app.HandleFunc("DELETE /reproduction/insemination/entries/delete/{id}", handler.Delete)
    util.LogDomainsInit("Entradas de Inseminação")
}
