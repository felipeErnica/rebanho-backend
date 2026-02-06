package auth

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/config/middlewares"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitAuth(app *app.App) {
	repository := NewRepostory(app.DBconn)
	handler := UserHandler{repository}

	app.HandleFuncNoMiddleware("POST /login", middlewares.CorsMiddleware(handler.Authenticate))
	app.HandleFuncNoMiddleware("POST /register", middlewares.CorsMiddleware(handler.Register))
	log.LogDomainsInit("Autenticação")
}
