package handlersUtil

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
)

/*
Envia uma resposta JSON contendo informações resultantes de operações aritiméticas utilizando SQL.
*/
func SendTotalEntity[E any, F any](
    w http.ResponseWriter, 
    r *http.Request, 
    filterType F,
    totalFunc func (string, F) (*E, error),
) {
	filter, ok := DecodeFilter(w, r, filterType)
	if !ok {
		return
	}
	userId, ok := GetUserId(w, r)
	if !ok {
		return
	}
	total, err := totalFunc(userId, filter)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	SendEntity(w, total)
}

/*
Envia uma resposta JSON contendo uma lista de informações resultante de um agrupamento de informações.
*/
func SendGroupedList[E any, F any](
    w http.ResponseWriter, 
    r *http.Request, 
    filterType F, 
    groupBy func(string, F) (*[]E, error),
) {
	filter, ok := DecodeFilter(w, r, filterType)
	if !ok {
		return
	}
	userId, ok := GetUserId(w, r)
	if !ok {
		return
	}
	obj, err := groupBy(userId, filter)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	SendList(w, obj)
}
