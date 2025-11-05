package pastureEntries

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitPastureEntries(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := PastureEntryHandler{repository}

	app.HandleFunc("GET /farm-area/pastures/{pastureId}/search/animals", handler.SearchPastureAnimals)
	app.HandleFunc("POST /farm-area/pastures/{pastureId}/entries", handler.FindByPasture)
	app.HandleFunc("POST /farm-area/pastures/{pastureId}/entries/total", handler.FindByPastureTotal)
	app.HandleFunc("GET /farm-area/pasture-entries/animal/{animalId}", handler.FindByAnimalId)
	util.LogDomainsInit("Entradas no Lote")
}
