package util

import (
	"bytes"
	"fmt"
)

type UpdateQuery struct {
	fields []string
	from   string
}

func NewUpdateQuery(tablename string, fields ...string) *UpdateQuery {
	return &UpdateQuery{
		fields: fields,
		from: tablename,
	}
}

func (q *UpdateQuery) AppendFields(fields ...string) *UpdateQuery {
	q.fields = append(q.fields, fields...)
	return q
}

func (q *UpdateQuery) Build() string {
	var buffer bytes.Buffer
	buffer.WriteString("update " + q.from)
	buffer.WriteString(" set " + q.fields[0] + " = $2")
	numParams := 3
	for i := 1; i < len(q.fields); i++ {
		statement := fmt.Sprintf(", %s = %d", q.fields[i], numParams)
		buffer.WriteString(statement)
		numParams++
	}
	buffer.WriteString(" where id = $1")
	return buffer.String()
}
