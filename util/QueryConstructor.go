package util

import (
	"bytes"
	"fmt"
	"log/slog/internal/buffer"
)

type StatementType int

const (
	MAX StatementType = iota
	MIN
	SELECT
	UPDATE
	DELETE
	INSERT
)

var selectMap = map[StatementType]string{
	MAX:    "max",
	MIN:    "min",
	SELECT: "select",
	UPDATE: "update",
	DELETE: "delete",
	INSERT: "insert",
}

type FieldGroup struct {
	alias  string
	fields []string
}

func NewSelectGroup(fields ...string) *FieldGroup {
	return &FieldGroup{
		fields: fields,
	}
}

func (s *FieldGroup) Alias(alias string) *FieldGroup {
	s.alias = alias
	return s
}

type QueryConstructor struct {
	statementType   StatementType
	fieldsGroups    []FieldGroup
	fromStatement   string
	joinGroups      []string
	whereStatements []string
	orderField      []string
	limit           int
}

func NewStatement(statementType StatementType, groups ...FieldGroup) *QueryConstructor {
	q := &QueryConstructor{
		statementType: statementType,
		fieldsGroups:  groups,
	}
	return q
}

func (q *QueryConstructor) AppendFields(groups ...FieldGroup) *QueryConstructor {
	q.fieldsGroups = append(q.fieldsGroups, groups...)
	return q
}

func (q *QueryConstructor) From(statement string) *QueryConstructor {
	q.fromStatement = statement
	return q
}

func (q *QueryConstructor) Joins(joins ...string) *QueryConstructor {
	q.joinGroups = joins
	return q
}

func (q *QueryConstructor) AppendJoins(joins ...string) *QueryConstructor {
	q.joinGroups = append(q.joinGroups, joins...)
	return q
}

func (q *QueryConstructor) Where(whereStatements ...string) *QueryConstructor {
	q.whereStatements = whereStatements
	return q
}

func (q *QueryConstructor) AppendWhere(whereStatements ...string) *QueryConstructor {
	q.whereStatements = append(q.whereStatements, whereStatements...)
	return q
}

func (q *QueryConstructor) OrderBy(orderBy ...string) *QueryConstructor {
	q.orderField = orderBy
	return q
}

func (q *QueryConstructor) AppendOrders(orderBy ...string) *QueryConstructor {
	q.orderField = append(q.orderField, orderBy...)
	return q
}

func (q *QueryConstructor) Limit(limit int) *QueryConstructor {
	q.limit = limit
	return q
}

func (q *QueryConstructor) Build() string {
	var buffer bytes.Buffer

	return buffer.String()
}

func (q *QueryConstructor) buildMainStatement(buffer bytes.Buffer) bytes.Buffer {

	switch q.statementType {
	case SELECT:
		buffer = q.selectBuild(buffer)
	case UPDATE:

	}

	buffer.WriteString(fieldDeclaration)
	for i := 1; i < len(fields); i++ {
		fieldDeclaration := fmt.Sprintf(", %s.%s", alias, fields[i])
		buffer.WriteString(fieldDeclaration)
	}
	q.mainBody = buffer.String()
	return q
}

func (q *QueryConstructor) selectBuild(buffer bytes.Buffer) bytes.Buffer {
	buffer.WriteString("SELECT ")
	buffer = q.selectGroupDeclaration(q.fieldsGroups[0], buffer)
	for i := 1; i < len(q.fieldsGroups); i++ {
		buffer.WriteString(", ")
		buffer = q.selectGroupDeclaration(q.fieldsGroups[i], buffer)
	}
	return buffer
}

func (q *QueryConstructor) selectGroupDeclaration(group FieldGroup, buffer bytes.Buffer) bytes.Buffer {
	firstDeclaration := q.selectFieldDeclaration(group.fields[0], group.alias)
	buffer.WriteString(firstDeclaration)
	for i := 1; i < len(group.fields); i++ {
		fieldDeclaration := fmt.Sprintf(", %s", q.selectFieldDeclaration(group.fields[i], group.alias))
		buffer.WriteString(fieldDeclaration)
	}
	return buffer
}

func (q *QueryConstructor) selectFieldDeclaration(field string, alias string) string {
	fieldDeclaration := field
	if alias != "" {
		fieldDeclaration = fmt.Sprintf("%s.%s", alias, field)
	}
	if q.statementType == MAX || q.statementType == MIN {
		fieldDeclaration = fmt.Sprintf("%s(%s)", q.statementType, fieldDeclaration)
	}
	return field
}

func (q *QueryConstructor) updateBuild(buffer bytes.Buffer) bytes.Buffer {
	buffer.WriteString("SELECT ")
	buffer = q.selectGroupDeclaration(q.fieldsGroups[0], buffer)
	for i := 1; i < len(q.fieldsGroups); i++ {
		buffer.WriteString(", ")
		buffer = q.selectGroupDeclaration(q.fieldsGroups[i], buffer)
	}
	return buffer
}

func (q *QueryConstructor) selectGroupDeclaration(group FieldGroup, buffer bytes.Buffer) bytes.Buffer {
	firstDeclaration := q.selectFieldDeclaration(group.fields[0], group.alias)
	buffer.WriteString(firstDeclaration)
	for i := 1; i < len(group.fields); i++ {
		fieldDeclaration := fmt.Sprintf(", %s", q.selectFieldDeclaration(group.fields[i], group.alias))
		buffer.WriteString(fieldDeclaration)
	}
	return buffer
}




