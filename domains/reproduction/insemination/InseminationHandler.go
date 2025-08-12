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
    
    handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) GetPregnantsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetPregnantNumbers(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendEntity(w, result)
}

func (h *InseminationHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *InseminationHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *InseminationHandler) GetBestBull(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetBestBull(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
    cursor := r.URL.Query().Get("cursor")
    sort := r.URL.Query().Get("sort")
    order := r.URL.Query().Get("order")

    filter, ok := handlersUtil.DecodeFilter(w, r, InseminationEntryFilter{}); if !ok {
        return
    }

	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.FindEntriesPage(userId, filter, sort, order, cursor)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendEntity(w, result)
}

func (h *InseminationHandler) FindEntriesByGroup(w http.ResponseWriter, r *http.Request) {
    groupId := r.PathValue("groupId")
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.FindEntriesByGroup(userId, groupId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.FindGroups(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) GetEntriesByGroupFoot(w http.ResponseWriter, r *http.Request) {
    groupId := r.PathValue("groupId")
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetEntriesByGroupFoot(userId, groupId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendEntity(w, result)
}

func (h *InseminationHandler) GetGroupsFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetGroupsFoot(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
    
    handlersUtil.SendEntity(w, result)
}
