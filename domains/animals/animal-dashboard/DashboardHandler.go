package animalDashboard

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type DashboardHandler struct {
	Repository *DashboardRepository
}

func (h *DashboardHandler) TotalBySex(w http.ResponseWriter, r *http.Request) {
    filter, ok := handlersUtil.DecodeFilter(w ,r , AnimalsDashboardFilter{}); if !ok {
        return
    }
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }
    total, err := h.Repository.TotalBySex(userId, filter)
    if err != nil {
        serverErrors.DatabaseGetError(err , w)
        return
    }
    handlersUtil.SendEntity(w, total)
}

func (h *DashboardHandler) GroupByAgeAndFarm(w http.ResponseWriter, r *http.Request) {
    filter, ok := handlersUtil.DecodeFilter(w, r, AnimalsDashboardFilter{}); if !ok {
        return
    }
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }
    obj, err := h.Repository.GroupByAgeAndFarm(userId, filter)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    handlersUtil.SendList(w, obj)
}

func (h *DashboardHandler) GroupByAge(w http.ResponseWriter, r *http.Request) {
    filter, ok := handlersUtil.DecodeFilter(w, r, AnimalsDashboardFilter{}); if !ok {
        return
    }
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }
    obj, err := h.Repository.GroupByAge(userId, filter)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    handlersUtil.SendList(w, obj)
}
