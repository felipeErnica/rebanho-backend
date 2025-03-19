package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type PastureEntryHandler struct {
    Impl        HandlerImpl[entity.PastureEntry]
    Repository  repositories.PastureEntryRepository
}

func InitPastureEntryHandler(mux *http.ServeMux) {
    repository:= new(repositories.PastureEntryRepository)
    repository.Init()
    impl:=HandlerImpl[entity.PastureEntry]{
        Repository: *repository.Base.Base,
    }
    handler:=PastureEntryHandler{
        Impl: impl,
        Repository: *repository,
    }

    mux.HandleFunc("GET /pastureEntries/pasture/{pastureId}/page", handler.FindByPastureId)
    mux.HandleFunc("GET /pastureEntries/pasture/{pastureId}/deleted/page", handler.FindDeletedByPastureId)
    mux.HandleFunc("GET /pastureEntries/animal/{animalId}", handler.FindByAnimalId)
    mux.HandleFunc("POST /pastureEntries/", handler.Add)
    mux.HandleFunc("POST /pastureEntries/save", handler.Save)
    mux.HandleFunc("DELETE /pastureEntries/{id}", handler.Delete)
    LogControllersInit("Entradas no Lote")
}

func (h *PastureEntryHandler) FindByPastureId(w http.ResponseWriter, r *http.Request) {
    cursor, sort, order:= h.Impl.GetPageParameters(r)
    pastureId:=r.PathValue("pastureId")
    page, err:=h.Repository.FindByPastureId(cursor, sort, order, pastureId)
    h.Impl.SendPage(w, page, err)
}

func (h *PastureEntryHandler) FindDeletedByPastureId(w http.ResponseWriter, r *http.Request) {
    cursor, sort, order:= h.Impl.GetPageParameters(r)
    pastureId:=r.PathValue("pastureId")
    page, err:=h.Repository.FindByDeletedPasturePage(cursor, sort, order, pastureId)
    h.Impl.SendPage(w, page, err)
}

func (h *PastureEntryHandler) FindByAnimalId(w http.ResponseWriter, r *http.Request) {
    animalId:=r.PathValue("animalId")
    list, err:=h.Repository.FindByAnimalId(animalId)
    h.Impl.SendList(w, list, err)
}

func (h *PastureEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *PastureEntryHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *PastureEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
