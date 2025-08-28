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
	farmId := r.URL.Query().Get("farmId")
	farmArray := handlersUtil.ParseArray(farmId)

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	list, err := h.Repository.SearchPasture(userId, farmArray)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendList(w, list)
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
