package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type InseminationEntryHandler struct {
	Impl HandlerImpl[entity.InseminationEntry]
    Repository repositories.InseminationEntryRepository
}

func InitInseminationEntry(app *app.App) {
    repository:=new(repositories.InseminationEntryRepository)
    repository.Init()
    impl:=HandlerImpl[entity.InseminationEntry] {
        Repository: repository.Impl,
    }
    handler:=InseminationEntryHandler{
        Repository: *repository,
        Impl: impl,
    }

    app.HandleFunc("GET insemination/groups/{groupId}/entries", handler.FindByGroupId)
    app.HandleFunc("GET insemination/entries/{id}", handler.FindById)
    app.HandleFunc("POST insemination/entries", handler.Add)
    app.HandleFunc("POST insemination/entries/save", handler.Save)
    app.HandleFunc("DELETE insemination/entries/{id}", handler.Delete)
    LogControllersInit("Entradas de Inseminação")
}

func (h *InseminationEntryHandler) FindByGroupId(w http.ResponseWriter, r *http.Request) {
    groupId:=r.PathValue("groupId")
    list, err:= h.Repository.FindByGroupId(groupId)
    h.Impl.SendList(w, list, err)
}

func (h *InseminationEntryHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *InseminationEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *InseminationEntryHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *InseminationEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
