package util

import (
	"fmt"
)

type DeleteQuery struct {
	from   string
}

func NewDeleteQuery(tablename string) string {
	return fmt.Sprintf("update %s set deleted_at = $1 where id = $2", tablename)
}
