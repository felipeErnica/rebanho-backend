package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
)

func InitLactation(mux *http.ServeMux) {
    handler:=LactationHandler{
        Repository: repositories.LactationRepository{},
    }
    
    mux.HandleFunc("GET /lactation/page", handler.GetPage)
    LogControllersInit("Lactações")
}

type LactationHandler struct {
	Repository repositories.LactationRepository
}

func (l *LactationHandler) GetPage(w http.ResponseWriter, r *http.Request) {
    cursor:=r.URL.Query().Get("cursor")
    sort:=r.URL.Query().Get("sort")
    direction:=r.URL.Query().Get("order")
    
    var page *entity.LactationPage
    var err error

    if cursor == "" {
        page, err = l.Repository.GetFirstPage(sort, direction)
    } else {
        page, err = l.Repository.GetNextPage(cursor, sort, direction)
    }

    
    if err != nil {
        DatabaseSendError(err, w)
    }
    
    response, err:= json.Marshal(page)
    if err != nil {
        JsonServerError(err, w)
    }

    w.Write(response)
}

