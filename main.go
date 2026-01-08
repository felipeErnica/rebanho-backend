package main

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/config/middlewares"
	"github.com/felipeErnica/rebanho-backend/db"
	"github.com/felipeErnica/rebanho-backend/domains"
	"github.com/felipeErnica/rebanho-backend/apiError"
	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	util.LogInfo("Iniciando server....", false)

	dataBaseInfo := db.ConnectPostgres().ReturnDatabaseInfo()
	db, err := sqlx.Open("postgres", dataBaseInfo)
	db.SetMaxOpenConns(15)
	db.SetMaxIdleConns(15)
	db.SetConnMaxLifetime(5 * time.Minute)
	defer db.Close()

	app := app.NewApp(db)
	app.CreateMiddlewaresGroup(
		middlewares.CorsMiddleware,
		middlewares.AuthenticationMiddleware,
	)

	if err != nil {
		util.LogError("Não foi possível conectar ao banco de dados!")
		apiError.InitServerError(err)
	}

	err = db.Ping()

	if err != nil {
		util.LogError("Não foi possível conectar ao banco de dados!")
		apiError.InitServerError(err)
	}

    domains.InitDomains(app)
	err = app.ListenAndServe("localhost:8080")

	if err != nil {
		util.LogError("Não foi possível conectar a porta do host especificado!")
		apiError.InitServerError(err)
	}

	util.LogInfo("Server encerrado com sucesso!", false)
}
