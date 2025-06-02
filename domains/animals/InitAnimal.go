package animals

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitAnimal(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := AnimalHandler{repository}

	app.HandleFunc("POST /animals/page", handler.FindPage)
	app.HandleFunc("GET /animals/{id}", handler.FindById)
	app.HandleFunc("GET /animals/name/{name}", handler.FindByName)
	app.HandleFunc("GET /animals/number/{number}", handler.FindByNumber)
	app.HandleFunc("GET /animals/father/{fatherId}", handler.FindByFatherId)
	app.HandleFunc("GET /animals/mother/{motherId}", handler.FindByMotherId)
	app.HandleFunc("GET /animals/pasture/{pastureId}/page", handler.FindByPastureId)
	app.HandleFunc("POST /animals", handler.Add)
	app.HandleFunc("POST /animals/save", handler.Update)
	app.HandleFunc("DELETE /animals/{id}", handler.Delete)
	util.LogDomainsInit("Animais")
}
