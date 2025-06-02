package slaughterEntry

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitSlaughterEntry(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := SlaughterEntryHandler{repository}

	app.HandleFunc("POST slaughter-area/slaughter/entries/page", handler.FindPage)
	app.HandleFunc("GET slaughter-area/slaughter/entries/group/{groupId}", handler.FindByGroupId)
	app.HandleFunc("GET slaughter-area/slaughter/entries/{id}", handler.FindById)
	app.HandleFunc("POST slaughter-area/slaughter/entries/add", handler.Add)
	app.HandleFunc("POST slaughter-area/slaughter/entries/save", handler.Save)
	app.HandleFunc("DELETE slaughter-area/slaughter/entries/delete/{id}", handler.Delete)
	util.LogDomainsInit("Entradas de Abate")
}
