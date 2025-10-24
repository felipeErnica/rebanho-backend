package app

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/config/middlewares"
	"github.com/jmoiron/sqlx"
)

type App struct {
	mux         *http.ServeMux
	middlewares []middlewares.Middleware
	DBconn      *sqlx.DB
}

func NewApp() *App {
	return &App{mux: http.NewServeMux()}
}

func (a *App) CreateMiddlewaresGroup(mids ...middlewares.Middleware) {
	a.middlewares = append(a.middlewares, mids...)
}

func (a *App) HandleFunc(pattern string, handler http.HandlerFunc) {
	finalHandler := handler
	for _, middleware := range a.middlewares {
		finalHandler = middleware(finalHandler)
	}
	a.mux.HandleFunc(pattern, finalHandler)
}

func (a *App) HandleFuncNoMiddleware(pattern string, handler http.HandlerFunc) {
	a.mux.HandleFunc(pattern, handler)
}

func (a *App) ListenAndServe(address string) error {
	return http.ListenAndServe(address, a.mux)
}
