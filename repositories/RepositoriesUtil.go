package repositories

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB

const PAGE_LIMIT int = 200

var userId *string

func InitRepository(dbConn *sqlx.DB) {
	db = dbConn
	util.LogInfo("O Repositório foi iniciado com sucesso!")
}

func SetUserId(id *string) {
	userId = id
}

func GetUserId() string {
	return *userId
}

func HasNextPage[E any](list []E) bool {
	return len(list) >= PAGE_LIMIT
}

func selectQueryList(query string, args ...any) (*sql.Rows, error) {
	query = strings.Join(strings.Fields(query), " ")
	println()
	util.LogInfo("Enviando query->   " + query)
	sql, err := db.Query(query, args...)
	return sql, err
}

func selectQueryOne(query string, args ...any) *sql.Row {
	query = strings.Join(strings.Fields(query), " ")
	println()
	util.LogInfo("Enviando query->   " + query)
	return db.QueryRow(query, args...)
}

func execQuery(query string, args ...any) error {
	query = strings.Join(strings.Fields(query), " ")
	util.LogInfo("Enviando query->   " + query)
	println()
	_, err := db.Exec(query, args...)
	return err
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

func decodeCursorTime(cursor string) (parsedTime *time.Time, second string, err error) {
	byt, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return
	}

	arrKey := strings.Split(string(byt), ",")
	if len(arrKey) != 2 {
		err = errors.New("cursor is invalid")
		return
	}

	if arrKey[0] == "null" {
		return nil, arrKey[1], err
	}

	parsed, err := time.Parse(time.RFC3339Nano, arrKey[0])
	if err != nil {
		return
	}

	return &parsed, arrKey[1], err
}

type SelectionProps[E any] struct {
	Dest        *[]E
	Query       string
	IsFiltered  bool
	IsFirstPage bool
	FirstParam  any
	SecondParam time.Time
	FilterArgs  []any
}

func SelectPage[E any](props SelectionProps[E]) error {
	query := strings.Join(strings.Fields(props.Query), " ")

    fmt.Println(query)
	if props.IsFiltered && props.IsFirstPage {
		err := db.Select(props.Dest, query, GetUserId(), props.FilterArgs)
		return err
	}

	if !props.IsFiltered && props.IsFirstPage {
		err := db.Select(props.Dest, query, GetUserId())
		return err
	}

	if !props.IsFiltered && !props.IsFirstPage {
		err := db.Select(props.Dest, query, GetUserId(), props.FirstParam, props.SecondParam)
		return err
	}

	if props.IsFiltered && !props.IsFirstPage {
		err := db.Select(props.Dest, query, GetUserId(), props.FirstParam, props.SecondParam, props.FilterArgs)
		return err
	}

	return errors.New("Formato Inválido")
}
