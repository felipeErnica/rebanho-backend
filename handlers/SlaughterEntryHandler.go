package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type SlaughterEntryHandler struct {
    Impl HandlerImpl[entity.SlaughterEntry]
    Repository repositories.SlaughterEntryRepository
}

func InitSlaughterEntry(mux *http.ServeMux) {
    repository:=new(repositories.SlaughterEntryRepository)
    repository.Init()

    impl:= HandlerImpl[entity.SlaughterEntry]{
        Repository: repository.Impl,
    }

    handler:=SlaughterEntryHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    mux.HandleFunc("GET /slaughter/groups/{groupId}/entries", handler.FindByGroupId)
    mux.HandleFunc("GET /slaughter/entries/{id}", handler.FindById)
    mux.HandleFunc("POST /slaughter/entries", handler.Add)
    mux.HandleFunc("POST /slaughter/entries/save", handler.Save)
    mux.HandleFunc("DELETE /slaughter/entries/{id}", handler.Delete)
    LogControllersInit("Entradas de Abate")
}

func (h *SlaughterEntryHandler) FindByGroupId(w http.ResponseWriter, r *http.Request) {
    groupId:=r.PathValue("groupId")
    list, err:=h.Repository.FindByGroupId(groupId)
    h.Impl.SendList(w, list, err)
} 

func (h *SlaughterEntryHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *SlaughterEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *SlaughterEntryHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *SlaughterEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
