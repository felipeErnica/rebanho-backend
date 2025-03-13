package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type AnimalHandler struct {
    Repository repositories.AnimalRepository
}

func InitAnimal(mux *http.ServeMux) {

    repository:=new(repositories.AnimalRepository)
    repository.Init()

    handler:=AnimalHandler{ 
        Repository: *repository,
    }

    mux.HandleFunc("GET /animais/page", handler.GetPage)
    mux.HandleFunc("GET /animais/{id}", handler.GetById)
    mux.HandleFunc("GET /animais/father/{fatherId}", handler.FindByFatherId)
    mux.HandleFunc("GET /animais/mother/{motherId}", handler.FindByMotherId)
    mux.HandleFunc("GET /animais/pasture/{pastureId}/page", handler.FindByPastureId)
    mux.HandleFunc("GET /animais/deleted/page", handler.FindDeleted)
    mux.HandleFunc("POST /animais", handler.Add)
    mux.HandleFunc("POST /animais/save", handler.Save)
    mux.HandleFunc("DELETE /animais/{id}", handler.Delete)
    LogControllersInit("Animais")
}

func (h *AnimalHandler) GetPage(w http.ResponseWriter, r *http.Request)  {
    cursor:=r.URL.Query().Get("cursor")
    sort:= r.URL.Query().Get("sort")
    order:= r.URL.Query().Get("order")

    animals, err:= h.Repository.FindPage(sort, order, cursor)

    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animals)
    if err != nil {
        JsonServerError(err, w)
        return
    }
    writeResponse(w, response)   
}

func (h *AnimalHandler) GetById(w http.ResponseWriter, r *http.Request)  {
    id:= r.PathValue("id")
    animal, err:= h.Repository.FindById(id)
    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animal)
    if err != nil {
        JsonServerError(err, w)
        return
    }

    writeResponse(w, response)
}

func (h *AnimalHandler) FindByFatherId(w http.ResponseWriter, r *http.Request)  {
    fatherId:= r.PathValue("fatherId")
    animal, err:= h.Repository.FindByFatherId(fatherId)
    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animal)
    if err != nil {
        JsonServerError(err, w)
        return
    }

    writeResponse(w, response)
}

func (h *AnimalHandler) FindByMotherId(w http.ResponseWriter, r *http.Request)  {
    motherId:= r.PathValue("motherId")
    animal, err:= h.Repository.FindById(motherId)
    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animal)
    if err != nil {
        JsonServerError(err, w)
        return
    }

    writeResponse(w, response)
}

func (h *AnimalHandler) FindByPastureId(w http.ResponseWriter, r *http.Request)  {
    cursor:=r.URL.Query().Get("cursor")
    sort:= r.URL.Query().Get("sort")
    order:= r.URL.Query().Get("order")
    pastureId:=r.PathValue("pastureId")

    animals, err:= h.Repository.FindByPastureId(sort, order, cursor, pastureId)

    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animals)
    if err != nil {
        JsonServerError(err, w)
        return
    }
    writeResponse(w, response)   
}

func (h *AnimalHandler) FindDeleted(w http.ResponseWriter, r *http.Request)  {
    cursor:=r.URL.Query().Get("cursor")
    sort:= r.URL.Query().Get("sort")
    order:= r.URL.Query().Get("order")

    animals, err:= h.Repository.FindDeletedPage(sort, order, cursor)

    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animals)
    if err != nil {
        JsonServerError(err, w)
        return
    }
    writeResponse(w, response)   
}

func (h *AnimalHandler) Add(w http.ResponseWriter, r *http.Request) {
    var create entity.Animal
    if err:= json.NewDecoder(r.Body).Decode(&create); err != nil {
        JsonServerError(err, w)
        return
    }

    animal, err:= h.Repository.Add(create)
    if err != nil {
        DatabaseSendError(err, w)
        return
    }

    response, err:= json.Marshal(animal)
    if err != nil {
        JsonServerError(err,w)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write(response)
}

func (h *AnimalHandler) Save(w http.ResponseWriter, r *http.Request) {
    var animal entity.Animal
    if err:= json.NewDecoder(r.Body).Decode(&animal); err != nil {
        JsonServerError(err, w)
        return
    }

    err:= h.Repository.Save(&animal)
    if err != nil {
        DatabaseSendError(err, w)
        return
    }

}

func (h *AnimalHandler) Delete(w http.ResponseWriter, r *http.Request) {
    id:= r.PathValue("id")
    err:=h.Repository.Delete(id)
    if err != nil {
        DatabaseSendError(err, w)
    }
}
