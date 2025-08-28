package repositoriesUtil

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/jmoiron/sqlx"
)

/*Exclui um objeto da Tabela SQL usando o id como parâmetro*/
func Delete(db *sqlx.DB, tableName string, id string) error {
	deletedAt := time.Now()
	query := fmt.Sprintf("update %s set deleted_at = $1 where id = $2", tableName)
	_, err := db.Exec(query, deletedAt, id)
	return err
}

/*Salva e atualiza um objeto da Tabela SQL*/
func Update[E any](db *sqlx.DB, tableName string, object *E) error {
	query := generateUpdateQuery(object, tableName)
	_, err := db.NamedExec(query, object)
	return err
}

/*Gera uma consulta de SQL do tipo UPDATE com base no objeto fornecido*/
func generateUpdateQuery(object any, tableName string) string {
	fieldNames := getFieldsNames(object)
	var buffer bytes.Buffer
	for _, fieldName := range fieldNames {
		if fieldName != "id" {
			statement := fmt.Sprintf("%[1]s = :%[1]s, ", fieldName)
			buffer.WriteString(statement)
		}
	}
	valuesFields := buffer.String()
	valuesFields = strings.TrimSuffix(valuesFields, ", ")

	query := fmt.Sprintf("update %s set %s where id = :id", tableName, valuesFields)
	return query
}

/*Adiciona um objeto a Tabela SQL*/
func Add[E any](db *sqlx.DB, tableName string, object *E) (*E, error) {
	query := generateInsertQuery(object, tableName)
	fmt.Println(query)
	_, err := db.NamedExec(query, object)
	return object, err
}

/*Gera uma consulta de SQL do tipo INSERT com base no objeto fornecido*/
func generateInsertQuery(object any, tableName string) string {
	fieldNames := getFieldsNames(object)

	var buffer bytes.Buffer
    for _, fieldName := range fieldNames {
		buffer.WriteString(":" + fieldName + ", ")
	}
	valuesFields := buffer.String()
    valuesFields = strings.TrimSuffix(valuesFields, ", ")

	query := fmt.Sprintf("insert into %s values(%s)", tableName, valuesFields)
	return query
}

/*Obtém os nomes do campo da tabela SQL relacionada ao tipo*/
func getFieldsNames(object any) []string {
	objectTypes := reflect.TypeOf(object)
	if objectTypes.Kind() == reflect.Pointer {
		objectTypes = objectTypes.Elem()
	}

	fieldNames := []string{}

	for i := range objectTypes.NumField() {
		field := objectTypes.Field(i)
		fieldNames = append(fieldNames, field.Tag.Get("db"))
	}
	return fieldNames
}

/*Retorna um objeto da Tabela SQL de acordo com os parâmetros informados*/
func GetOne[E any](db *sqlx.DB, query string, args ...any) (*E, error) {
	util.LogInfo(strings.Join(strings.Fields(query), " "), true)
	var object E
	err := db.Get(&object, query, args...)
	if err != nil {
		return nil, err
	}
	return &object, err
}

/*Retorna uma lista de objetos da Tabela SQL de acordo com os parâmetros informados*/
func GetList[E any](db *sqlx.DB, query string, args ...any) (*[]E, error) {
	util.LogInfo(strings.Join(strings.Fields(query), " "), true)
	object := []E{}
	err := db.Select(&object, query, args...)
	if err != nil {
		return nil, err
	}
	return &object, err
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
	query = query + fmt.Sprintf(" limit %d", limit)
	util.LogInfo(strings.Join(strings.Fields(query), " "), true)

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
