package app

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/middlewares"
)

type App struct {
	mux *http.ServeMux
    middlewares []middlewares.Middleware
}

func NewApp() *App {
    return &App{mux: http.NewServeMux()}
}

func (a *App) UseGroup(mids... middlewares.Middleware) {
    a.middlewares = append(a.middlewares, mids...)
}

func (a *App) HandleFunc(pattern string, handler http.HandlerFunc) {
    finalHandler:=handler
    for _, m:= range a.middlewares {
        finalHandler = m(finalHandler)
    }
    a.mux.HandleFunc(pattern, finalHandler)
}

func (a *App) HandleFuncNoMiddleware(pattern string, handler http.HandlerFunc) {
    a.mux.HandleFunc(pattern, handler)
}

func (a *App) ListenAndServe(address string) error {
    return http.ListenAndServe(address, a.mux)
}
