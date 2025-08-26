package lactation

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type LactationHandler struct {
	Repository *LactationRepository
}

func (h *LactationHandler) GetYearlyMilk(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetYearlyMilk(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetMonthMilk(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetMonthMilk(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetAnimalsAverage(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetAnimalsAverage(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetMilkProduction(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetMilkProduction(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetBestAnimals(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetBestAnimals(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetWorstAnimals(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetWorstAnimals(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetLastEntries(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}
