package testEntries

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type TestEntryHandler struct {
	Repository *TestEntryRepository
}

func (h *TestEntryHandler) FindByGroupId(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("groupId")
	list, err := h.Repository.FindByGroupId(groupId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendList(w, list)
}

func (h *TestEntryHandler) FindByAnimalId(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("animalId")
	list, err := h.Repository.FindByAnimalId(groupId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendList(w, list)
}

func (h *TestEntryHandler) FindById(w http.ResponseWriter, r *http.Request) {
	handlersUtil.FindById(w, r, h.Repository)
}

func (h *TestEntryHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w, r, h.Repository)
}

func (h *TestEntryHandler) Update(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Update(w, r, h.Repository)
}

func (h *TestEntryHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Delete(w, r, h.Repository)
}
