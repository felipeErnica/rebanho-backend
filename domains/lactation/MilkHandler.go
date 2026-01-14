package lactation

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type MilkHandler struct {
	Repository *MilkRepository
}

func (h *MilkHandler) FindGroupsPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")

	filter, err := handlersUtil.DecodeFilter(r, LactationGroupFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindGroupsPage(filter, order, cursor, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *MilkHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")

	filter, err := handlersUtil.DecodeFilter(r, MilkEntryFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindEntriesPage(filter, sort, order, cursor, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *MilkHandler) GetEntriesPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, err := handlersUtil.DecodeFilter(r, MilkEntryFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetEntriesPageFoot(filter, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *MilkHandler) GetGroupEntries(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateVar)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetGroupEntries(userId, entryDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *MilkHandler) GetGroupEntriesFoot(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateVar)

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetGroupEntriesFoot(userId, entryDate)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *MilkHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.PathValue("entryDate")
	entryDate, parseErr := time.Parse(time.RFC3339Nano, entryDateVar)
	if parseErr != nil {
		apiError.WriteError(w, parseErr)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	groupEntry, ok := handlersUtil.DecodeEntity(w, r, &LactationGroupSave{})
	if !ok {
		return
	}

	groupEntry.UserId = userId
	res, err := h.Repository.UpdateGroup(entryDate, groupEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *MilkHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.PathValue("entryDate")
	entryDate, parseErr := time.Parse(time.RFC3339Nano, entryDateVar)
	if parseErr != nil {
		apiError.WriteError(w, parseErr)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteGroup(entryDate, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}

func (h *MilkHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	milkEntry, ok := handlersUtil.DecodeEntity(w, r, &MilkEntrySave{})
	if !ok {
		return
	}

	milkEntry.UserId = userId
	err := h.Repository.Add(milkEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *MilkHandler) Replace(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	milkEntry, ok := handlersUtil.DecodeEntity(w, r, &MilkEntrySave{})
	if !ok {
		return
	}

	milkEntry.UserId = userId
	err := h.Repository.Replace(milkEntry, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteUpdateResponse(w)
}

func (h *MilkHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	milkEntry, ok := handlersUtil.DecodeEntity(w, r, &MilkEntry{})
	if !ok {
		return
	}

	milkEntry.UserId = userId
	res, err := h.Repository.Update(milkEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *MilkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.Repository.Delete(id)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}
