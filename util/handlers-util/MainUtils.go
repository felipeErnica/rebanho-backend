package handlersUtil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	authConfig "github.com/felipeErnica/rebanho-backend/config/auth-config"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
	"github.com/google/uuid"
)

/*
Decodifica a entide contida no corpo da solicitação HTTP
e retorna um erro caso o formato esteja incorreto.
*/
func DecodeEntity[E any](w http.ResponseWriter, r *http.Request, entity *E) bool {
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		err = errors.New(fmt.Sprintf("Falha na decodificação da entidade: %s", err.Error()))
		serverErrors.JsonServerError(err, w)
		return false
	}
	return true
}

/*
Retorna o ID do usuário estocado no Contexto da Requisição pelo middleware de autenticação.
Em caso de erro, retorna um booleano para o cancelamento da funcão requisitante.
*/
func GetUserId(w http.ResponseWriter, r *http.Request) (string, bool) {
	userId, err := authConfig.GetUserId(r)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return userId, false
	}

	return userId, true
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
e retorna um erro caso o formato esteja incorreto. Como parâmetro,
utiliza um filtro vazio do mesmo tipo do filtro a ser retornado.
*/
func DecodeFilter[F any](w http.ResponseWriter, r *http.Request, filter F) (F, bool) {
	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		err = errors.New(fmt.Sprintf("Falha na decodificação do filtro: %s", err.Error()))
		serverErrors.JsonServerError(err, w)
		return filter, false
	}
	return filter, true
}

/*
Retorna uma lista como resposta HTTP ao texto a ser pesquisado, o repositório deve conter uma função Search.
*/
func ReturnSearchResults[E any](
	w http.ResponseWriter,
	r *http.Request,
	search func(string, string) (*[]E, error),
) {
	input := "%" + r.URL.Query().Get("input") + "%"
	userId, ok := GetUserId(w, r)
	if !ok {
		return
	}
	list, err := search(userId, input)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	SendList(w, list)
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
	userId, ok := GetUserId(w, r)
	if !ok {
		return
	}
	list, err := repository.FindAll(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	SendList(w, list)
}

/*
Preenche os campos da entidade criada com os valores respectivos:
Id - Gera uma nova chave aleatória UUID
CreatedAt - Preenche o momento da criação
UserId - Recupera o ID do usuário salvo no Contexto da requisição HTTP
*/
func fillCreationFields[E any](w http.ResponseWriter, r *http.Request, obj *E) bool {
	value := reflect.ValueOf(obj)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}

	fieldId := value.FieldByName("Id")
	fieldUserId := value.FieldByName("UserId")
	fieldCreatedAt := value.FieldByName("CreatedAt")

	if !fieldId.CanSet() || !fieldUserId.CanSet() || !fieldCreatedAt.CanSet() {
		err := errors.New("Formato de estrura não suporta adições!")
		serverErrors.DatabaseSendError(err, w)
		return false
	}

	id := uuid.NewString()
	userId, ok := GetUserId(w, r)
	if !ok {
		return ok
	}
	createdAt := reflect.ValueOf(time.Now())

	fieldId.SetString(id)
	fieldCreatedAt.Set(createdAt)
	fieldUserId.SetString(userId)
	return true
}

/*
Salva uma nova entidade no banco, e retorna o resultado na Resposta HTTP
*/
func Add[E any](w http.ResponseWriter, r *http.Request, repository RepositoryAdd[E]) {
	var obj E
	DecodeEntity(w, r, &obj)
	ok := fillCreationFields(w, r, &obj)
	if !ok {
		return
	}

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
func Update[E any](w http.ResponseWriter, r *http.Request, repository RepositoryUpdate[E]) {
	var obj E
	DecodeEntity(w, r, &obj)
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
