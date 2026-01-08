package animals

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitAnimal(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := AnimalHandler{repository}

	app.HandleFunc("GET /animals/stats/dairy", handler.GetDairyHist)
	app.HandleFunc("GET /animals/stats/birth", handler.GetBirthHist)
	app.HandleFunc("GET /animals/stats/death", handler.GetDeathHist)
	app.HandleFunc("GET /animals/stats/slaughter", handler.GetSlaughterHist)
	app.HandleFunc("GET /animals/stats/types", handler.GetAnimalTypes)
	app.HandleFunc("GET /animals/stats/last-deaths", handler.GetLastDeaths)
    app.HandleFunc("GET /animals/stats/age-and-sex", handler.GetAgeAndSex)

	app.HandleFunc("GET /animals/{id}", handler.FindById)
	app.HandleFunc("POST /animals/page/foot", handler.GetPageFoot)
	app.HandleFunc("POST /animals/page", handler.FindPage)
	app.HandleFunc("POST /animals/search", handler.Search)

	app.HandleFunc("DELETE /animals/{id}", handler.Delete)
	app.HandleFunc("PUT /animals/{id}", handler.Update)
	app.HandleFunc("POST /animals", handler.Add)

	util.LogDomainsInit("Animais")
}
