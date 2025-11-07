package animalTable

import (
	"net/http"
	"github.com/felipeErnica/rebanho-backend/apiError"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type AnimalHandler struct {
	Repository *AnimalRepository
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
	handlersUtil.FindById(w, r, h.Repository)
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

func (h *AnimalHandler) SearchFather(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
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
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}


	result, err := h.Repository.SearchAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}

func (h *AnimalHandler) SearchMother(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
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
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
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

	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}


	result, err := h.Repository.SearchDairyAnimals(userId)
	if err != nil {
		apiError.WriteError(err, w)
		return
	}

	handlersUtil.WriteEntity(w, result)
}
