package main

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/db"
	"github.com/felipeErnica/rebanho-backend/handlers"
	"github.com/felipeErnica/rebanho-backend/middlewares"
	"github.com/felipeErnica/rebanho-backend/repositories"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
	"github.com/felipeErnica/rebanho-backend/util"
	_ "github.com/lib/pq"
)

func main() {
    
    util.LogInfo("Iniciando server....")

    app:=app.NewApp()
    app.UseGroup(middlewares.AuthenticationMiddleware)

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

    handlers.InitHandlers(app)
    repositories.InitRepository(db)
    
    err = app.ListenAndServe("localhost:8080")

    if err != nil {
        util.LogError("Não foi possível conectar a porta do host especificado!")
        serverErrors.InitServerError(err)
    }

    util.LogInfo("Server encerrado com sucesso!")

}
