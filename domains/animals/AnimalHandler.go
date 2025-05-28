package animals

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/serverErrors"
)

type AnimalHandler struct {
	Repository *AnimalRepository
}

func (h *AnimalHandler) FindPage(w http.ResponseWriter, r *http.Request) {
	var filter entity.AnimalFilter

	err := json.NewDecoder(r.Body).Decode(&filter)
	if err != nil {
		err = errors.New(fmt.Sprintf("Falha na decodificação do filtro: %s", err.Error()))
		serverErrors.JsonServerError(err, w)
		return
	}

	cursor, sort, order := h.Impl.GetPageParameters(r)
	page, err := h.Repository.FindPage(sort, order, cursor, &filter)
	h.Impl.SendPage(w, page, err)
}

func (h *AnimalHandler) FindById(w http.ResponseWriter, r *http.Request) {
	h.Impl.FindById(w, r)
}

func (h *AnimalHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	animals, err := h.Repository.FindByName(name)
	h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByNumber(w http.ResponseWriter, r *http.Request) {
	number := r.PathValue("number")
	animals, err := h.Repository.FindByIdentificationNumber(number)
	h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByFatherId(w http.ResponseWriter, r *http.Request) {
	fatherId := r.PathValue("fatherId")
	animals, err := h.Repository.FindByFatherId(fatherId)
	h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByMotherId(w http.ResponseWriter, r *http.Request) {
	motherId := r.PathValue("motherId")
	animals, err := h.Repository.FindByMotherId(motherId)
	h.Impl.SendList(w, animals, err)
}

func (h *AnimalHandler) FindByPastureId(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("pastureId")
	cursor, sort, order := h.Impl.GetPageParameters(r)
	pasturePage, err := h.Repository.FindByPastureId(sort, order, cursor, pastureId)
	h.Impl.SendPage(w, pasturePage, err)
}

func (h *AnimalHandler) Add(w http.ResponseWriter, r *http.Request) {
	h.Impl.Add(w, r)
}

func (h *AnimalHandler) Save(w http.ResponseWriter, r *http.Request) {
	h.Impl.Save(w, r)
}

func (h *AnimalHandler) Delete(w http.ResponseWriter, r *http.Request) {
	h.Impl.Delete(w, r)
}
