package loss

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type LossHandler struct {
	Repository *LossRepository
}

func (h *LossHandler) GetLossRate(w http.ResponseWriter, r *http.Request) {

    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetLossRate(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    
    handlersUtil.SendEntity(w, result)
}

func (h *LossHandler) GetLossesHist(w http.ResponseWriter, r *http.Request) {

    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetLossesHist(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    
    handlersUtil.SendEntity(w, result)
}

func (h *LossHandler) GetMostLossesAnimals(w http.ResponseWriter, r *http.Request) {

    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetMostLossesAnimals(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }
    
    handlersUtil.SendEntity(w, result)
}

func (h *LossHandler) FindPage(w http.ResponseWriter, r *http.Request) {
    sort := r.URL.Query().Get("sort")
    order := r.URL.Query().Get("order")
    cursor := r.URL.Query().Get("cursor")

    filter, ok := handlersUtil.DecodeFilter(w, r, LossFilter{}); if !ok {
        return
    }

    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.FindPage(filter, cursor, sort, order, userId)
    if err != nil {
        serverErrors.DatabaseGetError(err , w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LossHandler) GetPageFoot(w http.ResponseWriter, r *http.Request) {
    filter, ok := handlersUtil.DecodeFilter(w, r, LossFilter{}); if !ok {
        return
    }

    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetPageFoot(filter, userId)
    if err != nil {
        serverErrors.DatabaseGetError(err , w)
        return
    }

    handlersUtil.SendEntity(w, result)
}
