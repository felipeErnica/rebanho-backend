package repositoriesUtil

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
	"github.com/jmoiron/sqlx"
)

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

/*Gera uma consulta de SQL do tipo INSERT com base no objeto fornecido*/
func generateInsertQuery(object any, tableName string) string {
	fieldNames := getFieldsNames(object)

	var buffer bytes.Buffer
	buffer.WriteString(":" + fieldNames[0])
	for i := 1; i < len(fieldNames); i++ {
		buffer.WriteString(", :" + fieldNames[i])
	}
	valuesFields := buffer.String()

	query := fmt.Sprintf("INSERT INTO %s VALUES(%s)", tableName, valuesFields)
	return query
}

/*Gera uma consulta de SQL do tipo UPDATE com base no objeto fornecido*/
func generateUpdateQuery(object any, tableName string) string {
	fieldNames := getFieldsNames(object)

	var buffer bytes.Buffer
	statement := fmt.Sprintf("%[1]s = :%[1]s", fieldNames[1])
	buffer.WriteString(statement)
	for i := 2; i < len(fieldNames); i++ {
		statement := fmt.Sprintf("%[1]s = :%[1]s", fieldNames[i])
		buffer.WriteString(statement)
	}
	valuesFields := buffer.String()

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = :id", tableName, valuesFields)
	return query
}

/*Exclui um objeto da Tabela SQL usando o id como parâmetro*/
func Delete(db *sqlx.DB, tableName string, id string) error {
	query := fmt.Sprintf("DELETE FROM %s id = $1", tableName)
	_, err := db.Exec(query)
	return err
}

/*Salva e atualiza um objeto da Tabela SQL*/
func Update[E any](db *sqlx.DB, tableName string, object *E) error {
	query := generateUpdateQuery(object, tableName)
	_, err := db.NamedExec(query, object)
	return err
}

/*Adiciona um objeto a Tabela SQL*/
func Add[E any](db *sqlx.DB, tableName string, object *E) (*E, error) {
	query := generateInsertQuery(object, tableName)
	fmt.Println(query)
	_, err := db.NamedExec(query, object)
	return object, err
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

/*Retorna uma página, contendo uma lista de objetos da Consulta SQL,
o número total de linhas,
um booleano indicando se há uma próxima página e
um cursor codificado indicando a próxima página.*/
func GetPage[E any](
    db *sqlx.DB, 
    query string, 
    countQuery string, 
    sort string, 
    cursorArgs []any, 
    args ...any,
) (*entity.Page[E], error) {
	util.LogInfo(strings.Join(strings.Fields(query), " "), true)
	util.LogInfo(strings.Join(strings.Fields(countQuery), " "), true)

    listArgs := append(args, cursorArgs...)
	list := []E{}
	err := db.Select(&list, query, listArgs...)
	if err != nil {
        err = errors.New("Erro na lista: " + err.Error())
		return nil, err
	}

    totalResult := db.QueryRow(countQuery, args...)
    var total int
    err = totalResult.Scan(&total)
    if err != nil {
        err = errors.New("Erro na contagem: " + err.Error())
        return nil, err
    }
    
    cursor, err := CreateCursorKey(sort, list)
    if err != nil {
        return nil, err
    }
    fmt.Println("cursor: ", cursor)

    page := entity.Page[E]{
        List: &list,
        NextCursor: cursor,
        HasNextPage: cursor != "",
        Total: total,
    }
	return &page, err
}
