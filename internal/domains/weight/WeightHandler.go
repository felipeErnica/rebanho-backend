package weight

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type WeightHandler struct {
	Repository *WeightRepository
}

func (h *WeightHandler) GetWeightGainHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetWeightGainHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetWeightHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetWeightHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetLastWeightGain(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastWeightGain(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetLastWeight(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastWeight(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastEntries(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastGroups(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetBestMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestMothers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetBestFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestFathers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, WeightFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.FindEntriesPage(sort, order, cursor, filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, WeightFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Repository.GetEntriesPageFoot(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindGroups(userId, order)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) FindEntriesByDate(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")
	entryDateStr := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateStr)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindEntriesByDate(entryDate, userId, order, sort)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) GetEntriesByDateFoot(w http.ResponseWriter, r *http.Request) {
	entryDateStr := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateStr)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetEntriesByDateFoot(entryDate, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *WeightHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.Delete(id, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}

func (h *WeightHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &WeightEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	response, err := h.Repository.Update(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, response)
}
