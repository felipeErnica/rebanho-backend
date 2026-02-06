package animals

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type AnimalHandler struct {
	Service *AnimalService
}

func (h *AnimalHandler) GetDairyHist(w http.ResponseWriter, r *http.Request) {

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Service.GetDairyHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, res)
}

func (h *AnimalHandler) GetBirthHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Service.GetBirthHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, res)
}

func (h *AnimalHandler) GetDeathHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Service.GetDeathHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, res)
}

func (h *AnimalHandler) GetSlaughterHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Service.GetSlaughterHist(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, res)
}

func (h *AnimalHandler) GetAnimalTypes(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Service.GetAnimalTypes(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, res)
}

func (h *AnimalHandler) GetLastDeaths(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Service.GetLastDeaths(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, res)
}

func (h *AnimalHandler) GetAgeAndSex(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Service.GetAgeAndSex(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, res)
}

func (h *AnimalHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, AnimalFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.FindPage(userId, cursor, sort, order, 100, filter)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *AnimalHandler) GetPageFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, AnimalFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.GetPageFoot(userId, filter)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *AnimalHandler) FindById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	response, err := h.Service.FindById(id, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, response)
}

func (h *AnimalHandler) Search(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, AnimalFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.Search(sort, order, filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *AnimalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	skipValidation, parseErr := util.ParseBool(r.URL.Query().Get("skipValidation"))
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.Delete(skipValidation, id, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}

func (h *AnimalHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	newEntry, ok := util.DecodeEntity(w, r, &AnimalSave{})
	if !ok {
		return
	}

	newEntry.UserId = userId
	response, err := h.Service.Update(newEntry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, response)
}

func (h *AnimalHandler) Add(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	newEntry, ok := util.DecodeEntity(w, r, &AnimalSave{})
	if !ok {
		return
	}

	newEntry.UserId = userId
	apiErr := h.Service.Add(newEntry)
	if apiErr != nil {
		log.WriteAPIError(apiErr, w)
		return
	}

	util.WriteCreatedResponse(w)
}
