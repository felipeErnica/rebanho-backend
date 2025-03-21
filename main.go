package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/db"
	"github.com/felipeErnica/rebanho-backend/handlers"
	"github.com/felipeErnica/rebanho-backend/repositories"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
	"github.com/felipeErnica/rebanho-backend/util"
	_ "github.com/lib/pq"
)

func main() {
    
    util.LogInfo("Iniciando server....")

    dataBaseInfo:= db.ConnectPostgres().ReturnDatabaseInfo()
    db, err := sql.Open("postgres", dataBaseInfo)
    db.SetMaxOpenConns(15)
    db.SetMaxIdleConns(15)
    db.SetConnMaxLifetime(5 * time.Minute)
    defer db.Close()

    if err != nil {
        util.LogError("Não foi possível conectar ao banco de dados!")
        serverErrors.InitServerError(err)
    }

    err = db.Ping()

    if err != nil {
        util.LogError("Não foi possível conectar ao banco de dados!")
        serverErrors.InitServerError(err)
    }

	mux:= http.NewServeMux()
    handlers.InitHandlers(mux)
    repositories.InitRepository(db)
    
    err = http.ListenAndServe("localhost:8080", mux)

    if err != nil {
        util.LogError("Não foi possível conectar a porta do host especificado!")
        serverErrors.InitServerError(err)
    }

    util.LogInfo("Server encerrado com sucesso!")

}
