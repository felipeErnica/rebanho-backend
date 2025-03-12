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

func encodeCursor(key string) string {
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

func decodeCursorTime(cursor string) (first time.Time, second string, err error) {
	byt, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return
	}

	arrKey := strings.Split(string(byt), ",")
	if len(arrKey) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

    if first, err = time.Parse(time.RFC3339Nano, arrKey[0]); err != nil {
        return
    }

    return first, arrKey[1], err
}

func isTimeField(field string, timeFields []string) bool {

    for i:=0; i < len(timeFields); i++ {
        if field == timeFields[i] {
            return true
        }
    }
    
    return false
}

func getNullStatement(direction string) string {
    nullStatment:="DESC NULLS LAST"
    if direction == "asc" {
        nullStatment = "ASC NULLS FIRST" 
    }
    return nullStatment
}

func getNextPageCriteria(firstField string, secondField string, direction string, isNullValue bool) string {
    orderDirection:= "DESC NULLS LAST"
    signal:="<"
    if direction == "asc" {
        orderDirection = "ASC NULLS FIRST"
        signal = ">"
    }

    order:= fmt.Sprintf("%s %[2]s, %s %[2]s", firstField, orderDirection)
    
    var where string
    if isNullValue {
        where = fmt.Sprintf("WHERE %s IS NULL AND %s %s $2", firstField, secondField, signal)
    } else {
        where = fmt.Sprintf("WHERE (%s,%s) %s ($1, $2)", firstField, secondField, signal)
    }

    return where + "\n" + order
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
