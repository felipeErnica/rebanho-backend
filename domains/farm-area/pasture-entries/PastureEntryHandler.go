package pastureEntries

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type PastureEntryHandler struct {
	Repository *PastureEntryRepository
}

func (h *PastureEntryHandler) SearchPastureAnimals(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchPastureAnimals(pastureId, userId)
	if err != nil {
		apiError.WriteError(err, w)
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
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *PastureEntryHandler) FindByPastureTotal(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, ok := handlersUtil.DecodeFilter(w, r, PastureEntryFilter{})
	if !ok {
		return
	}

	result, err := h.Repository.FindByPastureTotal(pastureId, userId, filter)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *PastureEntryHandler) FindByAnimalId(w http.ResponseWriter, r *http.Request) {
	animalId := r.PathValue("animalId")
	list, err := h.Repository.FindByAnimalId(animalId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.SendList(w, list)
}
