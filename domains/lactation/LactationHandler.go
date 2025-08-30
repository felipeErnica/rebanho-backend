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

func (h *LactationHandler) GetBestMothers(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetBestMothers(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetWorstMothers(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetWorstMothers(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetBestFathers(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetBestFathers(userId)
    if err != nil {
        serverErrors.DatabaseGetError(err, w)
        return
    }

    handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetWorstFathers(w http.ResponseWriter, r *http.Request) {
    userId, ok := handlersUtil.GetUserId(w, r); if !ok {
        return
    }

    result, err := h.Repository.GetWorstFathers(userId)
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

func (h *LactationHandler) FindLactationPage(w http.ResponseWriter, r *http.Request) {
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")
	cursor := r.URL.Query().Get("cursor")

	filter, ok := handlersUtil.DecodeFilter(w, r, LactationHistFilter{}); if !ok {
		return
	}
	
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.FindLactationPage(filter, sort, order, cursor, userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	
	handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetLactationPageFoot(w http.ResponseWriter, r *http.Request) {
	filter, ok := handlersUtil.DecodeFilter(w, r, LactationHistFilter{}); if !ok {
		return
	}
	
	userId, ok := handlersUtil.GetUserId(w, r); if !ok {
		return
	}

	result, err := h.Repository.GetLactationPageFoot(filter, userId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}
	
	handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetLactationEntries(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")
	result, err := h.Repository.GetLactationEntries(lacId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}

func (h *LactationHandler) GetLactationEntriesFoot(w http.ResponseWriter, r *http.Request) {
	lacId := r.PathValue("id")
	result, err := h.Repository.GetLactationEntriesFoot(lacId)
	if err != nil {
		serverErrors.DatabaseGetError(err, w)
		return
	}

	handlersUtil.SendEntity(w, result)
}
