package animals

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitAnimal(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := AnimalHandler{repository}

	app.HandleFunc("GET /animals/dashboard/dairy-hist", handler.GetDairyHist)
	app.HandleFunc("GET /animals/dashboard/birth-hist", handler.GetBirthHist)
	app.HandleFunc("GET /animals/dashboard/death-hist", handler.GetDeathHist)
	app.HandleFunc("GET /animals/dashboard/slaughter-hist", handler.GetSlaughterHist)
	app.HandleFunc("GET /animals/dashboard/animal-types", handler.GetAnimalTypes)





//-------------------------------------------- Links Legados ------------------------------------------------------------------//
    app.HandleFunc("POST /animals/dashboard/total-general", handler.TotalBySex)
    app.HandleFunc("POST /animals/dashboard/types", handler.TotalByType)
    app.HandleFunc("POST /animals/dashboard/group-age-farm", handler.GroupByAgeAndFarm)
    app.HandleFunc("POST /animals/dashboard/group-pasture", handler.GroupByAgeAndPasture)
    app.HandleFunc("POST /animals/dashboard/group-age", handler.GroupByAge)
    app.HandleFunc("POST /animals/dashboard/group-year", handler.GroupByYear)

	app.HandleFunc("POST /animals/info/page", handler.FindPage)
	app.HandleFunc("GET /animals/info/id/{id}", handler.FindById)
	app.HandleFunc("GET /animals/info/name/{name}", handler.FindByName)
	app.HandleFunc("GET /animals/info/number/{number}", handler.FindByNumber)
	app.HandleFunc("GET /animals/info/father/{fatherId}", handler.FindByFatherId)
	app.HandleFunc("GET /animals/info/mother/{motherId}", handler.FindByMotherId)

	app.HandleFunc("DELETE /animals/delete/{id}", handler.DeleteAnimal)
	app.HandleFunc("DELETE /animals/delete-no-validation/{id}", handler.DeleteNoValidation)
	app.HandleFunc("PUT /animals/update", handler.Update)
	app.HandleFunc("PUT /animals/update-no-validation", handler.UpdateNoValidation)
	app.HandleFunc("PUT /animals/add", handler.Add)
	app.HandleFunc("PUT /animals/add-no-validation", handler.AddNoValidation)
	app.HandleFunc("PUT /animals/replace", handler.Replace)

	app.HandleFunc("GET /animals/{id}", handler.FindById)
	app.HandleFunc("GET /animals/{id}/male-offspring", handler.FindMaleOffspring)
	app.HandleFunc("GET /animals/{id}/female-offspring", handler.FindFemaleOffspring)

	app.HandleFunc("GET /animals/info/search/father", handler.SearchFather)
	app.HandleFunc("GET /animals/info/search/mother-all", handler.SearchAllMothers)
	app.HandleFunc("GET /animals/info/search/mother", handler.SearchMother)
	app.HandleFunc("GET /animals/info/search/bull", handler.SearchBull)
	app.HandleFunc("GET /animals/info/search/animal", handler.SearchAnimal)
	app.HandleFunc("GET /animals/info/search/dairy-animal", handler.SearchDairyAnimal)
//-------------------------------------------------------------------------------------------------------------------------------------------------//



	util.LogDomainsInit("Animais")
}
