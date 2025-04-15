package util

import (
	"bytes"
	"fmt"
)

type StatementType int

const (
	MAX StatementType = iota
	MIN
	SELECT
)

var selectMap = map[StatementType]string{
	MAX:    "max",
	MIN:    "min",
	SELECT: "select",
}

type FieldGroup struct {
	alias  string
	fields []string
}

func NewGroup(fields ...string) *FieldGroup {
	return &FieldGroup{
		fields: fields,
	}
}

func NewNamedGroup(alias string, fields ...string) *FieldGroup {
	return &FieldGroup{
		fields: fields,
		alias: alias,
	}
}

type SelectConstructor struct {
	statementType   StatementType
	fieldsGroups    []FieldGroup
	fromStatement   string
	joinStatements  []string
	whereStatements []string
	orderStatements      []string
	limit           int
}

func NewSelectQuery(statementType StatementType, groups ...FieldGroup) *SelectConstructor {
	q := &SelectConstructor{
		statementType: statementType,
		fieldsGroups:  groups,
	}
	return q
}

func (q *SelectConstructor) AppendFields(groups ...FieldGroup) *SelectConstructor {
	q.fieldsGroups = append(q.fieldsGroups, groups...)
	return q
}

func (q *SelectConstructor) From(statement string) *SelectConstructor {
	q.fromStatement = statement
	return q
}

func (q *SelectConstructor) Joins(joins ...string) *SelectConstructor {
	q.joinStatements = joins
	return q
}

func (q *SelectConstructor) AppendJoins(joins ...string) *SelectConstructor {
	q.joinStatements = append(q.joinStatements, joins...)
	return q
}

func (q *SelectConstructor) Where(whereStatements ...string) *SelectConstructor {
	q.whereStatements = whereStatements
	return q
}

func (q *SelectConstructor) AppendWhere(whereStatements ...string) *SelectConstructor {
	q.whereStatements = append(q.whereStatements, whereStatements...)
	return q
}

func (q *SelectConstructor) OrderBy(orderBy ...string) *SelectConstructor {
	q.orderStatements = orderBy
	return q
}

func (q *SelectConstructor) AppendOrders(orderBy ...string) *SelectConstructor {
	q.orderStatements = append(q.orderStatements, orderBy...)
	return q
}

func (q *SelectConstructor) Limit(limit int) *SelectConstructor {
	q.limit = limit
	return q
}

func (q *SelectConstructor) Build() string {
	var buffer bytes.Buffer
	buffer = q.buildMainStatement(buffer)
	buffer.WriteString(" from " + q.fromStatement)

	if len(q.joinStatements) != 0 {
        buffer = q.buildJoinStatements(buffer)
    }

    if len(q.whereStatements) != 0 {
        buffer = q.buildWhereStatements(buffer)
    }

    if len(q.orderStatements) != 0 {
        buffer = q.buildOrderStatements(buffer)
    }

    if q.limit > 0 {
        limitStatement := fmt.Sprintf(" limit %d", q.limit) 
        buffer.WriteString(limitStatement)
    }

	return buffer.String()
}

func (q *SelectConstructor) buildMainStatement(buffer bytes.Buffer) bytes.Buffer {
	buffer.WriteString("select ")
	buffer = q.selectGroupDeclaration(q.fieldsGroups[0], buffer)
	for i := 1; i < len(q.fieldsGroups); i++ {
		buffer.WriteString(", ")
		buffer = q.selectGroupDeclaration(q.fieldsGroups[i], buffer)
	}
	return buffer
}

func (q *SelectConstructor) selectGroupDeclaration(group FieldGroup, buffer bytes.Buffer) bytes.Buffer {
	firstDeclaration := q.selectFieldDeclaration(group.fields[0], group.alias)
	buffer.WriteString(firstDeclaration)
	for i := 1; i < len(group.fields); i++ {
		fieldDeclaration := fmt.Sprintf(", %s", q.selectFieldDeclaration(group.fields[i], group.alias))
		buffer.WriteString(fieldDeclaration)
	}
	return buffer
}

func (q *SelectConstructor) selectFieldDeclaration(field string, alias string) string {
	fieldDeclaration := field
	if alias != "" {
		fieldDeclaration = fmt.Sprintf("%s.%s", alias, field)
	}
	if q.statementType == MAX || q.statementType == MIN {
		fieldDeclaration = fmt.Sprintf("%s(%s)", selectMap[q.statementType], fieldDeclaration)
	}
	return fieldDeclaration
}

func (q *SelectConstructor) buildJoinStatements(buffer bytes.Buffer) bytes.Buffer {
    buffer.WriteString(" ")
	for _, statement := range q.joinStatements {
		buffer.WriteString(statement + " ")
	}
	return buffer
}

func (q *SelectConstructor) buildWhereStatements(buffer bytes.Buffer) bytes.Buffer {
	buffer.WriteString(" where")
	for _, statement := range q.whereStatements {
		buffer.WriteString(" " + statement)
	}
	return buffer
}

func (q *SelectConstructor) buildOrderStatements(buffer bytes.Buffer) bytes.Buffer {
	buffer.WriteString(" order by ")
	buffer.WriteString(q.orderStatements[0])
	for i := 1; i < len(q.orderStatements); i++ {
		buffer.WriteString(", " + q.orderStatements[i])
	}
	return buffer
}
