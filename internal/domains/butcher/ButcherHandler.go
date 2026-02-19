package butcher

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type ButcherHandler struct {
	Service *ButcherService
}

func (h *ButcherHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindAll(userId)
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

	result, err := h.Service.FindById(id, userId)
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
	err := h.Service.Add(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
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
	response, err := h.Service.Update(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, response)
}

func (h *ButcherHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, parseErr := util.DecodeURL(r,ButcherDelete{})
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	entry.UserId = userId
	err := h.Service.Delete(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}
