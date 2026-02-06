package insemination

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type InseminationHandler struct {
	Service *InseminationService
}

func (h *InseminationHandler) GetBirthRateStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetBirthRateStats(userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) GetPregnancyRateStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetPregnancyRateStats(userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) GetInseminationHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetInseminationStats(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *InseminationHandler) GetAnimalsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetAnimalsNumber(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) GetFutureBirths(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetFutureBirths(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetLastGroups(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetLastEntries(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) GetBestBull(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBestBull(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *InseminationHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	filter, err := util.DecodeFilter(r, InseminationEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindEntriesPage(userId, filter, sort, order, cursor)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) FindEntriesByGroup(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("inseminationDate")

	inseminationDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindEntriesByGroup(userId, inseminationDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *InseminationHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindGroups(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *InseminationHandler) GetEntriesByGroupFoot(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("inseminationDate")
	inseminationDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetEntriesByGroupFoot(userId, inseminationDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, InseminationEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.Repo.GetEntriesFoot(userId, filter)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &InseminationEntrySave{})
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

func (h *InseminationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	params, parseErr := util.DecodeURL(r, InseminationEntryDelete{})
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	params.UserId = userId
	err := h.Service.Delete(params)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}

func (h *InseminationHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &InseminationEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	res, err := h.Service.Update(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, res)
}

func (h *InseminationHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	group, ok := util.DecodeEntity(w, r, &InseminationGroupSave{})
	if !ok {
		return
	}

	group.UserId = userId
	result, err := h.Service.UpdateGroup(group)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, result)
}

func (h *InseminationHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	params, parseErr := util.DecodeURL(r, InseminationGroupDelete{})
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	params.UserId = userId
	err := h.Service.DeleteGroup(params)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}
