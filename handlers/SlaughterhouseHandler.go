package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type SlaughterhouseHandler struct {
    Impl        HandlerImpl[entity.Slaughterhouse]
    Repository  repositories.SlaughterhouseRepository  
}

func InitSlaugherhouses(mux *http.ServeMux) {
    repository:=new(repositories.SlaughterhouseRepository)
    repository.Init()

    impl:= HandlerImpl[entity.Slaughterhouse]{
        Repository: repository.Impl,
    }

    handler:=SlaughterhouseHandler{ 
        Repository: *repository,
        Impl: impl,
    }

    mux.HandleFunc("GET /slaughterhouses", handler.FindAll)
    mux.HandleFunc("GET /slaughterhouses/{id}", handler.FindById)
    mux.HandleFunc("POST /slaughterhouses", handler.Add)
    mux.HandleFunc("POST /slaughterhouses/save", handler.Save)
    mux.HandleFunc("DELETE /slaughterhouses/{id}", handler.Delete)
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
