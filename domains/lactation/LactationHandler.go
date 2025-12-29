package lactation

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type LactationHandler struct {
	Repository *LactationRepository
}

func (h *LactationHandler) GetLastMilk(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastMilk(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetLastAverageMilk(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastAverageMilk(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetLastAnimalsCount(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastAnimalsCount(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetMilkProduction(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetMilkProduction(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetYearMilkProduction(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetYearMilk(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetYearAverageMilk(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetYearAverageMilk(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetBestAnimals(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetWorstAnimals(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetWorstAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetBestMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestMothers(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetWorstMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetWorstMothers(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetBestFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestFathers(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetWorstFathers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetWorstFathers(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastEntries(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLastGroups(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) FindLactationPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, ok := handlersUtil.DecodeFilter(w, r, LactationHistFilter{})
	if !ok {
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindLactationPage(filter, sort, order, cursor, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) FindById(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindById(lacId, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetLactationPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, ok := handlersUtil.DecodeFilter(w, r, LactationHistFilter{})
	if !ok {
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetLactationPageFoot(filter, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetLactationEntries(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")
	result, err := h.Repository.GetLactationEntries(lacId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) GetLactationEntriesFoot(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")
	result, err := h.Repository.GetLactationEntriesFoot(lacId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) SearchLactatingAnimals(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchLactatingAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) SearchDryAnimals(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchDryAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) SearchNewLactationCalf(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchNewLactationCalf(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) SearchLactationCalf(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchLactationCalf(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *LactationHandler) AddLactation(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	lac, ok := handlersUtil.DecodeEntity(w, r, &AddLactationStruct{})
	if !ok {
		return
	}

	lac.UserId = userId
	err := h.Repository.AddLactation(lac)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *LactationHandler) UpdateLacAndTransfer(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	lac, ok := handlersUtil.DecodeEntity(w, r, &AddLactationStruct{})
	if !ok {
		return
	}

	lac.UserId = userId
	err := h.Repository.EndLactation(lac)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteUpdateResponse(w)
}

func (h *LactationHandler) UpdateLactation(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	lac, ok := handlersUtil.DecodeEntity(w, r, &LactationHist{})
	if !ok {
		return
	}

	lac.UserId = userId
	res, err := h.Repository.UpdateLactation(lac)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *LactationHandler) DeleteLactation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteLactation(id, userId)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteDeleteResponse(w)
}
