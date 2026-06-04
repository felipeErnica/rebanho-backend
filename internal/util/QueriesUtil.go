package util

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/jmoiron/sqlx"
)

/*
Executa um comando SQL,
usando um objeto mapeado como parâmetro
e retorna o id alterado
*/
func NamedExecReturningId[E any](db *sqlx.DB, query string, obj *E) (string, error) {
	query += " RETURNING id"
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	row, err := db.NamedQuery(query, obj)
	if err != nil {
		return "", err
	}

	id := ""
	if row.Next() {
		err := row.Scan(&id)
		if err != nil {
			return "", err
		}
	}

	return id, nil
}

/*Executa um comando SQL, usando um objeto mapeado como parâmetro*/
func NamedExec[E any](db *sqlx.DB, query string, obj *E) error {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	_, err := db.NamedExec(query, obj)
	if err != nil {
		return err
	}
	return nil
}

/*Envia um comando SQL a uma transação, usando um objeto mapeado como parâmetro*/
func NamedExecTx[E any](tx *sqlx.Tx, query string, obj *E) error {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	_, err := tx.NamedExec(query, obj)
	if err != nil {
		return err
	}
	return nil
}

/*
Envia um comando SQL a uma transição,
usando um objeto mapeado como parâmetro
e retorna o id alterado
*/
func NamedExecReturningIdTx[E any](tx *sqlx.Tx, query string, obj *E) (string, error) {
	query += " RETURNING id"
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	row, err := tx.NamedQuery(query, obj)
	if err != nil {
		return "", err
	}

	id := ""
	if row.Next() {
		err := row.Scan(&id)
		if err != nil {
			return "", err
		}
	}
	row.Close()

	return id, nil
}

/*Executa um comando SQL*/
func Exec(db *sqlx.DB, query string, args ...any) error {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	_, err := db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

/*Executa um comando SQL e envia a uma transação*/
func ExecTx(tx *sqlx.Tx, query string, args ...any) error {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	_, err := tx.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}

/*Retorna um objeto da Tabela SQL de acordo com os parâmetros informados*/
func GetOne[E any](db *sqlx.DB, query string, args ...any) (*E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	var list []E
	err := db.Select(&list, query, args...)
	if err != nil {
		return nil, err
	}

	if len(list) == 0 {
		return nil, nil
	}

	return &list[0], err
}

/*Retorna um objeto do banco, dentro de uma transação, de acordo com os parâmetros informados*/
func GetOneTx[E any](tx *sqlx.Tx, query string, object *E, args ...any) (*E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	err := tx.Get(object, query, args...)
	if err != nil {
		return nil, err
	}
	return object, err
}

/*Retorna um objeto do banco, dentro de uma transação, de acordo com o objeto enviado*/
func NamedGetTx[E any](tx *sqlx.Tx, query string, object E, arg any) (*E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	rows, err := tx.NamedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	var result E
	if rows.Next() {
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, err
}

/*Retorna um objeto do banco, dentro de uma transação, de acordo com o objeto enviado*/
func NamedGet[E any](db *sqlx.DB, query string, object E, arg any) (*E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	rows, err := db.NamedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	var result E
	if rows.Next() {
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
	}

	return &result, err
}

/*Retorna uma list do banco, dentro de uma transação, de acordo com o objeto enviado*/
func NamedQueryTx[E any, F any](tx *sqlx.Tx, query string, object E, arg F) (*[]E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	rows, err := tx.NamedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	list := []E{} 
	for rows.Next() {
		var result E
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		list = append(list, result)
	}

	return &list, err
}

/*Retorna um objeto do banco, dentro de uma transação, de acordo com o objeto enviado*/
func NamedQuery[E any, F any](db *sqlx.DB, query string, object E, arg F) (*[]E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	rows, err := db.NamedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	list := []E{} 
	for rows.Next() {
		var result E
		err := rows.StructScan(&result)
		if err != nil {
			return nil, err
		}
		list = append(list, result)
	}

	return &list, err
}

/*Retorna um objeto da Tabela SQL de acordo com os parâmetros informados*/
func GetPrimitive(db *sqlx.DB, query string, dest any, args ...any) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Pointer {
		return errors.New("A variável deve ser um ponteiro")
	}

	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	err := db.Get(dest, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

/*Retorna um objeto da Tabela SQL de acordo com os parâmetros informados*/
func NamedPrimitive(db *sqlx.DB, query string, dest any, arg any) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Pointer {
		return errors.New("A variável deve ser um ponteiro")
	}

	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	rows, err := db.NamedQuery(query, arg)
	if err != nil {
		return err
	}

	if rows.Next() {
		err = rows.Scan(dest)
		if err != nil {
			return err
		}
	}

	return nil
}

/*Retorna um objeto da Tabela SQL, através de uma transação, de acordo com os parâmetros informados*/
func GetPrimitiveTx(tx *sqlx.Tx, query string, dest any, args ...any) error {
	t := reflect.TypeOf(dest)
	if t.Kind() != reflect.Pointer {
		return errors.New("A variável deve ser um ponteiro")
	}

	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	err := tx.Get(dest, query, args...)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	return nil
}

/*Retorna uma lista de objetos da Tabela SQL de acordo com os parâmetros informados*/
func GetList[E any](db *sqlx.DB, query string, args ...any) (*[]E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	object := []E{}
	err := db.Select(&object, query, args...)
	if err != nil {
		return nil, err
	}
	return &object, err
}

/*Retorna uma lista de objetos, por transação, de acordo com os parâmetros informados*/
func GetListTx[E any](tx *sqlx.Tx, query string, dest E, args ...any) (*[]E, error) {
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)
	list := []E{}
	err := tx.Select(&list, query, args...)
	if err != nil {
		return nil, err
	}
	return &list, err
}

/*
Retorna uma página, contendo uma lista de objetos da Consulta SQL,
o número total de linhas,
um booleano indicando se há uma próxima página e
um cursor codificado indicando a próxima página.
*/
func GetPage[E any](
	db *sqlx.DB,
	query string,
	sort string,
	limit int,
	args ...any,
) (*entity.Page[E], error) {
	query = query + fmt.Sprintf(" LIMIT %d", limit)
	log.LogInfo(strings.Join(strings.Fields(query), " "), true)

	list := []E{}
	err := db.Select(&list, query, args...)
	if err != nil {
		err = errors.New("Erro na lista: " + err.Error())
		return nil, err
	}

	cursor, err := CreateCursorKey(sort, list)
	if err != nil {
		return nil, err
	}

	page := entity.Page[E]{
		List:        &list,
		NextCursor:  cursor,
		HasNextPage: len(list) == limit,
	}
	return &page, err
}
