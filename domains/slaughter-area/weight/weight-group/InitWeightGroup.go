package weightGroup

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitWeightGroup(app *app.App) {
	repository := NewRepository(app.DBconn, app.GetUserId())
	handler := WeightGroupHandler{repository}

	app.HandleFunc("GET slaughter-area/weight/groups/", handler.FindAll)
	app.HandleFunc("GET slaughter-area/weight/groups/{id}", handler.FindById)
	app.HandleFunc("POST slaughter-area/weight/groups/add", handler.Add)
	app.HandleFunc("POST slaughter-area/weight/groups/save", handler.Update)
	app.HandleFunc("DELETE slaughter-area/weight/groups/delete/{id}", handler.Delete)
	util.LogDomainsInit("Grupos de Pesagem")
}
