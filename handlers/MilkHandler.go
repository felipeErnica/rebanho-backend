package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type MilkHandler struct {
    Repository repositories.MilkRepository
}

func InitMilk(mux *http.ServeMux) {
    repository:= repositories.MilkRepository{}
    repository.Init()
    handler:= MilkHandler{
        Repository: repository,
    }

    mux.HandleFunc("GET /milkEntries/page", handler.FindPage)
    mux.HandleFunc("GET /milkEntries/animal/{animalId}", handler.FindByCow)
    mux.HandleFunc("GET /milkEntries/entryDate/{entryDate}", handler.FindByEntryDate)
    mux.HandleFunc("POST /milkEntries/", handler.Add)
    mux.HandleFunc("POST /milkEntries/save", handler.Save)
    mux.HandleFunc("DELETE /milkEntries/{id}", handler.Delete)
    LogControllersInit("Entradas de Leite")
}

func (h *MilkHandler) FindPage(w http.ResponseWriter, r *http.Request) {
    cursor:=r.URL.Query().Get("cursor")
    sort:=r.URL.Query().Get("sort")
    direction:=r.URL.Query().Get("order")
    page, err:=h.Repository.FindPage(sort, direction, cursor)
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(page)
    if err != nil {
        JsonServerError(err, w)
    }
    w.Write(response)
}

func (h *MilkHandler) FindByCow(w http.ResponseWriter, r *http.Request) {
    animalId:=r.PathValue("animalId")
    animalsList, err:=h.Repository.FindByAnimal(animalId)
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(animalsList)
    if err != nil {
        JsonServerError(err, w)
    }
    w.Write(response)
}

func (h *MilkHandler) FindByEntryDate(w http.ResponseWriter, r *http.Request) {
    dateString:=r.PathValue("entryDate")
    entryDate, err:=time.Parse(time.RFC3339Nano, dateString)
    if err != nil {
        return
    }

    animalsList, err:=h.Repository.FindByEntryDate(entryDate)
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(animalsList)
    if err != nil {
        JsonServerError(err, w)
    }
    w.Write(response)
}

func (h *MilkHandler) Add(w http.ResponseWriter, r *http.Request) {
    var newMilk entity.MilkEntry;
    err:= json.NewDecoder(r.Body).Decode(&newMilk)
    if err != nil {
        JsonServerError(err, w)
    }

    lactation, err:=h.Repository.Add(newMilk)
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

func (h *MilkHandler) Save(w http.ResponseWriter, r *http.Request) {
    var milk entity.MilkEntry;
    err:= json.NewDecoder(r.Body).Decode(&milk)
    if err != nil {
        JsonServerError(err, w)
    }

    err = h.Repository.Save(&milk)
    if err != nil {
        DatabaseSendError(err, w)
    }
}

func (h *MilkHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id:=r.PathValue("id")
    err:=h.Repository.Delete(id)
    if err != nil {
        DatabaseSendError(err, w)
    }
}
