package inseminationGroup

import (
	"net/http"

	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type GroupHandler struct {
	Repository *GroupRepository
}

func (h *GroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindAll(w, r, h.Repository)
}

func (h *GroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById(w, r, h.Repository)
}

func (h *GroupHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w, r, h.Repository)
}

func (h *GroupHandler) Save(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w, r, h.Repository)
}

func (h *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
