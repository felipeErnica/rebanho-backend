package lactation

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type LactationHandler struct {
	Service *LactationService
}

func (h *LactationHandler) GetLastLactating(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastLactating(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetLastDry(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastDry(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetDairyTypes(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetDairyTypes(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetBestAnimals(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBestAnimals(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetWorstAnimals(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetWorstAnimals(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetBestMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBestMothers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetWorstMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetWorstMothers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetBestFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBestFathers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetWorstFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetWorstFathers(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) FindLactationPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, err := util.DecodeFilter(r, LactationHistFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindLactationPage(filter, sort, order, cursor, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetLactationPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, err := util.DecodeFilter(r, LactationHistFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetLactationPageFoot(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) FindAnimalsPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, err := util.DecodeFilter(r, LactationAnimalFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindAnimalsPage(filter, sort, order, cursor, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) GetAnimalsPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, err := util.DecodeFilter(r, LactationAnimalFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetAnimalsPageFoot(filter, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}
func (h *LactationHandler) FindById(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindById(lacId, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *LactationHandler) AddLactation(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	lac, ok := util.DecodeEntity(w, r, &LactationHistSave{})
	if !ok {
		return
	}

	lac.UserId = userId
	apiErr := h.Service.AddLactation(lac)
	if apiErr != nil {
		log.WriteAPIError(apiErr, w)
		return
	}

	util.WriteCreatedResponse(w)
}

func (h *LactationHandler) UpdateLactation(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	lac, ok := util.DecodeEntity(w, r, &LactationHistSave{})
	if !ok {
		return
	}

	lac.UserId = userId
	res, err := h.Service.UpdateLactation(lac)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, res)
}

func (h *LactationHandler) DeleteLactation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.Repo.DeleteLactation(id, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}
