package lactation

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitLactationArea(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := LactationHandler{repository}

	app.HandleFunc("GET /lactation/stats/long-lactations", handler.GetLongLactations)
	app.HandleFunc("GET /lactation/stats/last-milk", handler.GetLastMilk)
	app.HandleFunc("GET /lactation/stats/last-avg-milk", handler.GetLastAverageMilk)
	app.HandleFunc("GET /lactation/stats/last-lactating", handler.GetLastLactating)
	app.HandleFunc("GET /lactation/stats/last-dry", handler.GetLastDry)
	app.HandleFunc("GET /lactation/stats/milk-production", handler.GetMilkProduction)
	app.HandleFunc("GET /lactation/stats/dairy-types", handler.GetDairyTypes)
	app.HandleFunc("GET /lactation/stats/year-milk", handler.GetYearMilkProduction)
	app.HandleFunc("GET /lactation/stats/year-avg-milk", handler.GetYearAverageMilk)
	app.HandleFunc("GET /lactation/stats/best-animals", handler.GetBestAnimals)
	app.HandleFunc("GET /lactation/stats/worst-animals", handler.GetWorstAnimals)
	app.HandleFunc("GET /lactation/stats/best-mothers", handler.GetBestMothers)
	app.HandleFunc("GET /lactation/stats/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /lactation/stats/worst-mothers", handler.GetWorstMothers)
	app.HandleFunc("GET /lactation/stats/worst-fathers", handler.GetWorstFathers)
	app.HandleFunc("GET /lactation/stats/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /lactation/stats/last-groups", handler.GetLastGroups)

	app.HandleFunc("GET /lactation/page", handler.FindLactationPage)
	app.HandleFunc("GET /lactation/page/foot", handler.GetLactationPageFoot)
	app.HandleFunc("GET /lactation/dry-animals/page", handler.FindDryAnimalsPage)
	app.HandleFunc("GET /lactation/dry-animals/page/foot", handler.GetDryAnimalsFoot)
	app.HandleFunc("GET /lactation/lac-animals/page", handler.FindLacAnimalsPage)
	app.HandleFunc("GET /lactation/lac-animals/page/foot", handler.GetLacAnimalsFoot)
	app.HandleFunc("GET /lactation/{id}", handler.FindById)
	app.HandleFunc("GET /lactation/{id}/entries", handler.GetLactationEntries)
	app.HandleFunc("GET /lactation/{id}/entries/foot", handler.GetLactationEntriesFoot)
	app.HandleFunc("GET /lactation/search-lactating", handler.SearchLactatingAnimals)
	app.HandleFunc("GET /lactation/search-dry", handler.SearchDryAnimals)
	app.HandleFunc("GET /lactation/search-calfs-new", handler.SearchNewLactationCalf)
	app.HandleFunc("GET /lactation/search-calfs", handler.SearchLactationCalf)


	app.HandleFunc("POST /lactation", handler.AddLactation)
	app.HandleFunc("PUT /lactation", handler.UpdateLactation)
	app.HandleFunc("DELETE /lactation/{id}/delete", handler.DeleteLactation)
	util.LogDomainsInit("Lactações")
}

func InitMilk(app *app.App) {
	repository := NewMilkRepository(app.DBconn)
	handler := MilkHandler{repository}

	app.HandleFunc("GET /lactation/groups/page", handler.FindGroupsPage)
	app.HandleFunc("GET /lactation/groups/{entryDate}/entries", handler.GetGroupEntries)
	app.HandleFunc("GET /lactation/groups/{entryDate}/entries/foot", handler.GetGroupEntriesFoot)

	app.HandleFunc("PUT /lactation/groups/{entryDate}/update", handler.UpdateGroup)
	app.HandleFunc("DELETE /lactation/groups/{entryDate}/delete", handler.DeleteGroup)

	app.HandleFunc("GET /lactation/entries/page", handler.FindEntriesPage)
	app.HandleFunc("GET /lactation/entries/page/foot", handler.GetEntriesPageFoot)

	app.HandleFunc("PUT /lactation/entries/add", handler.Add)
	app.HandleFunc("DELETE /lactation/entries/{id}/delete", handler.Delete)
	app.HandleFunc("PUT /lactation/entries/update", handler.Update)
	app.HandleFunc("PUT /lactation/entries/replace", handler.Replace)
	util.LogDomainsInit("Entrada de Leite")
}
