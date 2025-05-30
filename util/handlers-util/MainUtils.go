package handlersUtil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
)

/*
Retorna os parâmetros da página.
*/
func getPageParameters(r *http.Request) (sort string, order string, cursor string) {
	cursor = r.URL.Query().Get("cursor")
	sort = r.URL.Query().Get("sort")
	order = r.URL.Query().Get("order")

	if sort == "" {
		sort = "id"
	}

	if order == "" {
		order = "asc"
	}
	return cursor, sort, order
}

/*
Decodifica a entidae contida no corpo da solicitação HTTP
e retorna um erro caso o formato esteja incorreto.
*/
func decodeEntity[E any](w http.ResponseWriter, r *http.Request, entity *E) {
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		err = errors.New(fmt.Sprintf("Falha na decodificação da entidade: %s", err.Error()))
		serverErrors.JsonServerError(err, w)
		return
	}
}

/*
Envia a entidade como resposta HTTP.
*/
func SendEntity[E any](w http.ResponseWriter, model *E) {
	response, err := json.Marshal(model)
	if err != nil {
		serverErrors.JsonServerError(err, w)
		return
	}
	writeResponse(w, response)
}

/*
Envia uma lista como resposta HTTP.
*/
func SendList[E any](w http.ResponseWriter, model *[]E) {
	response, err := json.Marshal(model)
	if err != nil {
		serverErrors.JsonServerError(err, w)
		return
	}
	writeResponse(w, response)
}

/*
Decodifica o filtro contido no corpo da solicitação HTTP
e retorna um erro caso o formato esteja incorreto.
*/
func DecodeFilter[F any](w http.ResponseWriter, r *http.Request, filter *F) {
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		err = errors.New(fmt.Sprintf("Falha na decodificação do filtro: %s", err.Error()))
		serverErrors.JsonServerError(err, w)
		return
	}
}

/*
Retorna uma página como resposta HTTP, usando o repositório e o filtro como parâmetros.
O repositório deve conter uma função FindPage.
*/
func ReturnPage[E any, F any](w http.ResponseWriter, r *http.Request, repository PageRepository[E, F], filter F) {
	sort, order, cursor := getPageParameters(r)
	page, err := repository.FindPage(sort, order, cursor, filter)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	SendEntity(w, page)
}

/*
Retorna uma entidade como resposta HTTP pelo id fornecido, o repositório deve conter uma função FindById.
*/
func FindById[E any](w http.ResponseWriter, r *http.Request, repository RepositoryFindById[E]) {
	id := r.PathValue("id")
	entity, err := repository.FindById(id)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	SendEntity(w, entity)
}

/*
Retorna uma lista com todas as entidades no banco como resposta HTTP, o repositório deve conter uma função FindAll.
*/
func FindAll[E any](w http.ResponseWriter, r *http.Request, repository RepositoryFindAll[E]) {
	list, err := repository.FindAll()
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	SendList(w, list)
}

/*
Salva uma nova entidade no banco, e retorna o resultado na Resposta HTTP
*/
func Add[E any](w http.ResponseWriter, r *http.Request, repository RepositoryAdd[E]) {
	var obj E
	decodeEntity(w, r, &obj)

	model, err := repository.Add(&obj)
	if err != nil {
		serverErrors.DatabaseSendError(err, w)
		return
	}
	response, err := json.Marshal(model)
	if err != nil {
		serverErrors.JsonServerError(err, w)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

/*
Atualiza uma entidade no banco, e retorna o resultado na Resposta HTTP
*/
func Update[E any](w http.ResponseWriter, r *http.Request, repository RepositoryAdd[E]) {
	var obj E
	decodeEntity(w, r, &obj)
	err := repository.Update(&obj)
	if err != nil {
		serverErrors.DatabaseSendError(err, w)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

/*
Deleta uma entidade na tabela, correspondente ao Id fornecido na URL.
*/
func Delete(w http.ResponseWriter, r *http.Request, repository RepositoryDelete) {
	id := r.PathValue("id")
	err := repository.Delete(id)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
	}
}
