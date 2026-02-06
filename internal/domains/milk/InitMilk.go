package milk

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitMilk(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewService(repository)
	handler := MilkHandler{service}

	app.HandleFunc("GET /milk/stats/last-milk", handler.GetLastMilk)
	app.HandleFunc("GET /milk/stats/last-avg-milk", handler.GetLastAverageMilk)
	app.HandleFunc("GET /milk/stats/milk-production", handler.GetMilkProduction)
	app.HandleFunc("GET /milk/stats/year-milk", handler.GetYearMilkProduction)
	app.HandleFunc("GET /milk/stats/year-avg-milk", handler.GetYearAverageMilk)
	app.HandleFunc("GET /milk/stats/last-entries", handler.GetLastEntries)
	app.HandleFunc("GET /milk/stats/last-groups", handler.GetLastGroups)

	app.HandleFunc("GET /lactation/{id}/entries", handler.GetLactationEntries)
	app.HandleFunc("GET /lactation/{id}/entries/foot", handler.GetLactationEntriesFoot)

	app.HandleFunc("GET /milk/groups/page", handler.FindGroupsPage)
	app.HandleFunc("GET /milk/groups/{entryDate}/entries", handler.GetGroupEntries)
	app.HandleFunc("GET /milk/groups/{entryDate}/entries/foot", handler.GetGroupEntriesFoot)

	app.HandleFunc("PUT /milk/groups", handler.UpdateGroup)
	app.HandleFunc("DELETE /milk/groups/{entryDate}", handler.DeleteGroup)

	app.HandleFunc("GET /milk/page", handler.FindEntriesPage)
	app.HandleFunc("GET /milk/page/foot", handler.GetEntriesPageFoot)

	app.HandleFunc("POST /milk", handler.Add)
	app.HandleFunc("PUT /milk", handler.Update)
	app.HandleFunc("DELETE /milk/{id}", handler.Delete)
	log.LogDomainsInit("Entrada de Leite")
}
