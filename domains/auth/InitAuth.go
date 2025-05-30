package auth

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/middlewares"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitAuth(app *app.App) {
	repository := NewRepostory(app.DBconn)
	handler := UserHandler{repository}

	app.HandleFuncNoMiddleware("POST /login", middlewares.CorsMiddleware(handler.Authenticate))
	app.HandleFuncNoMiddleware("POST /register", middlewares.CorsMiddleware(handler.Register))
	util.LogDomainsInit("Autenticação")
}
