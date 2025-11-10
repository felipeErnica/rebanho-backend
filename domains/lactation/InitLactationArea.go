package lactation

import (
	"github.com/felipeErnica/rebanho-backend/app"
)

func InitLactationArea(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := LactationHandler{repository}

    app.HandleFunc("GET /lactation/dashboard/last-milk", handler.GetLastMilk)
    app.HandleFunc("GET /lactation/dashboard/last-avg-milk", handler.GetLastAverageMilk)
    app.HandleFunc("GET /lactation/dashboard/last-count", handler.GetLastAnimalsCount)
    app.HandleFunc("GET /lactation/dashboard/milk-production", handler.GetMilkProduction)
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

    app.HandleFunc("POST /lactation/groups/page", handler.FindGroupsPage)
    app.HandleFunc("GET /lactation/groups/entries", handler.GetGroupEntries)
    app.HandleFunc("GET /lactation/groups/entries/foot", handler.GetGroupEntriesFoot)

    app.HandleFunc("POST /lactation/entries/page", handler.FindEntriesPage)
    app.HandleFunc("POST /lactation/entries/page/foot", handler.GetEntriesPageFoot)
    app.HandleFunc("PUT /lactation/entries/add", handler.AddMilkEntry)
    app.HandleFunc("PUT /lactation/entries/add-and-transfer", handler.AddMilkAndTransferPasture)
    app.HandleFunc("PUT /lactation/entries/add-no-transfer", handler.AddMilkEntryNoTransfer)
    app.HandleFunc("DELETE /lactation/entries/delete/{id}", handler.DeleteMilkEntry)
    app.HandleFunc("PUT /lactation/entries/update", handler.UpdateMilkEntry)
    app.HandleFunc("PUT /lactation/entries/replace", handler.ReplaceMilkEntry)

    app.HandleFunc("POST /lactation/lac-hist/page", handler.FindLactationPage)
    app.HandleFunc("POST /lactation/lac-hist/page/foot", handler.GetLactationPageFoot)
    app.HandleFunc("GET /lactation/lac-hist/{id}/entries", handler.GetLactationEntries)
    app.HandleFunc("GET /lactation/lac-hist/{id}/entries/foot", handler.GetLactationEntriesFoot)
    app.HandleFunc("GET /lactation/lac-hist/search-lactating", handler.SearchLactatingAnimals)
    app.HandleFunc("GET /lactation/lac-hist/search-dry", handler.SearchDryAnimals)
    app.HandleFunc("GET /lactation/lac-hist/search-calfs-new", handler.SearchNewLactationCalf)
    app.HandleFunc("GET /lactation/lac-hist/search-calfs", handler.SearchLactationCalf)
    app.HandleFunc("PUT /lactation/lac-hist/add", handler.AddLactation)
    app.HandleFunc("PUT /lactation/lac-hist/update", handler.UpdateLactation)
    app.HandleFunc("DELETE /lactation/lac-hist/delete/{id}", handler.DeleteLactation)
    app.HandleFunc("DELETE /lactation/lac-hist/delete-entries/{id}", handler.DeleteLactationAndEntries)
}
