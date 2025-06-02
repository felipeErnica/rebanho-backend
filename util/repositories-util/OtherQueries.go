package repositoriesUtil

import (
	"bytes"
	"fmt"
	"reflect"

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
		buffer.WriteString(":" + fieldNames[i] + ",")
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
	_, err := db.NamedExec(query, object)
	return object, err
}

/*Retorna um objeto da Tabela SQL de acordo com os parâmetros informados*/
func GetOne[E any](db *sqlx.DB, query string, args ...any) (*E, error) {
	var object E
	err := db.Get(&object, query, args...)
	if err != nil {
		return nil, err
	}
	return &object, err
}

/*Retorna uma lista de objetos da Tabela SQL de acordo com os parâmetros informados*/
func GetList[E any](db *sqlx.DB, query string, args ...any) (*[]E, error) {
	object := []E{}
	err := db.Select(&object, query, args...)
	if err != nil {
		return nil, err
	}
	return &object, err
}
