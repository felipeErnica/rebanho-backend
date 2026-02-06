package cors

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/config/middlewares"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

//Inicia o end-point necessário para verificar os testes de options do front-end

func InitCorsOptions(app *app.App) {
	app.HandleFuncNoMiddleware("OPTIONS /", middlewares.CorsMiddleware(HandleOptionsTest))
	log.LogDomainsInit("CORS")
}

func HandleOptionsTest(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
