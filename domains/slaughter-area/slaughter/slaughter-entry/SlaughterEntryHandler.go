package slaughterEntry

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type SlaughterEntryHandler struct {
	Repository *SlaughterEntryRepository
}

func (h *SlaughterEntryHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	filter := SlaughterEntryFilter{}
	handlersUtil.DecodeFilter(w, r, &filter)
	handlersUtil.ReturnPage(w, r, h.Repository, filter)
}

func (h *SlaughterEntryHandler) FindByGroupId(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("groupId")
	list, err := h.Repository.FindByGroupId(groupId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
	handlersUtil.SendList(w, list)
}

func (h *SlaughterEntryHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById(w, r, h.Repository)
}

func (h *SlaughterEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w, r, h.Repository)
}

func (h *SlaughterEntryHandler) Save(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w, r, h.Repository)
}

func (h *SlaughterEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
