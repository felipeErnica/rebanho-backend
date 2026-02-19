package slaughter

import (
	"net/http"
	"strings"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type SlaughterHandler struct {
	Service *SlaughterService
}

func (h *SlaughterHandler) GetLastAverageWeight(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastAverageWeight(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastDeadWeight(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastDeadWeight(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastPerformance(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastPerformance(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetWeightHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetWeightHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetRateHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetRateHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetBestFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetBestFathers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetBestMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetBestMothers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetBestButchers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetBestButchers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastEntries(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastGroups(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, SlaughterFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.FindPage(sort, order, cursor, filter, 100, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetPageFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, SlaughterFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.GetPageFoot(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) FindButcherPage(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	butcherId := r.PathValue("butcherId")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, err := util.DecodeFilter(r, SlaughterFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.FindButcherPage(sort, order, cursor, filter, butcherId, 100, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetButcherPageFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	butcherId := r.PathValue("butcherId")
	filter, err := util.DecodeFilter(r, SlaughterFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.GetButcherPageFoot(filter, butcherId, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	order := r.URL.Query().Get("order")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindGroups(order, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) FindEntries(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, SlaughterFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.FindEntries(sort, order, filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, SlaughterFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.GetEntriesFoot(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.Delete(id, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}

func (h *SlaughterHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("ids")
	ids := strings.Split(idStr, ",")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.DeleteBatch(ids, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}

func (h *SlaughterHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &SlaughterSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	result, err := h.Service.Update(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, result)
}

func (h *SlaughterHandler) UpdateBatch(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	batch, ok := util.DecodeEntity(w, r, &SlaughterSaveBatch{})
	if !ok {
		return
	}

	batch.UserId = userId
	err := h.Service.UpdateBatch(batch)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteUpdateResponse(w)
}

func (h *SlaughterHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &SlaughterSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Service.Add(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
}
