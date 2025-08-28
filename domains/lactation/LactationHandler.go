package lactation

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/serverErrors"
	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type LactationHandler struct {
	Repository *LactationRepository
}

func (h *LactationHandler) GetYearlyMilk(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetYearlyMilk(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetMonthMilk(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetMonthMilk(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetAnimalsAverage(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetAnimalsAverage(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetMilkProduction(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetMilkProduction(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetBestAnimals(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetBestAnimals(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetWorstAnimals(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetWorstAnimals(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetLastEntries(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetLastGroups(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) FindGroupsPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")

	filter, ok := handlersUtil.DecodeFilter(w, r, LactationGroupFilter{})
	if !ok {
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindGroupsPage(filter, order, cursor, userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	order := r.URL.Query().Get("order")
	sort := r.URL.Query().Get("sort")

	filter, ok := handlersUtil.DecodeFilter(w, r, MilkEntryFilter{})
	if !ok {
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.FindEntriesPage(filter, sort, order, cursor, userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetEntriesPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, ok := handlersUtil.DecodeFilter(w, r, MilkEntryFilter{})
	if !ok {
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetEntriesPageFoot(filter, userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetGroupEntries(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.URL.Query().Get("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateVar)

	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetGroupEntries(userId, entryDate)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetGroupEntriesFoot(w http.ResponseWriter, r *http.Request) {
	entryDateVar := r.URL.Query().Get("entryDate")
	entryDate, err := time.Parse(time.RFC3339Nano, entryDateVar)

	userId, ok := handlersUtil.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Repository.GetGroupEntriesFoot(userId, entryDate)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}
