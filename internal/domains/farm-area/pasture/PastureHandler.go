package pasture

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type PastureHandler struct {
	Repository *PastureRepository
}

func (h *PastureHandler) SearchPasture(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, PastureFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	list, err := h.Repository.SearchPasture(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, list)
}

func (h *PastureHandler) FindAnimalsByPasture(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("id")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}
	result, err := h.Repository.FindAnimalsByPasture(pastureId, userId, sort, order)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.SendList(w, result)
}
