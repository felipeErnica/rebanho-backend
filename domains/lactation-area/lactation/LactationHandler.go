package lactation

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type LactationHandler struct {
	Repository *LactationRepository
}

func (h *LactationHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	filter := LactationFilter{}
	if !handlersUtil.DecodeFilter(w, r, &filter) {
        return
    }
	handlersUtil.ReturnPage(w, r, h.Repository, &filter)
}

func (h *LactationHandler) FindByCow(w http.ResponseWriter, r *http.Request) {
	animalId := r.PathValue("animalId")
	animalsList, err := h.Repository.FindByAnimal(animalId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
	}
	handlersUtil.SendList(w, animalsList)
}

func (h *LactationHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w ,r, h.Repository)
}

func (h *LactationHandler) Save(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Update(w, r, h.Repository)
}

func (h *LactationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Delete(w, r, h.Repository)
}
