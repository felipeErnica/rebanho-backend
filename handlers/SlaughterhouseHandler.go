package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type SlaughterhouseHandler struct {
    Impl        HandlerImpl[entity.Slaughterhouse]
    Repository  repositories.SlaughterhouseRepository  
}

func InitSlaugherhouses(app *app.App) {
    repository:=new(repositories.SlaughterhouseRepository)
    repository.Init()

    impl:= HandlerImpl[entity.Slaughterhouse]{
        Repository: repository.Impl,
    }

    handler:=SlaughterhouseHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    app.HandleFunc("GET /slaughter/slaughterhouses", handler.FindAll)
    app.HandleFunc("GET /slaughter/slaughterhouses/{id}", handler.FindById)
    app.HandleFunc("POST /slaughter/slaughterhouses", handler.Add)
    app.HandleFunc("POST /slaughter/slaughterhouses/save", handler.Save)
    app.HandleFunc("DELETE /slaughterhouses/{id}", handler.Delete)
    LogControllersInit("Frigoríficos")
}

func (h *SlaughterhouseHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindAll(w,r)
}

func (h *SlaughterhouseHandler) FindById(w http.ResponseWriter, r *http.Request) {
    h.Impl.FindById(w,r)
}

func (h *SlaughterhouseHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *SlaughterhouseHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *SlaughterhouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
