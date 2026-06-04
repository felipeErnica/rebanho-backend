package birth

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type BirthHandler struct {
	Service *BirthService
}

func (h *BirthHandler) GetPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, err := util.DecodeFilter(r, BirthEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetPageFoot(userId, filter)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.WriteEntity(w, result)
}

func (h *BirthHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, err := util.DecodeFilter(r, BirthEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.FindPage(userId, sort, order, filter, cursor, 100)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	res, err := h.Service.GetById(id)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.WriteEntity(w, res)
}

func (h *BirthHandler) GetBestIntervals(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBestIntervals(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *BirthHandler) GetLastBirths(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastBirths(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *BirthHandler) GetLastBirthsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastBirthsNumber(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetYearBirthsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetYearBirthsNumber(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetYearDeathsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetYearDeathsNumber(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetWorstIntervals(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetWorstIntervals(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *BirthHandler) GetBirthIntervalHistory(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetBirthIntervalHistory(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetDeathIndex(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetDeathIndex(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetBirthHistory(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBirthHistory(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) TotalBySex(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.TotalBySex(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetYearBySex(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetYearBySex(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}
	util.WriteEntity(w, result)
}

func (h *BirthHandler) AddBirth(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, ok := util.DecodeEntity(w, r, &BirthEntrySave{})
	if !ok {
		return
	}

	birthEntry.UserId = userId
	err := h.Service.AddBirth(birthEntry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
}

func (h *BirthHandler) UpdateBirth(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, ok := util.DecodeEntity(w, r, &BirthEntrySave{})
	if !ok {
		return
	}

	birthEntry.UserId = userId
	result, err := h.Service.UpdateBirth(birthEntry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) DeleteBirth(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	skipValidation, parseErr := util.ParseBool(r.URL.Query().Get("skipValidation"))
	if parseErr != nil {
		log.WriteError(w, parseErr)
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.DeleteBirth(id, userId, skipValidation)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BirthHandler) GetPotentialFather(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	birthEntry, err := util.DecodeURL(r, BirthEntrySave{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	birthEntry.UserId = userId
	result, apiErr := h.Service.GetPotentialFather(birthEntry)
	if apiErr != nil {
		log.WriteAPIError(apiErr, w)
		return
	}

	util.WriteEntity(w, result)
}
