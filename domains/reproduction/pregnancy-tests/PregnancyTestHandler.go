package pregnancyTests

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type TestEntryHandler struct {
	Repository *TestEntryRepository
}

func (h *TestEntryHandler) GetPregnancyRates (w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetPregnancyRate(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *TestEntryHandler) GetBirthRates (w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetBirthRate(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *TestEntryHandler) GetTestHist (w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetPregnancyTestHist(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *TestEntryHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetLastGroups(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *TestEntryHandler) GetNextBirths(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetNextBirths(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *TestEntryHandler) GetBestResults(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetBestResults(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}
