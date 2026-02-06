package farm

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type FarmHandler struct {
	Repository *FarmRepository
}

func (h *FarmHandler) FindFarmAnimals(w http.ResponseWriter, r *http.Request) {
	farmId := r.PathValue("id")
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	cursor := r.URL.Query().Get("cursor")

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, FarmAnimalFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindFarmAnimals(farmId, userId, filter, sort, order, cursor)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.WriteEntity(w, result)
}

func (h *FarmHandler) FindFarmAnimalsTotal(w http.ResponseWriter, r *http.Request) {
	farmId := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, FarmAnimalFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindFarmAnimalsTotal(farmId, userId, filter)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.WriteEntity(w, result)
}
