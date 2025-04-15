package util

import (
	"bytes"
	"fmt"
)

type InsertQuery struct {
	fields []string
	from   string
}

func NewInsertQuery(tablename string, fields ...string) *InsertQuery {
	return &InsertQuery{
		fields: fields,
		from:   tablename,
	}
}

func (q *InsertQuery) AppendFields(fields ...string) *InsertQuery {
	q.fields = append(q.fields, fields...)
	return q
}

func (q *InsertQuery) Build() string {
	var buffer bytes.Buffer
	buffer.WriteString("insert into " + q.from)
	buffer = q.buildFieldStatement(buffer)
	buffer = q.buildValuesStatement(buffer)
	return buffer.String()
}

func (q *InsertQuery) buildFieldStatement(buffer bytes.Buffer) bytes.Buffer {
	buffer.WriteString("(" + q.fields[0])
	for i := 1; i < len(q.fields); i++ {
		buffer.WriteString(", " + q.fields[i])
	}
	buffer.WriteString(")")
	return buffer
}

func (q *InsertQuery) buildValuesStatement(buffer bytes.Buffer) bytes.Buffer {
	buffer.WriteString("($1")
	for i := 2; i == len(q.fields); i++ {
		value := fmt.Sprintf(", $%d", i)
		buffer.WriteString(value)
	}
	buffer.WriteString(")")
	return buffer
}
