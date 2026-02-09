package pasture

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitPasture(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewService(repository)
	handler := PastureHandler{service}

	app.HandleFunc("GET /pastures/search", handler.Search)
	app.HandleFunc("GET /pastures/{id}/animals", handler.FindAnimalsById)
	log.LogDomainsInit("Pastos")
}
