package animals

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type AnimalHandler struct {
	Repository *AnimalRepository
}

func (h *AnimalHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	filter := AnimalFilter{}
	handlersUtil.DecodeFilter(w, r, &filter)
    handlersUtil.ReturnPage(w, r, h.Repository, filter)
}

func (h *AnimalHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById[Animal](w, r, h.Repository)
}

func (h *AnimalHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	animals, err := h.Repository.FindByName(name)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindByNumber(w http.ResponseWriter, r *http.Request) {
	number := r.PathValue("number")
	animals, err := h.Repository.FindByIdentificationNumber(number)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindByFatherId(w http.ResponseWriter, r *http.Request) {
	fatherId := r.PathValue("fatherId")
	animals, err := h.Repository.FindByFatherId(fatherId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindByMotherId(w http.ResponseWriter, r *http.Request) {
	motherId := r.PathValue("motherId")
	animals, err := h.Repository.FindByMotherId(motherId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

	handlersUtil.SendList(w, animals)
}

func (h *AnimalHandler) FindByPastureId(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
    filter := AnimalFilter{
        IsFiltered: true,
        Pastures: &[]string{pastureId},
    }
    handlersUtil.ReturnPage(w, r, h.Repository, filter)
}

func (h *AnimalHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w, r, h.Repository)
}

func (h *AnimalHandler) Update(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Update(w, r, h.Repository)
}

func (h *AnimalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Delete(w, r, h.Repository)
}
