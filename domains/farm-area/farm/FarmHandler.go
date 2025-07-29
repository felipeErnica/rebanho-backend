package farm

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type FarmHandler struct {
	Repository *FarmRepository
}

func (h *FarmHandler) FindAnimalsByFarm(w http.ResponseWriter, r *http.Request) {
	farmId := r.PathValue("id")
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	cursor := r.URL.Query().Get("cursor")

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, ok := handlersUtil.DecodeFilter(w, r, FarmAnimalFilter{})
	if !ok {
		return
	}

	result, err := h.Repository.FindAnimalsByFarm(farmId, userId, filter, sort, order, cursor)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendEntity(w, result)
}

func (h *FarmHandler) SearchFarm(w http.ResponseWriter, r *http.Request) {
	handlersUtil.ReturnSearchResults(w, r, h.Repository.SearchFarmById, h.Repository.SearchFarm)
}

func (h *FarmHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w, r, h.Repository)
}
