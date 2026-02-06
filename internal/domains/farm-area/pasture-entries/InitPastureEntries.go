package pastureEntries

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitPastureEntries(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := PastureEntryHandler{repository}

	app.HandleFunc("GET /farm-area/pastures/{pastureId}/search/animals", handler.SearchPastureAnimals)
	app.HandleFunc("GET /farm-area/pastures/{pastureId}/entries", handler.FindByPasture)
	app.HandleFunc("GET /farm-area/pastures/{pastureId}/entries/total", handler.FindByPastureTotal)
	app.HandleFunc("GET /farm-area/entries/animal/{animalId}", handler.FindByAnimalId)
	app.HandleFunc("PUT /farm-area/entries/add", handler.AddEntry)
	app.HandleFunc("PUT /farm-area/entries/transfer", handler.TransferEntry)
	log.LogDomainsInit("Entradas no Lote")
}
