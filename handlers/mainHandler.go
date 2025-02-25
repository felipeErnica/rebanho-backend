package handlers

import (
	"database/sql"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/util"
)

func InitHandlers(mux *http.ServeMux, db *sql.DB) {
    InitAnimal(mux, db)
}

func LogControllersInit(name string) {
    util.LogInfo("Requisições de " + name + " iniciadas com sucesso!")
}
