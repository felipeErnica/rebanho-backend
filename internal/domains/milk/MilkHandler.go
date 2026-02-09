package milk

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type MilkHandler struct {
	Service *MilkService
}

func (h *MilkHandler) GetLastMilk(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastMilk(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetLastAverageMilk(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastAverageMilk(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetMilkProduction(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetMilkProduction(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetYearMilkProduction(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetYearMilk(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetYearAverageMilk(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetYearAverageMilk(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *MilkHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *MilkHandler) FindGroupsPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")

	filter, err := util.DecodeFilter(r, LactationGroupFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindGroupsPage(filter, order, cursor, 100, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetLactationEntries(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")
	result, err := h.Service.GetLactationEntries(lacId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetLactationEntriesFoot(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")
	result, err := h.Service.GetLactationEntriesFoot(lacId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")

	filter, err := util.DecodeFilter(r, MilkEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindPage(filter, sort, order, cursor, 100, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetEntriesPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, err := util.DecodeFilter(r, MilkEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetPageFoot(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetGroupEntries(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateVar)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetGroupEntries(userId, entryDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) GetGroupEntriesFoot(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.PathValue("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateVar)

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetGroupEntriesFoot(userId, entryDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *MilkHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	groupEntry, ok := util.DecodeEntity(w, r, &LactationGroupSave{})
	if !ok {
		return
	}

	groupEntry.UserId = userId
	res, err := h.Service.UpdateGroup(groupEntry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, res)
}

func (h *MilkHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	entryDate, parseErr := time.Parse(time.RFC3339Nano, r.PathValue("entryDate"))
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.DeleteGroup(entryDate, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteDeleteResponse(w)
}

func (h *MilkHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	milkEntry, ok := util.DecodeEntity(w, r, &MilkEntrySave{})
	if !ok {
		return
	}

	milkEntry.UserId = userId
	err := h.Service.Add(milkEntry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
}

func (h *MilkHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	milkEntry, ok := util.DecodeEntity(w, r, &MilkEntrySave{})
	if !ok {
		return
	}

	milkEntry.UserId = userId
	res, err := h.Service.Update(milkEntry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, res)
}

func (h *MilkHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.Service.Delete(id)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteDeleteResponse(w)
}
