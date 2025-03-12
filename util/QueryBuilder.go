package util

import "fmt"

type QueryBuilder struct{}

func (q QueryBuilder) getSignal(direction string) string {
    signal:=">"
    if direction == "desc" {
        signal = "<"
    }
    return signal
}

func (q QueryBuilder) GetFirstPage(mainQuery string ,firstField string, secondField string,
	deletedField string, direction string) string {
	return fmt.Sprintf(`
        %[1]s
        WHERE %[2]s IS NULL
        ORDER BY %[3]s %[4]s, %[5]s %[4]s
    `, mainQuery, deletedField, firstField, direction, secondField)
}

func (q QueryBuilder) GetNextPage(mainQuery string, firstField string, secondField string, 
    deletedField string, direction string) string {
    signal:= q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NULL
            AND (%[3]s, %[4]s) %[5]s ($1, $2)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        `, mainQuery, deletedField, firstField, secondField, signal, direction)
}

func (q QueryBuilder) GetNextPageNull(mainQuery string, firstField string, secondField string, 
    deletedField string, direction string) string {
    signal:= q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NULL
            AND (%[3]s IS NULL AND %[4]s %[5]s ($1)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        `, mainQuery, deletedField, firstField, secondField, signal, direction)
}
func (q QueryBuilder) GetDeletedFirstPage(firstField string, secondField string,
	deletedField string, direction string) string {
	return fmt.Sprintf(`
        WHERE %[1]s IS NOT NULL
        ORDER BY %[2]s %[3]s, %[4]s %[3]s
    `, deletedField, firstField, direction, secondField)
}

func (q QueryBuilder) GetDeletedNextPage(mainQuery string, firstField string, 
    secondField string, deletedField string, direction string) string {
    signal:= q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NOT NULL
            AND (%[3]s, %[4]s) %[5]s ($1, $2)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        `, mainQuery, deletedField, firstField, secondField, signal, direction)
}

func (q QueryBuilder) GetDeletedNextPageNull(mainQuery string, firstField string, secondField string, 
    deletedField string, direction string) string {
    signal:= q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NOT NULL
            AND (%[3]s IS NULL AND %[4]s %[5]s ($1)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        `, mainQuery, deletedField, firstField, secondField, signal, direction)
}
