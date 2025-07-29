package pastureEntries

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type PastureEntryHandler struct {
	Repository *PastureEntryRepository
}

func (h *PastureEntryHandler) SearchAnimalsForPasture(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

    result, err := h.Repository.SearchAnimalsForPasture(pastureId, userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendList(w, result)
}

func (h *PastureEntryHandler) FindByPasture(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
    cursor := r.URL.Query().Get("cursor")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

    filter, ok := handlersUtil.DecodeFilter(w, r, PastureEntryFilter{})
    if !ok {
        return
    }

    result, err := h.Repository.FindByPasture(pastureId, userId, filter, cursor, sort, order)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
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
