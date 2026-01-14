package slaughter

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type SlaughterHandler struct {
	Repository *SlaughterRepository
}

func (h *SlaughterHandler) GetLastAverageWeight(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastAverageWeight(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastDeadWeight(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastDeadWeight(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastPerformance(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastPerformance(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetWeightHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetWeightHist(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetRateHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetRateHist(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetBestFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestFathers(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetBestMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestMothers(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetBestSlaughterHouses(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestSlaughterhouses(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastEntries(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastGroups(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := handlersUtil.DecodeFilter(r, SlaughterEntryFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindEntriesPage(sort, order, cursor, filter, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetEntriesPageFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := handlersUtil.DecodeFilter(r, SlaughterEntryFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	result, err := h.Repository.GetEntriesPageFoot(filter, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindGroups(order, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) FindEntriesByDate(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entryDateStr := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateStr)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindEntriesByDate(sort, order, entryDate, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetEntriesByDateFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entryDateStr := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateStr)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}


	result, err := h.Repository.GetEntriesByDateFoot(entryDate, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.Delete(id, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}

func (h *SlaughterHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &SlaughterEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	result, err := h.Repository.Update(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *SlaughterHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &SlaughterEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.Add(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *SlaughterHandler) Replace(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &SlaughterEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.Replace(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteUpdateResponse(w)
}

