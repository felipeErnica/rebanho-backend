package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type WeightGroupHandler struct {
    Impl        HandlerImpl[entity.WeightGroup]
    Repository  repositories.WeightGroupRepository
}

func InitWeightGroups(mux *http.ServeMux) {
    repository:=new(repositories.WeightGroupRepository)
    repository.Init()

    impl:= HandlerImpl[entity.WeightGroup]{
        Repository: repository.Impl,
    }

    handler:=WeightGroupHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    mux.HandleFunc("GET /weightGroup/", handler.FindAll)
    mux.HandleFunc("GET /weightGroup/{id}", handler.FindById)
    mux.HandleFunc("POST /weightGroup", handler.Add)
    mux.HandleFunc("POST /weightGroup/save", handler.Save)
    mux.HandleFunc("DELETE /weightGroup/{id}", handler.Delete)
    LogControllersInit("Grupos de Pesagem")
}

func (h *WeightGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindAll(w,r)
}

func (h *WeightGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *WeightGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *WeightGroupHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *WeightGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
