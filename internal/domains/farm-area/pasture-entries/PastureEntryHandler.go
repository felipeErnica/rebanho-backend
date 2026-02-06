package pastureEntries

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type PastureEntryHandler struct {
	Repository *PastureEntryRepository
}

func (h *PastureEntryHandler) SearchPastureAnimals(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchPastureAnimals(pastureId, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.SendList(w, result)
}

func (h *PastureEntryHandler) FindByPasture(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
	cursor := r.URL.Query().Get("cursor")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, PastureEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindByPasture(pastureId, userId, filter, cursor, sort, order)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *PastureEntryHandler) FindByPastureTotal(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, PastureEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindByPastureTotal(pastureId, userId, filter)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *PastureEntryHandler) FindByAnimalId(w http.ResponseWriter, r *http.Request) {
	animalId := r.PathValue("animalId")
	list, err := h.Repository.FindByAnimalId(animalId)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.SendList(w, list)
}

func (h *PastureEntryHandler) AddEntry(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &PastureEntry{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.AddEntry(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
}

func (h *PastureEntryHandler) TransferEntry(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &PastureEntry{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.TransferEntry(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
}
