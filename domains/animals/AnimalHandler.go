package animals

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
	"net/http"
)

type AnimalHandler struct {
	Repository *AnimalRepository
}

func (h *AnimalHandler) GetDairyHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GetDairyHist(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GetBirthHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GetBirthHist(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GetDeathHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GetDeathHist(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GetSlaughterHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GetSlaughterHist(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GetAnimalTypes(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GetAnimalTypes(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GetLastDeaths(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GetLastDeaths(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GetAgeAndSex(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GetAgeAndSex(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, ok := handlersUtil.DecodeFilter(w, r, AnimalFilter{})
	if !ok {
		return
	}

	result, err := h.Repository.FindPage(userId, cursor, sort, order, filter)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) GetPageFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, ok := handlersUtil.DecodeFilter(w, r, AnimalFilter{})
	if !ok {
		return
	}

	result, err := h.Repository.GetPageFoot(userId, filter)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	response, err := h.Repository.FindById(id, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, response)
}

func (h *AnimalHandler) Search(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	filter, ok := handlersUtil.DecodeFilter(w, r, AnimalFilter{})
	if !ok {
		return
	}

	result, err := h.Repository.Search(sort, order, filter, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

func (h *AnimalHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	newEntry, ok := handlersUtil.DecodeEntity(w, r, &AnimalSave{})
	if !ok {
		return
	}

	newEntry.UserId = userId
	response, err := h.Repository.Update(newEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, response)
}

func (h *AnimalHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	newEntry, ok := handlersUtil.DecodeEntity(w, r, &AnimalSave{})
	if !ok {
		return
	}

	newEntry.UserId = userId
	err := h.Repository.Add(newEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}
