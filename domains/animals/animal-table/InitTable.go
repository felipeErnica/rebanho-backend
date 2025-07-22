package animalTable

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitTable(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := AnimalHandler{repository}

	app.HandleFunc("POST /animals/info/page", handler.FindPage)
	app.HandleFunc("GET /animals/info/id/{id}", handler.FindById)
	app.HandleFunc("GET /animals/info/name/{name}", handler.FindByName)
	app.HandleFunc("GET /animals/info/number/{number}", handler.FindByNumber)
	app.HandleFunc("GET /animals/info/father/{fatherId}", handler.FindByFatherId)
	app.HandleFunc("GET /animals/info/mother/{motherId}", handler.FindByMotherId)
	app.HandleFunc("GET /animals/info/pasture/{pastureId}/page", handler.FindByPastureId)
	app.HandleFunc("GET /animals/info/search/father", handler.SearchFather)
	app.HandleFunc("GET /animals/info/search/mother", handler.SearchMother)
	app.HandleFunc("GET /animals/info/search/bull", handler.SearchBull)
	app.HandleFunc("GET /animals/info/search/animal", handler.SearchAnimal)
	app.HandleFunc("POST /animals/info/add", handler.Add)
	app.HandleFunc("POST /animals/info/save", handler.Update)
	app.HandleFunc("DELETE /animals/info/{id}", handler.Delete)
	util.LogDomainsInit("Animais")
}
