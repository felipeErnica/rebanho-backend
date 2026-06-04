package slaughtergroups

import (
	"net/http"

	"github.com/felipeErnica/rebanho-backend/internal/log"
	"github.com/felipeErnica/rebanho-backend/internal/util"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service}
}

func (h *Handler) FindAll(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	response, err := h.Service.FindAll(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, response)
}

func (h *Handler) FindLast(w http.ResponseWriter, r *http.Request) {
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	response, err := h.Service.FindLast(userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, response)
}

func (h *Handler) FindById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userId, ok := util.GetUserId(w, r)
	if !ok {
		return
	}

	response, err := h.Service.FindById(id, userId)
	if err != nil {
		log.WriteError(w, err)
		return
	}

	util.WriteEntity(w, response)
}
