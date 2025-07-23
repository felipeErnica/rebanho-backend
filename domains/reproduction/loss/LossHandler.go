package loss

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type PregnancyLossHandler struct {
    Repository *LossRepository
}

func (h *PregnancyLossHandler) FindByAnimalId(w http.ResponseWriter, r *http.Request) {
    animalId:=r.PathValue("animalId")
    list, err:= h.Repository.FindByAnimalId(animalId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    handlersUtil.SendList(w, list)
}

func (h *PregnancyLossHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById(w, r, h.Repository)
}

func (h *PregnancyLossHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w, r, h.Repository)
}

func (h *PregnancyLossHandler) Update(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w, r, h.Repository)
}

func (h *PregnancyLossHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
