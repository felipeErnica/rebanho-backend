package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type InseminationGroupHandler struct {
	Impl HandlerImpl[entity.InseminationGroup]
    Repository repositories.InseminationGroupRepository
}

func InitInseminationGroup(mux *http.ServeMux) {
    repository:=new(repositories.InseminationGroupRepository)
    repository.Init()
    
    impl:=HandlerImpl[entity.InseminationGroup]{
        Repository: repository.Impl,
    }

    handler:=InseminationGroupHandler{
        Impl: impl,
        Repository: *repository,
    }

    mux.HandleFunc("GET insemination/groups", handler.FindAll)
    mux.HandleFunc("GET insemination/groups/{id}", handler.FindById)
    mux.HandleFunc("POST insemination/groups", handler.Add)
    mux.HandleFunc("POST insemination/groups/save", handler.Save)
    mux.HandleFunc("DELETE insemination/groups/{id}", handler.Delete)
    LogControllersInit("Grupos de Inseminação")
}

func (h *InseminationGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindAll(w,r)
}

func (h *InseminationGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *InseminationGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *InseminationGroupHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *InseminationGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
