package repositories

import (
	"database/sql"

	"github.com/felipeErnica/rebanho-backend/util"
)

var db *sql.DB

func InitRepository(dbConn *sql.DB) {
    db = dbConn
	util.LogInfo("O Repositório foi iniciado com sucesso!")
}

func SelectQueryList(query string, args ...any) (*sql.Rows, error) {
    util.LogInfo("Enviando query: " + query)
    return db.Query(query, args...)
}

func SelectQueryOne(query string, args ...any) *sql.Row {
    util.LogInfo("Enviando query: " + query)
    return db.QueryRow(query, args...)
}
func ExecQuery(query string, args ...any) error {
    util.LogInfo("Enviando query: " + query)
    _, err := db.Exec(query, args...)
    return err
}
