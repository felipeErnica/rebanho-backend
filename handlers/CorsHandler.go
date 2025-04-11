package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/middlewares"
)

//Inicia o end-point necessário para verificar os testes de options do front-end

func InitCorsOptions(app *app.App) {
    app.HandleFuncNoMiddleware("OPTIONS /", middlewares.CorsMiddleware(HandleOptionsTest))
    LogControllersInit("CORS")
}

func HandleOptionsTest(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}
