package pregnancyTests

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type TestEntryHandler struct {
	Service *TestService
}

func (h *TestEntryHandler) GetPregnancyRates(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetPregnancyRates(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetAnimalsNumber(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) GetBirthRates(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetBirthRates(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetTestHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetTestHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) GetNextBirths(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetNextBirths(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetRankedResults(w http.ResponseWriter, r *http.Request) {
	rankBy := r.URL.Query().Get("rankBy")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetRankedResults(rankBy, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, err := util.DecodeFilter(r, TestFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindEntriesPage(filter, sort, order, cursor, 100, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {

	filter, err := util.DecodeFilter(r, TestFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetEntriesFoot(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindGroups(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) FindEntriesByGroup(w http.ResponseWriter, r *http.Request) {
	testDateString := r.PathValue("testDate")
	testDate, err := time.Parse(time.RFC3339Nano, testDateString)

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindEntriesByGroup(sort, order, testDate, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetEntriesByGroupFoot(w http.ResponseWriter, r *http.Request) {
	testDateString := r.PathValue("testDate")
	testDate, err := time.Parse(time.RFC3339Nano, testDateString)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetEntriesByGroupFoot(testDate, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *TestEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &TestSave{})
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

func (h *TestEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &TestSave{})
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

func (h *TestEntryHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	group, ok := util.DecodeEntity(w, r, &TestGroupSave{})
	if !ok {
		return
	}

	group.UserId = userId
	response, err := h.Service.UpdateGroup(group)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, response)
}

func (h *TestEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	testDateString := r.PathValue("testDate")
	testDate, _ := time.Parse(time.RFC3339Nano, testDateString)
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.DeleteGroup(testDate, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}
