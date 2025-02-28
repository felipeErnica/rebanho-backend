package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/repositories"
)

func InitLactation(mux *http.ServeMux) {
    handler:=LactationHandler{
        Repository: repositories.LactationRepository{},
    }
    
    mux.HandleFunc("GET /lactation/firstPage", handler.GetFirstPage)
    LogControllersInit("Lactações")
}

type LactationHandler struct {
	Repository repositories.LactationRepository
}

func (l *LactationHandler) GetFirstPage(w http.ResponseWriter, r *http.Request) {
    page, err:= l.Repository.GetFirstPage()
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(page)
    if err != nil {
        JsonServerError(err, w)
    }

    w.Write(response)
}

