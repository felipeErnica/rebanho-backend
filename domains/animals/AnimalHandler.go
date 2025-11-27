package animals

import (
	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
	"net/http"
)

type AnimalHandler struct {
	Repository *AnimalRepository
}

func (h *AnimalHandler) GroupByYear(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GroupByYear(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) TotalBySex(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GroupByYear(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) TotalByType(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.TotalByType(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GroupByAgeAndFarm(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GroupByAgeAndFarm(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GroupByAgeAndPasture(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GroupByAgeAndPasture(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, res)
}

func (h *AnimalHandler) GroupByAge(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	res, err := h.Repository.GroupByAge(userId)
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

func (h *AnimalHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}
	animals, err := h.Repository.FindByName(name, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindByNumber(w http.ResponseWriter, r *http.Request) {
	number := r.PathValue("number")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}
	animals, err := h.Repository.FindByNumber(number, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}
	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindByFatherId(w http.ResponseWriter, r *http.Request) {
	fatherId := r.PathValue("fatherId")
	animals, err := h.Repository.FindByFatherId(fatherId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindByMotherId(w http.ResponseWriter, r *http.Request) {
	motherId := r.PathValue("motherId")
	animals, err := h.Repository.FindByMotherId(motherId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindMaleOffspring(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindMaleOffspring(id, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) FindFemaleOffspring(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindFemaleOffspring(id, userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) SearchFather(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchFather(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}


func (h *AnimalHandler) SearchAnimal(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) SearchAllMothers(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchAllMothers(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) SearchMother(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchMother(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) SearchBull(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchBull(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) SearchDairyAnimal(w http.ResponseWriter, r *http.Request) {

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.SearchDairyAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) DeleteAnimal(w http.ResponseWriter, r *http.Request) {
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

func (h *AnimalHandler) DeleteNoValidation(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Repository.DeleteNoValidation(id, userId)
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

func (h *AnimalHandler) UpdateNoValidation(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	newEntry, ok := handlersUtil.DecodeEntity(w, r, &AnimalSave{})
	if !ok {
		return
	}

	newEntry.UserId = userId
	response, err := h.Repository.UpdateNoValidation(newEntry)
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

func (h *AnimalHandler) AddNoValidation(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	newEntry, ok := handlersUtil.DecodeEntity(w, r, &AnimalSave{})
	if !ok {
		return
	}

	newEntry.UserId = userId
	err := h.Repository.AddNoValidation(newEntry)
	if err != nil {
		apiError.WriteAPIError(err, w)
		return
	}

	handlersUtil.WriteCreatedResponse(w)
}

func (h *AnimalHandler) Replace(w http.ResponseWriter, r *http.Request) {
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

	handlersUtil.WriteUpdateResponse(w)
}
