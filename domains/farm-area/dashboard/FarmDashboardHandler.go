package dashboard

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type FarmDashboardHandler struct {
	Repository *FarmDashboardRepository
}

func (h *FarmDashboardHandler) FarmInfo(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r)
    if !ok {
        return
    }
    result, err := h.Repository.GetFarmInfo(userId)
    if err != nil {
        apiError.WriteError(w, err )
        return
    }
    handlersUtil.SendList(w, result)
}

func (h *FarmDashboardHandler) PastureInfo(w http.ResponseWriter, r *http.Request) {
    farmId := r.URL.Query().Get("farmId")
    userId, ok := handlersUtil.GetUserId(w, r)
    if !ok {
        return
    }
    result, err := h.Repository.GetPastureInfo(userId, farmId)
    if err != nil {
        apiError.WriteError(w, err )
        return
    }
    handlersUtil.SendList(w, result)
}
