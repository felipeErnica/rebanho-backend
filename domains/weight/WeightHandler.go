package weight

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type WeightHandler struct {
	Repository *WeightRepository
}

func (h *WeightHandler) GetYearWeightGain(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetYearWeightGain(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *WeightHandler) GetYearWeight(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetYearWeight(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *WeightHandler) GetLastWeightGain(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetLastWeightGain(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *WeightHandler) GetLastWeight(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetLastWeight(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *WeightHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *WeightHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *WeightHandler) GetBestMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetBestMothers(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *WeightHandler) GetBestFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetBestFathers(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *WeightHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	filter, ok := handlersUtil.DecodeFilter(w, r, WeightFilter{}); if !ok {
		return
	}

	result, err := h.Repository.FindEntriesPage(sort, order, cursor, filter, userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}


func (h *WeightHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.FindGroups(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *WeightHandler) FindEntriesByDate(w http.ResponseWriter, r *http.Request) {
	entryDateStr := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateStr)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.FindEntriesByDate(entryDate, userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}
