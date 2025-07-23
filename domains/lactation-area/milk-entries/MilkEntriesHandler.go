package milkEntries

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type MilkHandler struct {
	Repository *MilkRepository
}

func (h *MilkHandler) FindPage(w http.ResponseWriter, r *http.Request) {
}

func (h *MilkHandler) FindByCow(w http.ResponseWriter, r *http.Request) {
	animalId := r.PathValue("animalId")
	milkList, err := h.Repository.FindByAnimal(animalId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    handlersUtil.SendList(w, milkList)
}

func (h *MilkHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w, r, h.Repository)
}

func (h *MilkHandler) Update(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Update(w, r, h.Repository)
}

func (h *MilkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Delete(w, r, h.Repository)
}
