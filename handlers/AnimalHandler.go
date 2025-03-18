package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type AnimalHandler struct {
    Repository  repositories.AnimalRepository
    Impl        HandlerImpl[entity.Animal]
}

func InitAnimal(mux *http.ServeMux) {

    repository:=new(repositories.AnimalRepository)
    repository.Init()

    impl:= HandlerImpl[entity.Animal]{
        Repository: *repository.Base.Base,
    }

    handler:=AnimalHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    mux.HandleFunc("GET /animals/page", handler.FindPage)
    mux.HandleFunc("GET /animals/{id}", handler.FindById)
    mux.HandleFunc("GET /animals/name/{name}", handler.FindByName)
    mux.HandleFunc("GET /animals/number/{number}", handler.FindByNumber)
    mux.HandleFunc("GET /animals/father/{fatherId}", handler.FindByFatherId)
    mux.HandleFunc("GET /animals/mother/{motherId}", handler.FindByMotherId)
    mux.HandleFunc("GET /animals/pasture/{pastureId}/page", handler.FindByPastureId)
    mux.HandleFunc("POST /animals", handler.Add)
    mux.HandleFunc("POST /animals/save", handler.Save)
    mux.HandleFunc("DELETE /animals/{id}", handler.Delete)
    LogControllersInit("Animais")
}

func (h *AnimalHandler) FindPage(w http.ResponseWriter, r *http.Request)  {
    cursor, sort, order:= h.Impl.GetPageParameters(r)
    page, err:= h.Repository.FindPage(sort, order, cursor)
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
