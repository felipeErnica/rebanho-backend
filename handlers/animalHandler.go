package handlers

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/felipeErnica/rebanho-backend/repositories"
)

type AnimalHandler struct {
    Repository repositories.AnimalRepository
}

func InitAnimal(mux *http.ServeMux, db *sql.DB) {

    repository:= repositories.AnimalRepository{ Db: db, }
    repositories.LogInitRepository("Animais")
    handler:=AnimalHandler{ Repository: repository, }

    mux.Handle("GET /animais", &handler)
    //mux.Handle("/animais/", &AnimalHandler{})
    LogControllersInit("Animais")
}

func (h *AnimalHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    switch {
        case r.Method == http.MethodGet && strings.EqualFold(r.URL.Path, "/animais"): h.GetAll(w, r)
    }
}

func (h *AnimalHandler) GetAll(w http.ResponseWriter, r *http.Request)  {
    w.Write([]byte("Lista de Animais"))
}
