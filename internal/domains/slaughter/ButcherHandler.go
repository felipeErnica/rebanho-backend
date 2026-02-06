package slaughter

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type ButcherHandler struct {
	Repository *ButcherRepository
}

func (h *ButcherHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindAll(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *ButcherHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindById(id, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *ButcherHandler) Search(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.Search(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *ButcherHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	butcherId := r.PathValue("id")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, err := util.DecodeFilter(r, SlaughterEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindEntriesPage(sort, order, cursor, filter, butcherId, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *ButcherHandler) FindPageFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	butcherId := r.PathValue("id")
	filter, err := util.DecodeFilter(r, SlaughterEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindEntriesPageFoot(filter, butcherId, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *ButcherHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &ButcherSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.Add(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
}

func (h *ButcherHandler) Replace(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &ButcherSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.Replace(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteUpdateResponse(w)
}

func (h *ButcherHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &ButcherSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	response, err := h.Repository.Update(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, response)
}

func (h *ButcherHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.Delete(id, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}
