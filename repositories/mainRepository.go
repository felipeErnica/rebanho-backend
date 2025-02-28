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
const PAGE_LIMIT int = 500

type Page struct {}

func encodeCursor(createdAt time.Time, uuid string) string {
	key := fmt.Sprintf("%s,%s", createdAt.Format(time.RFC3339Nano), uuid)
	return base64.StdEncoding.EncodeToString([]byte(key))
}

func decodeCursor(cursor string) (createdAt time.Time, id string, err error) {
	byt, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return
	}

	arrKey := strings.Split(string(byt), ",")
	if len(arrKey) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

    formatDate, err:= time.Parse(time.RFC3339Nano, arrKey[0])
    if err != nil {
        return
    }

	return formatDate, arrKey[1], err
}

func InitRepository(dbConn *sql.DB) {
    db = dbConn
	util.LogInfo("O Repositório foi iniciado com sucesso!")
}

func selectQueryList(query string, args ...any) (*sql.Rows, error) {
    util.LogInfo("Enviando query: " + strings.ReplaceAll(query, "\n", " "))
    return db.Query(query, args...)
}

func selectQueryOne(query string, args ...any) *sql.Row {
    util.LogInfo("Enviando query: " + strings.ReplaceAll(query, "\n", " "))
    return db.QueryRow(query, args...)
}
func execQuery(query string, args ...any) error {
    util.LogInfo("Enviando query: " + strings.ReplaceAll(query, "\n", " "))
    _, err := db.Exec(query, args...)
    return err
}
