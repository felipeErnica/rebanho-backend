package repositories

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/util"
)

var db *sql.DB
const PAGE_LIMIT uint16 = 500

type Page struct {}

func encodeCursor(createdAt time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", createdAt, uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func decodeCursor(cursor string) (createdAt string, id string, err error) {
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
    util.LogInfo("Enviando query: " + query)
    return db.Query(query, args...)
}

func selectQueryOne(query string, args ...any) *sql.Row {
    util.LogInfo("Enviando query: " + query)
    return db.QueryRow(query, args...)
}
func execQuery(query string, args ...any) error {
    util.LogInfo("Enviando query: " + query)
    _, err := db.Exec(query, args...)
    return err
}
