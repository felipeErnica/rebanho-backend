package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type WeightEntryHandler struct {
    Impl        HandlerImpl[entity.WeightEntry]
    Repository  repositories.WeightEntryRepository
}

func InitWeightEntries(mux *http.ServeMux) {
    repository:=new(repositories.WeightEntryRepository)
    repository.Init()

    impl:= HandlerImpl[entity.WeightEntry]{
        Repository: repository.Impl,
    }

    handler:=WeightEntryHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    mux.HandleFunc("GET /weight/group/{groupId}/entries", handler.FindByGroupId)
    mux.HandleFunc("GET /weight/entries/{id}", handler.FindById)
    mux.HandleFunc("POST /weight/entries", handler.Add)
    mux.HandleFunc("POST /weight/entries/save", handler.Save)
    mux.HandleFunc("DELETE /weight/entries/{id}", handler.Delete)
    LogControllersInit("Entradas de Peso")
}

func (h *WeightEntryHandler) FindByGroupId(w http.ResponseWriter, r *http.Request) {
    groupId:=r.PathValue("groupId")
    list, err:= h.Repository.FindByGroupId(groupId)
    h.Impl.SendList(w, list, err)
}

func (h *WeightEntryHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *WeightEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *WeightEntryHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *WeightEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
