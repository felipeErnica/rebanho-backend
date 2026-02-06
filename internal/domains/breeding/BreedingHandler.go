package breeding

import (
	"net/http"
	"time"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type BreedingHandler struct {
	Service *BreedingService
}

func (h *BreedingHandler) GetBirthRateStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetBirthRateStats(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) GetPregnancyRateStats(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetPregnancyRateStats(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) GetInseminationHist(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBreedingStats(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *BreedingHandler) GetAnimalsNumber(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetAnimalsNumber(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) GetFutureBirths(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetFutureBirths(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) GetLastGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetLastGroups(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) GetLastEntries(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.GetLastEntries(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) GetBestBull(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetBestBull(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *BreedingHandler) FindEntriesPage(w http.ResponseWriter, r *http.Request) {
	cursor := r.URL.Query().Get("cursor")
	sort := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	filter, err := util.DecodeFilter(r, BreedingEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindEntriesPage(userId, filter, sort, order, cursor)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) FindEntriesByGroup(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("breedingDate")

	breedingDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindEntriesByGroup(userId, breedingDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *BreedingHandler) FindGroups(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.FindGroups(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.SendList(w, result)
}

func (h *BreedingHandler) GetEntriesByGroupFoot(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("breedingDate")

	breedingDate, err := time.Parse(time.RFC3339Nano, queryDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	result, err := h.Service.Repo.GetEntriesByGroupFoot(userId, breedingDate)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) GetEntriesFoot(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	filter, err := util.DecodeFilter(r, BreedingEntryFilter{})
	if err != nil {
		log.WriteError(w, err)
		return
	}

	result, err := h.Service.Repo.GetEntriesFoot(userId, filter)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) AddBreedingBull(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.Repo.AddBreedingBull(id, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteUpdateResponse(w)
}

func (h *BreedingHandler) AddBreeding(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &BreedingEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	err := h.Service.AddBreeding(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteCreatedResponse(w)
}

func (h *BreedingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ignorePregnancy, parseErr := util.ParseBool(r.URL.Query().Get("ignorePregnancy"))
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	changeFather, parseErr := util.ParseBool(r.URL.Query().Get("changeFather"))
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.Delete(id, ignorePregnancy, changeFather, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}

func (h *BreedingHandler) Update(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	entry, ok := util.DecodeEntity(w, r, &BreedingEntrySave{})
	if !ok {
		return
	}

	entry.UserId = userId
	res, err := h.Service.Update(entry)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, res)
}

func (h *BreedingHandler) UpdateGroup(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("breedingDate")
	breedingDate, parseErr := time.Parse(time.RFC3339Nano, queryDate)
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	group, ok := util.DecodeEntity(w, r, &BreedingGroup{})
	if !ok {
		return
	}

	group.UserId = userId
	result, err := h.Service.UpdateGroup(breedingDate, group)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteEntity(w, result)
}

func (h *BreedingHandler) DeleteGroup(w http.ResponseWriter, r *http.Request) {
	queryDate := r.PathValue("breedingDate")
	breedingDate, parseErr := time.Parse(time.RFC3339Nano, queryDate)
	if parseErr != nil {
		log.WriteError(w, parseErr)
		return
	}

	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	err := h.Service.Repo.DeleteGroup(breedingDate, userId)
	if err != nil {
		log.WriteAPIError(err, w)
		return
	}

	util.WriteDeleteResponse(w)
}
