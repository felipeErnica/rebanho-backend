package pasture

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitPasture(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := PastureHandler{repository}

	app.HandleFunc("GET /farm-area/pastures/search", handler.SearchPasture)
	app.HandleFunc("GET /farm-area/pastures/search-all", handler.SearchAllPastures)
	app.HandleFunc("GET /farm-area/pastures/{id}/animals", handler.FindAnimalsByPasture)
	util.LogDomainsInit("Pastos")
}
