package handlers

import (
	"encoding/json"
	"fmt"
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
    mux.HandleFunc("POST /animais", handler.Add)
    mux.HandleFunc("POST /animais/save", handler.Save)
    LogControllersInit("Animais")
}

func (h *AnimalHandler) GetPage(w http.ResponseWriter, r *http.Request)  {
    cursor:=r.URL.Query().Get("cursor")
    sort:= r.URL.Query().Get("sort")
    order:= r.URL.Query().Get("order")

    animals, err:= h.Repository.GetPage(sort, order, cursor)

    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animals.GetPage())
    if err != nil {
        JsonServerError(err, w)
        return
    }
    
    w.Header().Set("Content-Type","application/json")
    w.Write(response)
}

func (h *AnimalHandler) GetById(w http.ResponseWriter, r *http.Request)  {
    id:= r.PathValue("id")
    animal, err:= h.Repository.GetById(id)
    if err != nil {
        DatabaseGetError(err, w)
        return
    }

    response, err:= json.Marshal(animal)
    if err != nil {
        JsonServerError(err, w)
        return
    }

    w.Write(response)
}

func (h *AnimalHandler) Add(w http.ResponseWriter, r *http.Request) {
    var createAnimal entity.CreateAnimal
    if err:= json.NewDecoder(r.Body).Decode(&createAnimal); err != nil {
        JsonServerError(err, w)
        return
    }

    var test entity.Animal
    if err:= json.NewDecoder(r.Body).Decode(&test); err != nil {
        JsonServerError(err, w)
        return
    }

    fmt.Println(json.Marshal(test))

    animal, err:= h.Repository.Add(&createAnimal)
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
