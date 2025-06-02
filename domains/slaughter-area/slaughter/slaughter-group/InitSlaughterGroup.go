package slaughterGroup

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitSlaughterGroup(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := SlaughterGroupHandler{repository}

	app.HandleFunc("GET /slaughter-area/slaughter/groups", handler.FindAll)
	app.HandleFunc("GET /slaughter-area/slaughter/groups/{id}", handler.FindById)
	app.HandleFunc("POST /slaughter-area/slaughter/groups/add", handler.Add)
	app.HandleFunc("POST /slaughter-area/slaughter/groups/save", handler.Update)
	app.HandleFunc("DELETE /slaughter-area/slaughter/groups/delete/{id}", handler.Delete)
	util.LogDomainsInit("Frigoríficos")
}
