package slaughtergroups

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitSlaughterGroup(app *app.App) {
	repo := NewRepository(app.DBconn)
	service := NewService(repo)
	handler := NewHandler(service)

	app.HandleFunc("GET /slaughter/group", handler.FindAll)
	app.HandleFunc("GET /slaughter/group/{id}", handler.FindById)

	log.LogDomainsInit("Grupos de Abate")
}
