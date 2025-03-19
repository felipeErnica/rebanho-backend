package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type PregnancyLossHandler struct {
    Impl       HandlerImpl[entity.PregnancyLoss]
    Repository repositories.PregnancyLossRepository
}

func InitPregnancyLossHandler(mux *http.ServeMux) {
    repository:= new(repositories.PregnancyLossRepository)
    repository.Init()

    impl:=HandlerImpl[entity.PregnancyLoss]{
        Repository: *repository.Impl.Base,
    }

    handler:=PregnancyLossHandler{
        Impl: impl,
        Repository: *repository,
    }

    mux.HandleFunc("GET /losses/page", handler.FindPage)
    mux.HandleFunc("GET /losses/animal/{animalId}", handler.FindByAnimalId)
    mux.HandleFunc("GET /losses/{id}", handler.FindById)
    mux.HandleFunc("POST /losses", handler.Add)
    mux.HandleFunc("POST /losses/save", handler.Save)
    mux.HandleFunc("DELETE /losses/{id}", handler.Delete)
    LogControllersInit("Perdas de Parição")
}

func (h *PregnancyLossHandler) FindPage(w http.ResponseWriter, r *http.Request) {
    cursor, sort, order:= h.Impl.GetPageParameters(r)
    page, err:= h.Repository.FindPage(cursor, sort, order)
    h.Impl.SendPage(w, page, err)
}

func (h *PregnancyLossHandler) FindByAnimalId(w http.ResponseWriter, r *http.Request) {
    animalId:=r.PathValue("animalId")
    list, err:= h.Repository.FindByAnimalId(animalId)
    h.Impl.SendList(w, list, err)
}

func (h *PregnancyLossHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *PregnancyLossHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *PregnancyLossHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *PregnancyLossHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
