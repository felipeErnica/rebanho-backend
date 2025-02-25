package main

import (
	"database/sql"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/db"
	server_errors "github.com/felipeErnica/rebanho-backend/errors"
	"github.com/felipeErnica/rebanho-backend/handlers"
	"github.com/felipeErnica/rebanho-backend/util"
)

func main() {
    
    util.LogInfo("Iniciando server....")

	mux := http.NewServeMux()
    handlers.InitHandlers(mux)
    
    dataBaseInfo:= db.ConnectPostgres().ReturnDatabaseInfo()
    db, err := sql.Open("postgres", dataBaseInfo)

    if err != nil {
        server_errors.InitServerError(err)
    }
    
    defer db.Close()

    err = db.Ping()

    if err != nil {
        server_errors.InitServerError(err)
    }

    err = http.ListenAndServe("localhost:8080", mux)

    if err != nil {
        server_errors.InitServerError(err)
    }

    util.LogInfo("Server encerrado com sucesso!")

}
