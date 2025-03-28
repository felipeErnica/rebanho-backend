package util

import (
	"bytes"
	"fmt"
)

type QueryConstructor struct {
    mainBody        string
    whereField      []string
    orderField      []string
    limit           string
}

func (q *QueryConstructor) FromQuery(query string) *QueryConstructor {
    q.mainBody = query
    return q
}

func (q *QueryConstructor) Select(alias string, fields... string) *QueryConstructor {
    var buffer bytes.Buffer
    buffer.WriteString("SELECT ")

    if alias != "" {
        fieldDeclaration:=fmt.Sprintf(" %s.%s", alias, fields[0])
        buffer.WriteString(fieldDeclaration)
        for i:=1; i < len(fields); i++ {
            fieldDeclaration:=fmt.Sprintf(", %s.%s", alias, fields[i])
            buffer.WriteString(fieldDeclaration)
        }
        q.mainBody = buffer.String()
        return q
    }

    fieldDeclaration:=fmt.Sprintf("%s", fields[0])
    buffer.WriteString(fieldDeclaration)
    for i:=1; i < len(fields); i++ {
        fieldDeclaration:=fmt.Sprintf(", %s", fields[i])
        buffer.WriteString(fieldDeclaration)
    }
    
    q.mainBody = buffer.String()
    return q
}

func (q *QueryConstructor) SelectMax(alias string, fields... string) *QueryConstructor {
    var buffer bytes.Buffer
    buffer.WriteString("SELECT ")

    if alias != "" {
        fieldDeclaration:=fmt.Sprintf(" max(%s.%s)", alias, fields[0])
        buffer.WriteString(fieldDeclaration)
        for i:=1; i < len(fields); i++ {
            fieldDeclaration:=fmt.Sprintf(", max(%s.%s)", alias, fields[i])
            buffer.WriteString(fieldDeclaration)
        }
        q.mainBody = buffer.String()
        return q
    }

    fieldDeclaration:=fmt.Sprintf("max(%s)", fields[0])
    buffer.WriteString(fieldDeclaration)
    for i:=1; i < len(fields); i++ {
        fieldDeclaration:=fmt.Sprintf(", max(%s)", fields[i])
        buffer.WriteString(fieldDeclaration)
    }
    
    q.mainBody = buffer.String()
    return q
}

func (q *QueryConstructor) SelectMin(alias string, fields... string) *QueryConstructor {
    var buffer bytes.Buffer
    buffer.WriteString("SELECT ")

    if alias != "" {
        fieldDeclaration:=fmt.Sprintf(" min(%s.%s)", alias, fields[0])
        buffer.WriteString(fieldDeclaration)
        for i:=1; i < len(fields); i++ {
            fieldDeclaration:=fmt.Sprintf(", min(%s.%s)", alias, fields[i])
            buffer.WriteString(fieldDeclaration)
        }
        q.mainBody = buffer.String()
        return q
    }

    fieldDeclaration:=fmt.Sprintf("min(%s)", fields[0])
    buffer.WriteString(fieldDeclaration)
    for i:=1; i < len(fields); i++ {
        fieldDeclaration:=fmt.Sprintf(", min(%s)", fields[i])
        buffer.WriteString(fieldDeclaration)
    }
    
    q.mainBody = buffer.String()
    return q
}

func (q *QueryConstructor) AndSelect(alias string, fields... string) *QueryConstructor {
    var buffer bytes.Buffer

    if alias != "" {
        for i:= range fields {
            fieldDeclaration:=fmt.Sprintf(", %s.%s", alias, fields[i])
            buffer.WriteString(fieldDeclaration)
        }
        q.mainBody+= buffer.String()
        return q
    }

    for i:= range fields {
        fieldDeclaration:=fmt.Sprintf(", %s.%s", alias, fields[i])
        buffer.WriteString(fieldDeclaration)
    }
    q.mainBody+= buffer.String()
    return q
}

func (q *QueryConstructor) From(tableName string, alias string) *QueryConstructor {
    fromField:= fmt.Sprintf("\nFROM %s", tableName)
    if alias != "" {
        fromField = fmt.Sprintf("%s AS %s", fromField, alias)
    }
    q.mainBody += fromField
    return q
}

func (q *QueryConstructor) LeftJoin(tableName string, alias string) *QueryConstructor {
    joinField:="\nLEFT JOIN " + tableName
    if alias != "" {
        joinField+= " AS " + alias 
    }
    q.mainBody+=joinField
    return q
}

func (q *QueryConstructor) RightJoin(tableName string, alias string) *QueryConstructor {
    joinField:="\nRIGHT JOIN " + tableName
    if alias != "" {
        joinField+= " AS " + alias 
    }
    q.mainBody+=joinField
    return q
}

func (q *QueryConstructor) Join(tableName string, alias string) *QueryConstructor {
    joinField:="\nJOIN " + tableName
    if alias != "" {
        joinField+= " AS " + alias 
    }
    q.mainBody+=joinField
    return q
}

func (q *QueryConstructor) On(firstField string, secondField string) *QueryConstructor {
    joinField:=fmt.Sprintf(" ON %s = %s", firstField, secondField)
    q.mainBody+=joinField
    return q
}

func (q *QueryConstructor) Limit(limit int) *QueryConstructor {
    q.limit+=fmt.Sprintf("\nLIMIT %d", limit)
    return q
}

func (q *QueryConstructor) Update(tableName string, fields... string) *QueryConstructor {
    var buffer bytes.Buffer
    numParam:=2
    fieldDeclaration:=fmt.Sprintf(" %s = $%d", fields[0], numParam)
    buffer.WriteString(fieldDeclaration)
    numParam++
    for i:=1; i < len(fields); i++ {
        fieldDeclaration:=fmt.Sprintf(", %s = $%d", fields[i], numParam)
        buffer.WriteString(fieldDeclaration)
        numParam++
    }
    q.mainBody = fmt.Sprintf("UPDATE %s SET %s WHERE id = $1", tableName, buffer.String())
    return q
}

func (q *QueryConstructor) Delete(tableName string) *QueryConstructor {
    q.mainBody = fmt.Sprintf("UPDATE %s SET deleted_at = $1 WHERE id = $2", tableName)
    return q
}

func (q *QueryConstructor) Insert(tableName string, fields... string) *QueryConstructor {
    var bufferFields bytes.Buffer
    bufferFields.WriteString(fields[0])

    var bufferValues bytes.Buffer
    numParam:=1
    valueDeclaration:=fmt.Sprintf("$%d", numParam)
    bufferValues.WriteString(valueDeclaration)
    numParam++

    for i:=1;i<len(fields);i++ {
        bufferFields.WriteString(", " + fields[i])
        valueDeclaration:=fmt.Sprintf(", $%d", numParam)
        bufferValues.WriteString(valueDeclaration)
        numParam++
    }

    q.mainBody = fmt.Sprintf("INSERT INTO %s (%s) VALUES(%s)", tableName, bufferFields.String(), bufferValues.String())
    return q
}

func (q *QueryConstructor) Where(condition string) *QueryConstructor {
    condition = "\nWHERE " + condition
    q.whereField = append(q.whereField, condition)
    return q
}

func (q *QueryConstructor) And(condition string) *QueryConstructor {
    condition = "\nAND " + condition
    q.whereField = append(q.whereField, condition)
    return q
}

func (q *QueryConstructor) Or(condition string) *QueryConstructor {
    condition = "\nOR " + condition
    q.whereField = append(q.whereField, condition)
    return q
}

func (q *QueryConstructor) Order(condition string) *QueryConstructor {
    condition = "\nORDER BY " + condition
    q.orderField = append(q.orderField, condition)
    return q
}

func (q *QueryConstructor) AndOrder(condition string) *QueryConstructor {
    condition = ", " + condition
    q.orderField = append(q.orderField, condition)
    return q
}

func (q *QueryConstructor) Build() string {
    var buffer bytes.Buffer
    buffer.WriteString(q.mainBody)
    for _, element:= range q.whereField {
        buffer.WriteString(element)
    }
    for _, element:= range q.orderField {
        buffer.WriteString(element)
    }
    if q.limit != "" {
        buffer.WriteString(q.limit)
    }
    return buffer.String()
}
