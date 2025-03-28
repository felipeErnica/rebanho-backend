package handlers

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

type MilkHandler struct {
    Impl        HandlerImpl[entity.MilkEntry]
    Repository  repositories.MilkRepository
}

func InitMilk(app *app.App) {
    repository:= repositories.MilkRepository{}
    repository.Init()

    impl:=HandlerImpl[entity.MilkEntry]{
        Repository: *repository.Impl.Base,
    }

    handler:= MilkHandler{
        Impl: impl,
        Repository: repository,
    }

    app.HandleFunc("GET /milkEntries/page", handler.FindPage)
    app.HandleFunc("GET /milkEntries/animal/{animalId}", handler.FindByCow)
    app.HandleFunc("GET /milkEntries/entryDate/{entryDate}", handler.FindByEntryDate)
    app.HandleFunc("POST /milkEntries/", handler.Add)
    app.HandleFunc("POST /milkEntries/save", handler.Save)
    app.HandleFunc("DELETE /milkEntries/{id}", handler.Delete)
    LogControllersInit("Entradas de Leite")
}

func (h *MilkHandler) FindPage(w http.ResponseWriter, r *http.Request) {
    cursor, sort, direction:=h.Impl.GetPageParameters(r)
    page, err:=h.Repository.FindPage(sort, direction, cursor)
    h.Impl.SendPage(w, page, err)
}

func (h *MilkHandler) FindByCow(w http.ResponseWriter, r *http.Request) {
    animalId:=r.PathValue("animalId")
    milkList, err:=h.Repository.FindByAnimal(animalId)
    h.Impl.SendList(w, milkList, err)
}

func (h *MilkHandler) FindByEntryDate(w http.ResponseWriter, r *http.Request) {
    dateString:=r.PathValue("entryDate")
    entryDate, err:=time.Parse(time.RFC3339Nano, dateString)
    if err != nil {
        return
    }
    milkList, err:=h.Repository.FindByEntryDate(entryDate)
    h.Impl.SendList(w, milkList, err)
}

func (h *MilkHandler) Add(w http.ResponseWriter, r *http.Request) {
    h.Impl.Add(w,r)
}

func (h *MilkHandler) Save(w http.ResponseWriter, r *http.Request) {
    h.Impl.Save(w,r)
}

func (h *MilkHandler) Delete(w http.ResponseWriter, r *http.Request) {
    h.Impl.Delete(w,r)
}
