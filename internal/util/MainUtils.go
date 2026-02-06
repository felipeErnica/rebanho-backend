package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	authConfig "github.com/felipeErnica/rebanho-backend/internal/config/auth-config"
	"github.com/gorilla/schema"
)

/*
Decodifica a entide contida no corpo da solicitação HTTP
e retorna um erro caso o formato esteja incorreto.
*/
func DecodeEntity[E any](w http.ResponseWriter, r *http.Request, entity *E) (*E, bool) {
	err := json.NewDecoder(r.Body).Decode(&entity)
	if err != nil {
		err = fmt.Errorf("Falha na decodificação da entidade: %s", err.Error())
		log.JsonServerError(err, w)
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
		log.WriteError(w, err)
		return userId, false
	}

	return userId, true
}

/*
Envia a entidade como resposta HTTP.
*/
func WriteEntity[E any](w http.ResponseWriter, model *E) {
	response, err := json.Marshal(model)
	if err != nil {
		log.JsonServerError(err, w)
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
		log.JsonServerError(err, w)
		return
	}
	writeResponse(w, response)
}

/*
Decodifica o filtro contido na solicitação HTTP
e retorna um erro caso o formato esteja incorreto. Como parâmetro,
utiliza um filtro vazio do mesmo tipo do filtro a ser retornado.
*/
func DecodeFilter[F any](r *http.Request, filter F) (*F, error) {
	urlVals := r.URL.Query()
	urlVals.Del("sort")
	urlVals.Del("cursor")
	urlVals.Del("order")

	if len(urlVals) <= 0 {
		return nil, nil
	}

	for key := range urlVals {
		values := strings.Split(urlVals.Get(key), ",")
		if len(values) > 1 {
			urlVals[key] = values
		}
	}

	err := schema.NewDecoder().Decode(&filter, urlVals)
	if err != nil {
		return nil, err
	}

	return &filter, nil
}

/*
Decodifica o objeto contido na URL da solicitação HTTP
e retorna um erro caso o formato esteja incorreto. Como parâmetro,
utiliza um objeto vazio do mesmo tipo do objeto a ser retornado.
*/
func DecodeURL[F any](r *http.Request, obj F) (*F, error) {
	urlVals := r.URL.Query()
	if len(urlVals) <= 0 {
		return nil, nil
	}

	for key := range urlVals {
		values := strings.Split(urlVals.Get(key), ",")
		if len(values) > 1 {
			urlVals[key] = values
		}
	}

	err := schema.NewDecoder().Decode(&obj, urlVals)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

func WriteUpdateResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(log.APIError{Title: "Sucesso", Message: "Atualizado com sucesso!"})
	if err != nil {
		log.JsonServerError(err, w)
		return
	}

	writeResponse(w, response)
}

/*
Criado uma entidade no banco, e retorna o resultado na Resposta HTTP
*/
func WriteCreatedResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)

	response, err := json.Marshal(log.APIError{Title: "Sucesso", Message: "Adicionado com sucesso!"})
	if err != nil {
		log.JsonServerError(err, w)
		return
	}

	writeResponse(w, response)
}

/*
Deleta uma entidade na tabela, correspondente ao Id fornecido na URL.
*/
func WriteDeleteResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)

	response, err := json.Marshal(log.APIError{Title: "Sucesso", Message: "Excluído com sucesso!"})
	if err != nil {
		log.JsonServerError(err, w)
		return
	}

	writeResponse(w, response)
}
