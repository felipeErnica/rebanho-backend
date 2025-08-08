package insemination

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type InseminationHandler struct {
	Repository *InseminationRepository
}

func (h *InseminationHandler) GetBirthRateStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetBirthRateStats(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendEntity(w, result)
}

func (h *InseminationHandler) GetInseminationHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetInseminationStats(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendEntity(w, result)
}

func (h *InseminationHandler) GetBestBull(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetBestBull(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendEntity(w, result)
}
