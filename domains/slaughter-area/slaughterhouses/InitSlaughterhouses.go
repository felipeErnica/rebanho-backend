package slaughterhouses

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitSlaughterhouse(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := SlaughterhouseHandler{repository}

	app.HandleFunc("GET /slaughter-area/slaughterhouses", handler.FindAll)
	app.HandleFunc("GET /slaughter-area/slaughterhouses/{id}", handler.FindById)
	app.HandleFunc("POST /slaughter-area/slaughterhouses/add", handler.Add)
	app.HandleFunc("POST /slaughter-area/slaughterhouses/save", handler.Update)
	app.HandleFunc("DELETE /slaughter-area/slaughterhouses/delete/{id}", handler.Delete)
	util.LogDomainsInit("Frigoríficos")
}
