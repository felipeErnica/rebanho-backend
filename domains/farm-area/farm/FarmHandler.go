package farm

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type FarmHandler struct {
	Repository *FarmRepository
}

func (h *FarmHandler) FindFarmAnimals(w http.ResponseWriter, r *http.Request) {
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

	result, err := h.Repository.FindFarmAnimals(farmId, userId, filter, sort, order, cursor)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.WriteEntity(w, result)
}

func (h *FarmHandler) FindFarmAnimalsTotal(w http.ResponseWriter, r *http.Request) {
	farmId := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, ok := handlersUtil.DecodeFilter(w, r, FarmAnimalFilter{})
	if !ok {
		return
	}

	result, err := h.Repository.FindFarmAnimalsTotal(farmId, userId, filter)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.WriteEntity(w, result)
}

func (h *FarmHandler) SearchFarm(w http.ResponseWriter, r *http.Request) {
	handlersUtil.ReturnSearchResults(w, r, h.Repository.SearchFarmById, h.Repository.SearchFarm)
}
