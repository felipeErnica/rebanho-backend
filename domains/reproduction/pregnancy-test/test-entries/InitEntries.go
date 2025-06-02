package testEntries

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitEntries(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := TestEntryHandler{repository}

	app.HandleFunc("GET /reproduction/pregnancy-test/entries/groups/{groupId}", handler.FindByGroupId)
	app.HandleFunc("GET /reproduction/pregnancy-test/entries/animals/{animalId}", handler.FindByAnimalId)
	app.HandleFunc("GET /reproduction/pregnancy-test/entries/{id}", handler.FindById)
	app.HandleFunc("POST /reproduction/pregnancy-test/entries/add", handler.Add)
	app.HandleFunc("POST /reproduction/pregnancy-test/entries/save", handler.Update)
	app.HandleFunc("DELETE /reproduction/pregnancy-test/entries/delete/{id}", handler.Delete)
	util.LogDomainsInit("Entradas de Toque")
}
