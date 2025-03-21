package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type SlaughterGroupHandler struct {
    Impl HandlerImpl[entity.SlaughterGroup]
    Repository repositories.SlaughterGroupRepository
}

func InitSlaughterGroup(mux *http.ServeMux) {
    repository:=new(repositories.SlaughterGroupRepository)
    repository.Init()

    impl:= HandlerImpl[entity.SlaughterGroup]{
        Repository: repository.Impl,
    }

    handler:=SlaughterGroupHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    mux.HandleFunc("GET /slaughter/groups", handler.FindAll)
    mux.HandleFunc("GET /slaughter/groups/{id}", handler.FindById)
    mux.HandleFunc("POST /slaughter/groups", handler.Add)
    mux.HandleFunc("POST /slaughter/groups/save", handler.Save)
    mux.HandleFunc("DELETE /slaughter/groups/{id}", handler.Delete)
    LogControllersInit("Frigoríficos")
}

func (h *SlaughterGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindAll(w,r)
}

func (h *SlaughterGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *SlaughterGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *SlaughterGroupHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *SlaughterGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}

