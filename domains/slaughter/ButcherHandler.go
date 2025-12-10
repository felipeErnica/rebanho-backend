package slaughter

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type ButcherHandler struct {
	Repository *ButcherRepository
}

func (h *ButcherHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindAll(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *ButcherHandler) Search(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.Search(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *ButcherHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &ButcherSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.Add(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *ButcherHandler) Replace(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &ButcherSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.Replace(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteUpdateResponse(w)
}

func (h *ButcherHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &ButcherSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	response, err := h.Repository.Update(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, response)
}

func (h *ButcherHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.Delete(id, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}
