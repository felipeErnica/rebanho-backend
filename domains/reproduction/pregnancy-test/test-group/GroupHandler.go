package testGroup

import (
	"net/http"

	handlersUtil "github.com/felipeErnica/rebanho-backend/util/handlers-util"
)

type TestGroupHandler struct {
    Repository *TestGroupRepository
}

func (h *TestGroupHandler) FindAll(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindAll(w, r, h.Repository)
} 

func (h *TestGroupHandler) FindById(w http.ResponseWriter, r *http.Request) {
    handlersUtil.FindById(w,r, h.Repository)
} 

func (h *TestGroupHandler) Add(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Add(w,r, h.Repository)
} 

func (h *TestGroupHandler) Update(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Update(w,r, h.Repository)
} 

func (h *TestGroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
    handlersUtil.Delete(w, r, h.Repository)
} 
