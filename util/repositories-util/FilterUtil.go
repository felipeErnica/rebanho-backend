package repositoriesUtil

import (
	"reflect"
)

func isFiltered(filter any) bool {
	filterValue := reflect.ValueOf(filter)
	if filterValue.Kind() == reflect.Pointer {
		filterValue = filterValue.Elem()
	}
	return filterValue.FieldByName("IsFiltered").Bool()
}
