package handlersUtil

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	authConfig "github.com/felipeErnica/rebanho-backend/config/auth-config"
)

/*
Decodifica a entide contida no corpo da solicitação HTTP
e retorna um erro caso o formato esteja incorreto.
*/
func DecodeEntity[E any](w http.ResponseWriter, r *http.Request, entity *E) (*E, bool) {
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		err = fmt.Errorf("Falha na decodificação da entidade: %s", err.Error())
		apiError.JsonServerError(err, w)
		return nil, false
	}
	return entity, true
}

/*
Retorna o ID do usuário estocado no Contexto da Requisição pelo middleware de autenticação.
Em caso de erro, retorna um booleano para o cancelamento da funcão requisitante.
*/
func GetUserId(w http.ResponseWriter, r *http.Request) (string, bool) {
	userId, err := authConfig.GetUserId(r)
	if err != nil {
		apiError.DatabaseGetError(err, w)
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
		apiError.JsonServerError(err, w)
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
		apiError.JsonServerError(err, w)
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
		err = fmt.Errorf("Falha na decodificação do filtro: %s", err.Error())
		apiError.JsonServerError(err, w)
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
	searchById func(string, []string) (*[]E, error),
	searchByInput func(string, string) (*[]E, error),
) {
	userId, ok := GetUserId(w, r)
	if !ok {
		return
	}

	id := r.URL.Query().Get("id")
	if id != "" {
		idList := ParseArray(id)
		list, err := searchById(userId, idList)
		if err != nil {
			apiError.DatabaseGetError(err, w)
			return
		}
		SendList(w, list)
		return
	}

	input := "%" + r.URL.Query().Get("input") + "%"
	list, err := searchByInput(userId, input)
	if err != nil {
		apiError.DatabaseGetError(err, w)
		return
	}
	SendList(w, list)
}

func ReturnSearchResultById[E any](
	id string,
	w http.ResponseWriter,
	r *http.Request,
	search func(string, []string) (*[]E, error),
) {
}

/*
Retorna uma entidade como resposta HTTP pelo id fornecido, o repositório deve conter uma função FindById.
*/
func FindById[E any](w http.ResponseWriter, r *http.Request, repository RepositoryFindById[E]) {
	id := r.PathValue("id")
	entity, err := repository.FindById(id)
	if err != nil {
		apiError.DatabaseGetError(err, w)
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
		apiError.DatabaseGetError(err, w)
		return
	}
	SendList(w, list)
}

/*
Atualiza uma entidade no banco, e retorna o resultado na Resposta HTTP
*/
func Update(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(apiError.APIError{Title: "Sucesso", Message: "Atualizado com sucesso!"})
	if err != nil {
		apiError.JsonServerError(err, w)
		return
	}

	writeResponse(w, response)
}


/*
Criado uma entidade no banco, e retorna o resultado na Resposta HTTP
*/
func Add(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)

	response, err := json.Marshal(apiError.APIError{Title: "Sucesso", Message: "Adicionado com sucesso!"})
	if err != nil {
		apiError.JsonServerError(err, w)
		return
	}

	writeResponse(w, response)
}

/*
Deleta uma entidade na tabela, correspondente ao Id fornecido na URL.
*/
func Delete(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(apiError.APIError{Title: "Sucesso", Message: "Excluído com sucesso!"})
	if err != nil {
		apiError.JsonServerError(err, w)
		return
	}

	writeResponse(w, response)
}
