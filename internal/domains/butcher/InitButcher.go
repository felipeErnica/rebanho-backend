package butcher

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitButcher(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewService(repository)
	handler := ButcherHandler{service}

	app.HandleFunc("POST /butchers", handler.Add)
	app.HandleFunc("PUT /butchers", handler.Update)
	app.HandleFunc("DELETE /butchers", handler.Delete)

	app.HandleFunc("GET /butchers", handler.FindAll)
	app.HandleFunc("GET /butchers/{id}", handler.FindById)

	log.LogDomainsInit("Frigoríficos")
}
