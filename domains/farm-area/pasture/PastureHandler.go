package pasture

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type PastureHandler struct {
	Repository *PastureRepository
}

func (h *PastureHandler) SearchPasture(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := handlersUtil.DecodeFilter(r, PastureFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	list, err := h.Repository.SearchPasture(filter, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, list)
}

func (h *PastureHandler) SearchAllPastures(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	list, err := h.Repository.SearchAllPastures(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, list)
}

func (h *PastureHandler) FindAnimalsByPasture(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("id")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}
	result, err := h.Repository.FindAnimalsByPasture(pastureId, userId, sort, order)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}
	handlersUtil.SendList(w, result)
}
