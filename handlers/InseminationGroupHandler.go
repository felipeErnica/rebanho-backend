package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type InseminationGroupHandler struct {
	Impl HandlerImpl[entity.InseminationGroup]
    Repository repositories.InseminationGroupRepository
}

func InitInseminationGroup(app *app.App) {
    repository:=new(repositories.InseminationGroupRepository)
    repository.Init()
    
    impl:=HandlerImpl[entity.InseminationGroup]{
        Repository: repository.Impl,
    }

    handler:=InseminationGroupHandler{
        Impl: impl,
        Repository: *repository,
    }

    app.HandleFunc("GET insemination/groups", handler.FindAll)
    app.HandleFunc("GET insemination/groups/{id}", handler.FindById)
    app.HandleFunc("POST insemination/groups", handler.Add)
    app.HandleFunc("POST insemination/groups/save", handler.Save)
    app.HandleFunc("DELETE insemination/groups/{id}", handler.Delete)
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
