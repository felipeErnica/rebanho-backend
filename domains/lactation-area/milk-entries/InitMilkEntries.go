package milkEntries

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitMilkEntries(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := MilkHandler{repository}

	app.HandleFunc("GET /lactation-area/milkEntries/page", handler.FindPage)
	app.HandleFunc("GET /lactation-area/milkEntries/animal/{animalId}", handler.FindByCow)
	app.HandleFunc("POST /lactation-area/milkEntries/add", handler.Add)
	app.HandleFunc("POST /lactation-area/milkEntries/save", handler.Update)
	app.HandleFunc("DELETE /lactation-area/milkEntries/delete/{id}", handler.Delete)
	util.LogDomainsInit("Entradas de Leite")
}
