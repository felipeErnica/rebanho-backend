package embryoTransfer

import "net/http"

type TransferHandler struct {
	Repository *TransferRepository
}

func (h *TransferHandler) FindAll(w http.ResponseWriter, r *http.Request) {

}
