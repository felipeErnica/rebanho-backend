package pasture

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitPasture(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := PastureHandler{repository}

	app.HandleFunc("GET /pastures/search", handler.SearchPasture)
	app.HandleFunc("GET /pastures/{id}/animals", handler.FindAnimalsByPasture)
	log.LogDomainsInit("Pastos")
}
