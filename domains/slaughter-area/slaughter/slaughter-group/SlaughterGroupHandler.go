package slaughterGroup

import (
	"net/http"

	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type SlaughterGroupHandler struct {
	Repository *SlaughterGroupRepository
}

func (h *SlaughterGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindAll(w, r, h.Repository)
}

func (h *SlaughterGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById(w, r, h.Repository)
}

func (h *SlaughterGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w, r, h.Repository)
}

func (h *SlaughterGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w, r, h.Repository)
}

func (h *SlaughterGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
