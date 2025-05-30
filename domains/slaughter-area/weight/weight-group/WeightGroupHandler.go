package weightGroup

import (
	"net/http"

	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type WeightGroupHandler struct {
	Repository *WeightGroupRepository
}

func (h *WeightGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	handlersUtil.FindAll(w, r, h.Repository)
}

func (h *WeightGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
	handlersUtil.FindById(w, r, h.Repository)
}

func (h *WeightGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w, r, h.Repository)
}

func (h *WeightGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Update(w, r, h.Repository)
}

func (h *WeightGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
}
