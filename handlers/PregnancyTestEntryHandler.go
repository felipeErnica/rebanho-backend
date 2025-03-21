package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type PregnancyTestEntryHandler struct {
    Impl HandlerImpl[entity.PregnancyTestEntry]
    Repository repositories.PregancyTestEntryRepository
}

func InitPregnancyTestEntry(mux *http.ServeMux) {
    repository:=new(repositories.PregancyTestEntryRepository)
    repository.Init()

    handler:=new(PregnancyTestEntryHandler)
    handler.Repository = *repository
    handler.Impl = HandlerImpl[entity.PregnancyTestEntry]{
        Repository: repository.Impl,
    }

    mux.HandleFunc("GET /pregnancy/test/groups/{groupId}/entries", handler.FindByGroupId)
    mux.HandleFunc("GET /pregnancy/test/entries/{id}", handler.FindById)
    mux.HandleFunc("POST /pregnancy/test/entries", handler.Add)
    mux.HandleFunc("POST /pregnancy/test/entries/save", handler.Save)
    mux.HandleFunc("DELETE /pregnancy/test/entries/{id}", handler.Delete)
    LogControllersInit("Entradas de Toque")
} 

func (h *PregnancyTestEntryHandler) FindByGroupId(w http.ResponseWriter, r *http.Request) {
    groupId:=r.PathValue("groupId")
    list, err:=h.Repository.FindByGroupId(groupId)
    h.Impl.SendList(w, list, err)
}

func (h *PregnancyTestEntryHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *PregnancyTestEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *PregnancyTestEntryHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *PregnancyTestEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
