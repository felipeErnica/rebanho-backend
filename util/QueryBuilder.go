package util

import "fmt"

type QueryBuilder struct{}

const PAGE_LIMIT int = 500

func (q QueryBuilder) getNullOrdering(direction string) string {
	nullDirection := "DESC NULLS LAST"
	if direction == "asc" {
		nullDirection = "ASC NULLS FIRST"
	}
	return nullDirection
}

func (q QueryBuilder) getSignal(direction string) string {
	signal := ">"
	if direction == "desc" {
		signal = "<"
	}
	return signal
}

func (q QueryBuilder) GetFirstPage(mainQuery string, firstField string, secondField string,
	deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE %[2]s IS NULL
        ORDER BY %[3]s %[4]s, %[5]s %[4]s
        LIMIT %[6]d
    `, mainQuery, deletedField, firstField, nullOrdering, secondField, PAGE_LIMIT)
}

func (q QueryBuilder) GetNextPage(mainQuery string, firstField string, secondField string,
	deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	signal := q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NULL
            AND (%[3]s, %[4]s) %[5]s ($1, $2)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        LIMIT %d
        `, mainQuery, deletedField, firstField, secondField, signal, nullOrdering, PAGE_LIMIT)
}

func (q QueryBuilder) GetNextPageNull(mainQuery string, firstField string, secondField string,
	deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	signal := q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NULL
            AND (%[3]s IS NULL AND %[4]s %[5]s $1)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        LIMIT %d
        `, mainQuery, deletedField, firstField, secondField, signal, nullOrdering, PAGE_LIMIT)
}

func (q *QueryBuilder) GetFirstPageConditional(mainQuery string,conditionalField string, firstField string, secondField string, 
    deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NULL
            AND %[3]s = $1
        ORDER BY %[4]s %[5]s, %[6]s %[5]s
        LIMIT %[7]d
    `, mainQuery, deletedField, conditionalField, firstField, nullOrdering, secondField, PAGE_LIMIT)
}

func (q *QueryBuilder) GetNextPageNullConditional(mainQuery string, conditionalField string, firstField string,
    secondField string, deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	signal := q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NULL
            AND %[3]s = $1
            AND (%[4]s IS NULL AND %[5]s %[6]s $2)
        ORDER BY %[4]s %[7]s, %[5]s %[7]s
        LIMIT %d
        `, mainQuery, deletedField, conditionalField, firstField, secondField, signal, nullOrdering, PAGE_LIMIT)
}

func (q *QueryBuilder) GetNextPageConditional(mainQuery string, conditionalField string, firstField string, secondField string, 
    deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	signal := q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NULL
            AND %[3]s = $1
            AND (%[4]s, %[5]s) %[6]s ($2, $3)
        ORDER BY %[4]s %[7]s, %[5]s %[7]s
        LIMIT %d
        `, mainQuery, deletedField, conditionalField, firstField, secondField, signal, nullOrdering, PAGE_LIMIT)
}

func (q QueryBuilder) GetDeletedFirstPage(mainQuery string, firstField string, secondField string,
	deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE %[2]s IS NOT NULL
        ORDER BY %[3]s %[4]s, %[5]s %[4]s
        LIMIT %[6]d
    `, mainQuery, deletedField, firstField, nullOrdering, secondField, PAGE_LIMIT)
}

func (q QueryBuilder) GetDeletedNextPage(mainQuery string, firstField string,
	secondField string, deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	signal := q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NOT NULL
            AND (%[3]s, %[4]s) %[5]s ($1, $2)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        LIMIT %d
        `, mainQuery, deletedField, firstField, secondField, signal, nullOrdering, PAGE_LIMIT)
}

func (q QueryBuilder) GetDeletedNextPageNull(mainQuery string, firstField string, secondField string,
	deletedField string, direction string) string {
	nullOrdering := q.getNullOrdering(direction)
	signal := q.getSignal(direction)
	return fmt.Sprintf(`
        %[1]s
        WHERE 
            %[2]s IS NOT NULL
            AND (%[3]s IS NULL AND %[4]s %[5]s $1)
        ORDER BY %[3]s %[6]s, %[4]s %[6]s
        LIMIT %d
        `, mainQuery, deletedField, firstField, secondField, signal, nullOrdering, PAGE_LIMIT)
}
