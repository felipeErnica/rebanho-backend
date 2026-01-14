package embryoTransfer

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type TransferHandler struct {
	Repository *TransferRepository
}

func (h *TransferHandler) GetBirthRateStats(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetPregnancyRateStats(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetTransferHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetTransferHist(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) GetAnimalsNumber(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetFutureBirths(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetBestBull(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetBestReceivers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestReceivers(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) GetBestDonors(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestDonors(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	filter, err := handlersUtil.DecodeFilter(r, TransferEntryFilter{})
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

func (h *TransferHandler) FindEntriesByGroup(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("transferDate")

	transferDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindEntriesByGroup(userId, transferDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) GetEntriesByGroupFoot(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("transferDate")

	transferDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetEntriesByGroupFoot(userId, transferDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TransferHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := handlersUtil.DecodeFilter(r, TransferEntryFilter{})
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

func (h *TransferHandler) SearchTransferBulls(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchTransferBulls(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) SearchNonTransferBulls(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchNonTransferBulls(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) UpdateAsTransferBulls(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.UpdateAsTransferBulls(id, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteUpdateResponse(w)
}

func (h *TransferHandler) SearchEmbryoDonors(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchEmbryoDonors(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) SearchNonEmbryoDonors(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchNonEmbryoDonors(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *TransferHandler) UpdateAsEmbryoDonors(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.UpdateAsEmbryoDonors(id, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteUpdateResponse(w)
}

func (h *TransferHandler) AddTransfer(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &EmbryoTransferSave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Repository.AddTransfer(entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *TransferHandler) Replace(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &EmbryoTransferSave{})
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

func (h *TransferHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *TransferHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &EmbryoTransferSave{})
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

func (h *TransferHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("transferDate")
	transferDate, _ := time.Parse(time.RFC3339Nano, dateStr)
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteGroup(transferDate, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}

func (h *TransferHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	dateStr := r.PathValue("transferDate")
	transferDate, _ := time.Parse(time.RFC3339Nano, dateStr)
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &TransferGroup{})
	if !ok {
		return
	}

	entry.UserId = userId
	res, err := h.Repository.UpdateGroup(transferDate, entry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}
