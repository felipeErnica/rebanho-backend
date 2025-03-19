package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type HandlerImpl[E entity.IEntity] struct {
    Repository  repositories.RepositoryImpl[E]
}

func (h *HandlerImpl[E]) GetPageParameters(r *http.Request) (cursor string, sort string, order string) {
    cursor = r.URL.Query().Get("cursor")
    sort = r.URL.Query().Get("sort")
    order = r.URL.Query().Get("order")
    
    if sort == "" {
        sort = "created_at"
    }

    if order == "" {
        order = "asc"
    }
    return cursor, sort, order
}

func (h *HandlerImpl[E]) SendPage(w http.ResponseWriter, page *entity.Page[E], err error)  {
    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(page)
    if err != nil {
        JsonServerError(err, w)
        return
    }
    writeResponse(w, response)   
}

func (h *HandlerImpl[E]) SendList(w http.ResponseWriter, list *[]E, err error)  {
    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(list)
    if err != nil {
        JsonServerError(err, w)
        return
    }
    writeResponse(w, response)   
}

func (h *HandlerImpl[E]) SendEntity(w http.ResponseWriter, model *E, err error)  {
    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(model)
    if err != nil {
        JsonServerError(err, w)
        return
    }
    writeResponse(w, response)   
}

func (h *HandlerImpl[E]) FindAll(w http.ResponseWriter, r *http.Request) {
    list, err:=h.Repository.FindAll()
    h.SendList(w, list, err)
}

func (h *HandlerImpl[E]) FindById(w http.ResponseWriter, r *http.Request)  {
    id:= r.PathValue("id")
    model, err:= h.Repository.FindById(id)
    h.SendEntity(w, model, err)
}

func (h *HandlerImpl[E]) Add(w http.ResponseWriter, r *http.Request) {
    var create E
    if err:= json.NewDecoder(r.Body).Decode(&create); err != nil {
        JsonServerError(err, w)
        return
    }

    model, err:= h.Repository.Add(create)
    if err != nil {
        DatabaseSendError(err, w)
        return
    }

    response, err:= json.Marshal(model)
    if err != nil {
        JsonServerError(err,w)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write(response)
}

func (h *HandlerImpl[E]) Save(w http.ResponseWriter, r *http.Request) {
    var model E
    if err:= json.NewDecoder(r.Body).Decode(&model); err != nil {
        JsonServerError(err, w)
        return
    }

    err:= h.Repository.Save(&model)
    if err != nil {
        DatabaseSendError(err, w)
        return
    }
}

func (h *HandlerImpl[E]) Delete(w http.ResponseWriter, r *http.Request) {
    id:= r.PathValue("id")
    err:=h.Repository.SoftDelete(id)
    if err != nil {
        DatabaseSendError(err, w)
    }
}
