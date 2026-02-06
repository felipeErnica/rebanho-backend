package main

import (
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/config/middlewares"
	"github.com/felipeErnica/rebanho-backend/internal/db"
	"github.com/felipeErnica/rebanho-backend/internal/domains"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {

	log.LogInfo("Iniciando server....", false)

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
		log.LogError("Não foi possível conectar ao banco de dados!")
		log.InitServerError(err)
	}

	err = db.Ping()

	if err != nil {
		log.LogError("Não foi possível conectar ao banco de dados!")
		log.InitServerError(err)
	}

    domains.InitDomains(app)
	err = app.ListenAndServe("localhost:8080")

	if err != nil {
		log.LogError("Não foi possível conectar a porta do host especificado!")
		log.InitServerError(err)
	}

	log.LogInfo("Server encerrado com sucesso!", false)
}
