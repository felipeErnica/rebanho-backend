package animalDashboard

import (
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
	"net/http"
)

type DashboardHandler struct {
	Repository *DashboardRepository
}

func (h *DashboardHandler) GroupByYear(w http.ResponseWriter, r *http.Request) {
	handlersUtil.SendTotalEntity(w, r, DashboardFilter{}, h.Repository.GroupByYear)
}

func (h *DashboardHandler) TotalBySex(w http.ResponseWriter, r *http.Request) {
	handlersUtil.SendTotalEntity(w, r, DashboardFilter{}, h.Repository.TotalBySex)
}

func (h *DashboardHandler) TotalByType(w http.ResponseWriter, r *http.Request) {
	handlersUtil.SendTotalEntity(w, r, DashboardFilter{}, h.Repository.TotalByType)
}

func (h *DashboardHandler) GroupByAgeAndFarm(w http.ResponseWriter, r *http.Request) {
    handlersUtil.SendGroupedList(w, r, DashboardFilter{},h.Repository.GroupByAgeAndFarm)
}

func (h *DashboardHandler) GroupByAgeAndPasture(w http.ResponseWriter, r *http.Request) {
    handlersUtil.SendGroupedList(w, r, DashboardFilter{},h.Repository.GroupByAgeAndPasture)
}

func (h *DashboardHandler) GroupByAge(w http.ResponseWriter, r *http.Request) {
    handlersUtil.SendGroupedList(w, r, DashboardFilter{},h.Repository.GroupByAge)
}
