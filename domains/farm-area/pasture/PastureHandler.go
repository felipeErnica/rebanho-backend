package pasture

import (
	"net/http"

	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type PastureHandler struct {
	Repository *PastureRepository
}

func (h *PastureHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindAll(w, r, h.Repository)
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
