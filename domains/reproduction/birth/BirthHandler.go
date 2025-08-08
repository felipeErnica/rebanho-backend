package birth

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type BirthHandler struct {
	Repository *BirthRepository
}

func (h *BirthHandler) FindPageFooter(w http.ResponseWriter, r *http.Request) {
    filter, ok := handlersUtil.DecodeFilter(w, r, BirthEntryFilter{}); if !ok {
        return
    }

    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.FindPageFooter(userId, filter)
    if err != nil {
        serverErrors.DatabaseGetError(err , w)
        return
    }
    handlersUtil.SendEntity(w, result)
}

func (h *BirthHandler) FindPage(w http.ResponseWriter, r *http.Request) {
    sort := r.URL.Query().Get("sort")
    order := r.URL.Query().Get("order")
    cursor := r.URL.Query().Get("cursor")

    filter, ok := handlersUtil.DecodeFilter(w, r, BirthEntryFilter{}); if !ok {
        return
    }

    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.FindPage(userId, sort, order, filter, cursor)
    if err != nil {
        serverErrors.DatabaseGetError(err , w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *BirthHandler) GetBestIntervals(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBestIntervals(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendList(w, result)
}

func (h *BirthHandler) GetBirthStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetBirthsStats(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendEntity(w, result)
}

func (h *BirthHandler) TotalBySex(w http.ResponseWriter, r *http.Request) {
	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.TotalBySex(userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendEntity(w, result)
}

func (h *BirthHandler) FindByMotherId(w http.ResponseWriter, r *http.Request) {
	motherId := r.PathValue("motherId")
	list, err := h.Repository.FindByMotherId(motherId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	handlersUtil.SendList(w, list)
}

func (h *BirthHandler) Delete(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Delete(w, r, h.Repository)
}
