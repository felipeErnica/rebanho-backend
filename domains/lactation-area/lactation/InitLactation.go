package lactation

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitLactation(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := LactationHandler{repository}

	app.HandleFunc("GET /lactation-area/lactation/page", handler.FindPage)
	app.HandleFunc("GET /lactation-area/lactation/animal/{animalId}", handler.FindByCow)
	app.HandleFunc("POST /lactation-area/lactation/add", handler.Add)
	app.HandleFunc("POST /lactation-area/lactation/save", handler.Save)
	app.HandleFunc("DELETE /lactation-area/lactation/delete/{id}", handler.Delete)
	util.LogDomainsInit("Lactações")
}
