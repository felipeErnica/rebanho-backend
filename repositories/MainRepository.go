package repositories

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/felipeErnica/rebanho-backend/util"
)

var db *sql.DB
const PAGE_LIMIT int = 500

func encodeCursor(param string, uuid string) string {
	key := fmt.Sprintf("%s,%s", param, uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func decodeCursor(cursor string) (first string, second string, err error) {
	byt, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return
	}

	arrKey := strings.Split(string(byt), ",")
	if len(arrKey) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

    return arrKey[0], arrKey[1], err
}

func InitRepository(dbConn *sql.DB) {
    db = dbConn
	util.LogInfo("O Repositório foi iniciado com sucesso!")
}

func selectQueryList(query string, args ...any) (*sql.Rows, error) {
    query = strings.Join(strings.Fields(query)," ")
    println()
    util.LogInfo("Enviando query->   " + query)
    sql, err:= db.Query(query, args...)
    return sql, err
}

func selectQueryOne(query string, args ...any) *sql.Row {
    query = strings.Join(strings.Fields(query)," ")
    println()
    util.LogInfo("Enviando query->   " + query)
    return db.QueryRow(query, args...)
}

func execQuery(query string, args ...any) error {
    query = strings.Join(strings.Fields(query)," ")
    util.LogInfo("Enviando query->   " + query)
    println()
    _, err := db.Exec(query, args...)
    return err
}
