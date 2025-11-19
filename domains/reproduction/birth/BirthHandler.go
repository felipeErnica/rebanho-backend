package birth

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type BirthHandler struct {
	Repository *BirthRepository
}

func (h *BirthHandler) FindPageFooter(w http.ResponseWriter, r *http.Request) {
	filter, ok := handlersUtil.DecodeFilter(w, r, BirthEntryFilter{})
	if !ok {
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindPageFooter(userId, filter)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, ok := handlersUtil.DecodeFilter(w, r, BirthEntryFilter{})
	if !ok {
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindPage(userId, sort, order, filter, cursor)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetBestIntervals(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestIntervals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *BirthHandler) GetLastBirths(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastBirths(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *BirthHandler) GetLastBirthsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastBirthsNumber(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetYearBirthsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetYearBirthsNumber(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetYearDeathsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetYearDeathsNumber(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetWorstIntervals(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetWorstIntervals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *BirthHandler) GetBirthIntervalHistory(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBirthIntervalHistory(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetDeathIndex(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetDeathIndex(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetBirthHistory(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBirthHistory(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) TotalBySex(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.TotalBySex(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetYearBySex(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetYearBySex(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) AddBirth(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, ok := handlersUtil.DecodeEntity(w, r, &BirthEntrySave{})
	if !ok {
		return
	}

	birthEntry.UserId = userId
	err := h.Repository.AddBirth(birthEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *BirthHandler) AddBirthNoValidation(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, ok := handlersUtil.DecodeEntity(w, r, &BirthEntrySave{})
	if !ok {
		return
	}

	birthEntry.UserId = userId
	err := h.Repository.AddBirthNoValidation(birthEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *BirthHandler) ReplaceBirth(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, ok := handlersUtil.DecodeEntity(w, r, &BirthEntrySave{})
	if !ok {
		return
	}

	birthEntry.UserId = userId
	err := h.Repository.ReplaceBirth(birthEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *BirthHandler) UpdateBirth(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, ok := handlersUtil.DecodeEntity(w, r, &BirthEntrySave{})
	if !ok {
		return
	}

	birthEntry.UserId = userId
	result, err := h.Repository.UpdateBirth(birthEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *BirthHandler) GetFather(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, ok := handlersUtil.DecodeEntity(w, r, &BirthEntrySave{})
	if !ok {
		return
	}

	birthEntry.UserId = userId
	result, err := h.Repository.GetFather(birthEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}
