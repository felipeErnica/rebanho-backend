package pregnancyTests

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type TestEntryHandler struct {
	Repository *TestEntryRepository
}

func (h *TestEntryHandler) GetPregnancyRates(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetPregnancyRate(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetAnimalsNumber(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) GetBirthRates(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBirthRate(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetTestHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetPregnancyTestHist(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) GetNextBirths(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetNextBirths(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetRankedResults(w http.ResponseWriter, r *http.Request) {
	rankBy := r.URL.Query().Get("rankBy")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	var result *[]TestAnimal
	var err error

	switch rankBy {
	case "best-results":
		result, err = h.Repository.GetBestResults(userId)
	case "worst-results":
		result, err = h.Repository.GetWorstResults(userId)
	}

	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, err := handlersUtil.DecodeFilter(r, TestEntryFilter{})
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

func (h *TestEntryHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {

	filter, err := handlersUtil.DecodeFilter(r, TestEntryFilter{})
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetEntriesFoot(filter, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindGroups(userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) FindEntriesByGroup(w http.ResponseWriter, r *http.Request) {
	testDateString := r.PathValue("testDate")
	testDate, err := time.Parse(time.RFC3339Nano, testDateString)

	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindEntriesByGroup(sort, order, testDate, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) GetEntriesByGroupFoot(w http.ResponseWriter, r *http.Request) {
	testDateString := r.PathValue("testDate")
	testDate, err := time.Parse(time.RFC3339Nano, testDateString)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetEntriesByGroupFoot(testDate, userId)
	if err != nil {
		apiError.WriteError(w, err)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *TestEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &TestEntrySave{})
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

func (h *TestEntryHandler) Replace(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &TestEntrySave{})
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

func (h *TestEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := handlersUtil.DecodeEntity(w, r, &TestEntrySave{})
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

func (h *TestEntryHandler) UpdateBatch(w http.ResponseWriter, r *http.Request) {
	testDateString := r.PathValue("testDate")
	testDate, _ := time.Parse(time.RFC3339Nano, testDateString)
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	group, ok := handlersUtil.DecodeEntity(w, r, &TestGroups{})
	if !ok {
		return
	}

	group.UserId = userId
	response, err := h.Repository.UpdateBatch(testDate, group)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, response)
}

func (h *TestEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *TestEntryHandler) DeleteBatch(w http.ResponseWriter, r *http.Request) {
	testDateString := r.PathValue("testDate")
	testDate, _ := time.Parse(time.RFC3339Nano, testDateString)
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteBatch(testDate, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}

