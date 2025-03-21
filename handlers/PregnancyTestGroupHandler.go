package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type PregnancyTestGroupHandler struct {
    Impl HandlerImpl[entity.PregancyTestGroup]
    Repository repositories.PregnancyTestGroupRepository
}

func InitPregnancyTestGroup(mux *http.ServeMux) {
    repository:= new(repositories.PregnancyTestGroupRepository)
    repository.Init()

    impl:=HandlerImpl[entity.PregancyTestGroup]{
        Repository: repository.Impl,
    }

    handler:= PregnancyTestGroupHandler{
        Impl: impl,
        Repository: *repository,
    }

    mux.HandleFunc("GET /pregnancy/groups",  handler.FindAll)
    mux.HandleFunc("GET /pregnancy/groups/{id}",  handler.FindById)
    mux.HandleFunc("POST /pregnancy/groups",  handler.Add)
    mux.HandleFunc("POST /pregnancy/groups/save",  handler.Save)
    mux.HandleFunc("DELETE /pregnancy/groups/{id}",  handler.Delete)
    LogControllersInit("Grupo de Toque")
}

func (h *PregnancyTestGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindAll(w,r)
} 

func (h *PregnancyTestGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
} 

func (h *PregnancyTestGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
} 

func (h *PregnancyTestGroupHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
} 

func (h *PregnancyTestGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
} 
