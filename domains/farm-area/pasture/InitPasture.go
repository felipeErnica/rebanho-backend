package pasture

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitPasture(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := PastureHandler{repository}

	app.HandleFunc("GET /farm-area/pastures", handler.FindAll)
	app.HandleFunc("GET /farm-area/pastures/{id}", handler.FindById)
	app.HandleFunc("POST /farm-area/pastures/add", handler.Add)
	app.HandleFunc("POST /farm-area/pastures/save", handler.Update)
	app.HandleFunc("DELETE /farm-area/pastures/delete/{id}", handler.Delete)
	util.LogDomainsInit("Pastos")
}
