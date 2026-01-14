package insemination

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type InseminationHandler struct {
	Repository *InseminationRepository
}

func (h *InseminationHandler) GetBirthRateStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBirthRateStats(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) GetPregnancyRateStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetPregnancyRateStats(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) GetInseminationHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetInseminationStats(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) GetAnimalsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetAnimalsNumber(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) GetFutureBirths(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetFutureBirths(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *InseminationHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *InseminationHandler) GetBestBull(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestBull(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	filter, err := handlersUtil.DecodeFilter(r, InseminationEntryFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindEntriesPage(userId, filter, sort, order, cursor)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) FindEntriesByGroup(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("inseminationDate")

	inseminationDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindEntriesByGroup(userId, inseminationDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindGroups(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) GetEntriesByGroupFoot(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("inseminationDate")
	inseminationDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetEntriesByGroupFoot(userId, inseminationDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := handlersUtil.DecodeFilter(r, InseminationEntryFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	result, err := h.Repository.GetEntriesFoot(userId, filter)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) SearchInseminationBulls(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchInseminationBulls(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) SearchNonInseminationBulls(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchNonInseminationBulls(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) SetAsInseminationBulls(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SetAsInseminationBull(id, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *InseminationHandler) AddInsemination(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &InseminationEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.AddInsemination(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *InseminationHandler) ReplaceInsemination(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &InseminationEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.ReplaceInsemination(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *InseminationHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *InseminationHandler) DeleteNoValidation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteNoValidation(id, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}

func (h *InseminationHandler) DeleteAndChangeFather(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteAndChangeFather(id, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}

func (h *InseminationHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &InseminationEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	res, err := h.Repository.Update(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *InseminationHandler) UpdateNoValidation(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &InseminationEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	res, err := h.Repository.UpdateNoValidation(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *InseminationHandler) UpdateBatch(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("inseminationDate")
	inseminationDate, parseErr := time.Parse(time.RFC3339Nano, queryDate)
	if parseErr != nil {
		apiError.WriteError(w, parseErr)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	group, ok := handlersUtil.DecodeEntity(w, r, &InseminationGroup{})
	if !ok {
		return
	}

	group.UserId = userId
	result, err := h.Repository.UpdateBatch(inseminationDate, group)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *InseminationHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("inseminationDate")
	inseminationDate, parseErr := time.Parse(time.RFC3339Nano, queryDate)
	if parseErr != nil {
		apiError.WriteError(w, parseErr)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteBatch(inseminationDate, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}
