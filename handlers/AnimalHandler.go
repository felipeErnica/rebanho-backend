package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type AnimalHandler struct {
    Repository  repositories.AnimalRepository
    Impl        HandlerImpl[entity.Animal]
}

func InitAnimal(app *app.App) {
    repository:=new(repositories.AnimalRepository)
    repository.Init()

    impl:= HandlerImpl[entity.Animal]{
        Repository: *repository.Impl.Base,
    }

    handler:=AnimalHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    app.HandleFunc("POST /animals/page", handler.FindPage)
    app.HandleFunc("GET /animals/{id}", handler.FindById)
    app.HandleFunc("GET /animals/name/{name}", handler.FindByName)
    app.HandleFunc("GET /animals/number/{number}", handler.FindByNumber)
    app.HandleFunc("GET /animals/father/{fatherId}", handler.FindByFatherId)
    app.HandleFunc("GET /animals/mother/{motherId}", handler.FindByMotherId)
    app.HandleFunc("GET /animals/pasture/{pastureId}/page", handler.FindByPastureId)
    app.HandleFunc("POST /animals", handler.Add)
    app.HandleFunc("POST /animals/save", handler.Save)
    app.HandleFunc("DELETE /animals/{id}", handler.Delete)
    LogControllersInit("Animais")
}

func (h *AnimalHandler) FindPage(w http.ResponseWriter, r *http.Request)  {
    var filter entity.AnimalFilter
    
    err := json.NewDecoder(r.Body).Decode(&filter)
    if err != nil {
        err = errors.New(fmt.Sprintf("Falha na decodificação do filtro: %s", err.Error()))
        JsonServerError(err, w)
        return
    }

    cursor, sort, order:= h.Impl.GetPageParameters(r)
    page, err:= h.Repository.FindPage(sort, order, cursor, &filter)
    h.Impl.SendPage(w, page, err)
}

func (h *AnimalHandler) FindById(w http.ResponseWriter, r *http.Request)  {
    h.Impl.FindById(w,r)
}

func (h *AnimalHandler) FindByName(w http.ResponseWriter, r *http.Request)  {
    name:= r.PathValue("name")
    animals, err:= h.Repository.FindByName(name)
    h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByNumber(w http.ResponseWriter, r *http.Request)  {
    number:= r.PathValue("number")
    animals, err:= h.Repository.FindByIdentificationNumber(number)
    h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByFatherId(w http.ResponseWriter, r *http.Request)  {
    fatherId:= r.PathValue("fatherId")
    animals, err:= h.Repository.FindByFatherId(fatherId)
    h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByMotherId(w http.ResponseWriter, r *http.Request)  {
    motherId:= r.PathValue("motherId")
    animals, err:= h.Repository.FindByMotherId(motherId)
    h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByPastureId(w http.ResponseWriter, r *http.Request)  {
    pastureId:=r.PathValue("pastureId")
    cursor, sort, order:= h.Impl.GetPageParameters(r)
    pasturePage, err:= h.Repository.FindByPastureId(sort, order, cursor, pastureId)
    h.Impl.SendPage(w, pasturePage, err)
}

func (h *AnimalHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w, r)
}

func (h *AnimalHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w, r)
}

func (h *AnimalHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
