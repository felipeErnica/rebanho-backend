package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

func InitLactation(mux *http.ServeMux) {
    repository:= repositories.LactationRepository{}
    repository.Init()
    handler:=LactationHandler{
        Repository: repository,
    }
    
    mux.HandleFunc("GET /lactation/page", handler.FindPage)
    mux.HandleFunc("GET /lactation/animal/{animalId}", handler.FindByCow)
    mux.HandleFunc("POST /lactation/", handler.Add)
    mux.HandleFunc("POST /lactation/save", handler.Save)
    mux.HandleFunc("DELETE /lactation/{id}", handler.Delete)
    LogControllersInit("Lactações")
}

type LactationHandler struct {
	Repository repositories.LactationRepository
}

func (l *LactationHandler) FindPage(w http.ResponseWriter, r *http.Request) {
    cursor:=r.URL.Query().Get("cursor")
    sort:=r.URL.Query().Get("sort")
    direction:=r.URL.Query().Get("order")
    page, err:=l.Repository.FindPage(sort, direction, cursor)
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(page)
    if err != nil {
        JsonServerError(err, w)
    }
    w.Write(response)
}

func (l *LactationHandler) FindByCow(w http.ResponseWriter, r *http.Request) {
    animalId:=r.PathValue("animalId")
    animalsList, err:=l.Repository.FindByAnimal(animalId)
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(animalsList)
    if err != nil {
        JsonServerError(err, w)
    }
    w.Write(response)
}

func (l *LactationHandler) Add(w http.ResponseWriter, r *http.Request) {
    var newLactation entity.Lactation;
    err:= json.NewDecoder(r.Body).Decode(&newLactation)
    if err != nil {
        JsonServerError(err, w)
    }

    lactation, err:=l.Repository.Add(newLactation)
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(lactation)
    if err != nil {
        JsonServerError(err, w)
    }
    w.WriteHeader(http.StatusCreated)
    w.Write(response)
}

func (l *LactationHandler) Save(w http.ResponseWriter, r *http.Request) {
    var lactation entity.Lactation;
    err:= json.NewDecoder(r.Body).Decode(&lactation)
    if err != nil {
        JsonServerError(err, w)
    }

    err = l.Repository.Save(&lactation)
    if err != nil {
        DatabaseSendError(err, w)
    }
}

func (l *LactationHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id:=r.PathValue("id")
    err:=l.Repository.Delete(id)
    if err != nil {
        DatabaseSendError(err, w)
    }
}
