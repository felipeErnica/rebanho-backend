package handlers

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type LactationHandler struct {
	Repository  repositories.LactationRepository
    Impl        HandlerImpl[entity.Lactation]
}

func InitLactation(app *app.App) {
    repository:= repositories.LactationRepository{}
    repository.Init()

    impl:=HandlerImpl[entity.Lactation]{
        Repository: *repository.Base.Base,
    }

    handler:=LactationHandler{
        Repository: repository,
        Impl: impl,
    }
    
    app.HandleFunc("GET /lactation/page", handler.FindPage)
    app.HandleFunc("GET /lactation/animal/{animalId}", handler.FindByCow)
    app.HandleFunc("POST /lactation/", handler.Add)
    app.HandleFunc("POST /lactation/save", handler.Save)
    app.HandleFunc("DELETE /lactation/{id}", handler.Delete)
    LogControllersInit("Lactações")
}

func (h *LactationHandler) FindPage(w http.ResponseWriter, r *http.Request) {
    cursor, sort, direction:= h.Impl.GetPageParameters(r)
    page, err:=h.Repository.FindPage(sort, direction, cursor)
    h.Impl.SendPage(w, page, err)
}

func (h *LactationHandler) FindByCow(w http.ResponseWriter, r *http.Request) {
    animalId:=r.PathValue("animalId")
    animalsList, err:=h.Repository.FindByAnimal(animalId)
    h.Impl.SendList(w, animalsList, err)
}

func (h *LactationHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w, r)
}

func (h *LactationHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w, r)
}

func (h *LactationHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
