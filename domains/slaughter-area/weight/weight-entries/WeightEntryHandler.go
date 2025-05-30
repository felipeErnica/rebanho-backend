package weightEntries

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type WeightEntryHandler struct {
	Repository *WeightEntryRepository
}

func (h *WeightEntryHandler) FindByGroupId(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("groupId")
	list, err := h.Repository.FindByGroupId(groupId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    handlersUtil.SendList(w, list)
}

func (h *WeightEntryHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById(w, r, h.Repository)
}

func (h *WeightEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w, r, h.Repository)
}

func (h *WeightEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w, r, h.Repository)
}

func (h *WeightEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
