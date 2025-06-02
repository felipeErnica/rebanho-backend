package birth

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitBirth(app *app.App) {
    repository := NewRepository(app.DBconn)
    handler := BirthHandler{repository}
    app.HandleFunc("GET /reproduction/births/mother/{motherId}", handler.FindByMotherId)
    app.HandleFunc("DELETE /reproduction/births/delete/{id}", handler.Delete)
        util.LogDomainsInit("Nascimentos")
}
