package pastureEntries

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type PastureEntryHandler struct {
	Repository *PastureEntryRepository
}

func (h *PastureEntryHandler) FindByAnimalId(w http.ResponseWriter, r *http.Request) {
	animalId := r.PathValue("animalId")
	list, err := h.Repository.FindByAnimalId(animalId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    handlersUtil.SendList(w, list)
}

func (h *PastureEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w, r, h.Repository)
}

func (h *PastureEntryHandler) Save(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w, r, h.Repository)
}

func (h *PastureEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
