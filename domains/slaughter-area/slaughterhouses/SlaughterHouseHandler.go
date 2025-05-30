package slaughterhouses

import (
	"net/http"

	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type SlaughterhouseHandler struct {
	Repository *SlaughterhouseRepository
}

func (h *SlaughterhouseHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindAll(w, r, h.Repository)
}

func (h *SlaughterhouseHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById(w, r, h.Repository)
}

func (h *SlaughterhouseHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w, r, h.Repository)
}

func (h *SlaughterhouseHandler) Update(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w, r, h.Repository)
}

func (h *SlaughterhouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
