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

    mux.HandleFunc("GET /weightEntries/group/{groupId}", handler.FindByGroupId)
    mux.HandleFunc("GET /weightEntries/{id}", handler.FindById)
    mux.HandleFunc("POST /weightEntries", handler.Add)
    mux.HandleFunc("POST /weightEntries/save", handler.Save)
    mux.HandleFunc("DELETE /weightEntries/{id}", handler.Delete)
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
