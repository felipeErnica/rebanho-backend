package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type PastureHandler struct {
    Impl        HandlerImpl[entity.Pasture]
    Repository  repositories.PastureRepository
}

func InitPastureHandler(app *app.App) {
    repository:= new(repositories.PastureRepository)
    repository.Init()
    impl:=HandlerImpl[entity.Pasture]{
        Repository: repository.Impl,
    }
    handler:=PastureHandler{
        Impl: impl,
        Repository: *repository,
    }

    app.HandleFunc("GET /pastures", handler.FindAll)
    app.HandleFunc("GET /pastures/{id}", handler.FindById)
    app.HandleFunc("POST /pastures", handler.Add)
    app.HandleFunc("POST /pastures/save", handler.Save)
    app.HandleFunc("DELETE /pastures/{id}", handler.Delete)
    LogControllersInit("Pastos")
}

func (h *PastureHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindAll(w,r)
}

func (h *PastureHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *PastureHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *PastureHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *PastureHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
