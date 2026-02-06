package dashboard

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type FarmDashboardHandler struct {
	Repository *FarmDashboardRepository
}

func (h *FarmDashboardHandler) FarmInfo(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}
	result, err := h.Repository.GetFarmInfo(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.SendList(w, result)
}

func (h *FarmDashboardHandler) PastureInfo(w http.ResponseWriter, r *http.Request) {
	farmId := r.URL.Query().Get("farmId")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}
	result, err := h.Repository.GetPastureInfo(userId, farmId)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.SendList(w, result)
}
