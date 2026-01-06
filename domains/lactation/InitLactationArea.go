package lactation

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitLactationArea(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := LactationHandler{repository}

	app.HandleFunc("GET /lactation/dashboard/last-milk", handler.GetLastMilk)
	app.HandleFunc("GET /lactation/dashboard/last-avg-milk", handler.GetLastAverageMilk)
	app.HandleFunc("GET /lactation/dashboard/last-lactating", handler.GetLastLactating)
	app.HandleFunc("GET /lactation/dashboard/last-dry", handler.GetLastDry)
	app.HandleFunc("GET /lactation/dashboard/milk-production", handler.GetMilkProduction)
	app.HandleFunc("GET /lactation/dashboard/dairy-types", handler.GetDairyTypes)
	app.HandleFunc("GET /lactation/dashboard/year-milk", handler.GetYearMilkProduction)
	app.HandleFunc("GET /lactation/dashboard/year-avg-milk", handler.GetYearAverageMilk)
	app.HandleFunc("GET /lactation/dashboard/best-animals", handler.GetBestAnimals)
	app.HandleFunc("GET /lactation/dashboard/worst-animals", handler.GetWorstAnimals)
	app.HandleFunc("GET /lactation/dashboard/best-mothers", handler.GetBestMothers)
	app.HandleFunc("GET /lactation/dashboard/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /lactation/dashboard/worst-mothers", handler.GetWorstMothers)
	app.HandleFunc("GET /lactation/dashboard/worst-fathers", handler.GetWorstFathers)
	app.HandleFunc("GET /lactation/dashboard/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /lactation/dashboard/last-groups", handler.GetLastGroups)

	app.HandleFunc("POST /lactation/lac-hist/page", handler.FindLactationPage)
	app.HandleFunc("POST /lactation/lac-hist/page/foot", handler.GetLactationPageFoot)
	app.HandleFunc("POST /lactation/lac-hist/dry-animals/page", handler.FindDryAnimalsPage)
	app.HandleFunc("POST /lactation/lac-hist/dry-animals/page/foot", handler.GetDryAnimalsFoot)
	app.HandleFunc("POST /lactation/lac-hist/lac-animals/page", handler.FindLacAnimalsPage)
	app.HandleFunc("POST /lactation/lac-hist/lac-animals/page/foot", handler.GetLacAnimalsFoot)
	app.HandleFunc("GET /lactation/lac-hist/{id}", handler.FindById)
	app.HandleFunc("GET /lactation/lac-hist/{id}/entries", handler.GetLactationEntries)
	app.HandleFunc("GET /lactation/lac-hist/{id}/entries/foot", handler.GetLactationEntriesFoot)
	app.HandleFunc("GET /lactation/lac-hist/search-lactating", handler.SearchLactatingAnimals)
	app.HandleFunc("GET /lactation/lac-hist/search-dry", handler.SearchDryAnimals)
	app.HandleFunc("GET /lactation/lac-hist/search-calfs-new", handler.SearchNewLactationCalf)
	app.HandleFunc("GET /lactation/lac-hist/search-calfs", handler.SearchLactationCalf)


	app.HandleFunc("PUT /lactation/lac-hist/add", handler.AddLactation)
	app.HandleFunc("PUT /lactation/lac-hist/update", handler.UpdateLactation)
	app.HandleFunc("DELETE /lactation/lac-hist/{id}/delete", handler.DeleteLactation)
	util.LogDomainsInit("Lactações")
}

func InitMilk(app *app.App) {
	repository := NewMilkRepository(app.DBconn)
	handler := MilkHandler{repository}

	app.HandleFunc("POST /lactation/groups/page", handler.FindGroupsPage)
	app.HandleFunc("GET /lactation/groups/{entryDate}/entries", handler.GetGroupEntries)
	app.HandleFunc("GET /lactation/groups/{entryDate}/entries/foot", handler.GetGroupEntriesFoot)

	app.HandleFunc("PUT /lactation/groups/{entryDate}/update", handler.UpdateGroup)
	app.HandleFunc("DELETE /lactation/groups/{entryDate}/delete", handler.DeleteGroup)

	app.HandleFunc("POST /lactation/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /lactation/entries/page/foot", handler.GetEntriesPageFoot)

	app.HandleFunc("PUT /lactation/entries/add", handler.Add)
	app.HandleFunc("DELETE /lactation/entries/{id}/delete", handler.Delete)
	app.HandleFunc("PUT /lactation/entries/update", handler.Update)
	app.HandleFunc("PUT /lactation/entries/replace", handler.Replace)
	util.LogDomainsInit("Entrada de Leite")
}
