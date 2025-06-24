package farm

import (
	"net/http"

	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type FarmHandler struct {
	Repository *FarmRepository
}

func (h *FarmHandler) SearchFarm(w http.ResponseWriter, r *http.Request) {
    handlersUtil.ReturnSearchResults(w, r, h.Repository.SearchFarm)
}

func (h *FarmHandler) Add(w http.ResponseWriter, r *http.Request) {
	handlersUtil.Add(w, r, h.Repository)
}
