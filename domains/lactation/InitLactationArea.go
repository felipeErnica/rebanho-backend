package lactation

import (
	"github.com/felipeErnica/rebanho-backend/app"
)

func InitLactationArea(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := LactationHandler{repository}

    app.HandleFunc("GET /lactation/dashboard/yearly-milk", handler.GetYearlyMilk)
    app.HandleFunc("GET /lactation/dashboard/month-milk", handler.GetMonthMilk)
    app.HandleFunc("GET /lactation/dashboard/animals-average", handler.GetAnimalsAverage)
    app.HandleFunc("GET /lactation/dashboard/milk-production", handler.GetMilkProduction)
    app.HandleFunc("GET /lactation/dashboard/best-animals", handler.GetBestAnimals)
    app.HandleFunc("GET /lactation/dashboard/worst-animals", handler.GetWorstAnimals)
    app.HandleFunc("GET /lactation/dashboard/last-entries", handler.GetLastEntries)
    app.HandleFunc("GET /lactation/dashboard/last-groups", handler.GetLastGroups)

    app.HandleFunc("POST /lactation/groups/page", handler.FindGroupsPage)
    app.HandleFunc("GET /lactation/groups/entries", handler.GetGroupEntries)
    app.HandleFunc("GET /lactation/groups/entries/foot", handler.GetGroupEntriesFoot)

    app.HandleFunc("POST /lactation/entries/page", handler.FindEntriesPage)
    app.HandleFunc("POST /lactation/entries/page/foot", handler.GetEntriesPageFoot)
}
