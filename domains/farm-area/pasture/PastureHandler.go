package pasture

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type PastureHandler struct {
	Repository *PastureRepository
}

func (h *PastureHandler) SearchPasture(w http.ResponseWriter, r *http.Request) {
	input := "%" + r.URL.Query().Get("input") + "%"
	farmId := r.URL.Query().Get("farmId")
	farmArray := handlersUtil.ParseArray(farmId)

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	id := r.URL.Query().Get("id")
	if id != "" {
		idList := handlersUtil.ParseArray(id)
		list, err := h.Repository.SearchPastureById(userId, farmArray, idList)
		if err != nil {
			serverErrors.DatabaseGetError(err, w)
			return
		}
		handlersUtil.SendList(w, list)
		return
	}

	list, err := h.Repository.SearchPasture(userId, input, farmArray)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendList(w, list)
}

func (h *PastureHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	handlersUtil.FindAll(w, r, h.Repository)
}

func (h *PastureHandler) FindAnimalsByPasture(w http.ResponseWriter, r *http.Request) {
	pastureId := r.PathValue("id")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}
	result, err := h.Repository.FindAnimalsByPasture(pastureId, userId, sort, order)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendList(w, result)
}

func (h *PastureHandler) FindById(w http.ResponseWriter, r *http.Request) {
	handlersUtil.FindById(w, r, h.Repository)
}

func (h *PastureHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w, r, h.Repository)
}

func (h *PastureHandler) Update(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Update(w, r, h.Repository)
}

func (h *PastureHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Delete(w, r, h.Repository)
}
