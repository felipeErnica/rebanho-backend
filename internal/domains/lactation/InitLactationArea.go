package lactation

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitLactationArea(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewLactationService(repository)
	handler := LactationHandler{service}

	app.HandleFunc("GET /lactation/stats/last-lactating", handler.GetLastLactating)
	app.HandleFunc("GET /lactation/stats/last-dry", handler.GetLastDry)
	app.HandleFunc("GET /lactation/stats/dairy-types", handler.GetDairyTypes)
	app.HandleFunc("GET /lactation/stats/best-animals", handler.GetBestAnimals)
	app.HandleFunc("GET /lactation/stats/worst-animals", handler.GetWorstAnimals)
	app.HandleFunc("GET /lactation/stats/best-mothers", handler.GetBestMothers)
	app.HandleFunc("GET /lactation/stats/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /lactation/stats/worst-mothers", handler.GetWorstMothers)
	app.HandleFunc("GET /lactation/stats/worst-fathers", handler.GetWorstFathers)

	app.HandleFunc("GET /lactation/page", handler.FindLactationPage)
	app.HandleFunc("GET /lactation/page/foot", handler.GetLactationPageFoot)
	app.HandleFunc("GET /lactation/animals/page", handler.FindAnimalsPage)
	app.HandleFunc("GET /lactation/animals/page/foot", handler.GetAnimalsPageFoot)
	app.HandleFunc("GET /lactation/{id}", handler.FindById)

	app.HandleFunc("POST /lactation", handler.AddLactation)
	app.HandleFunc("PUT /lactation", handler.UpdateLactation)
	app.HandleFunc("DELETE /lactation/{id}", handler.DeleteLactation)
	log.LogDomainsInit("Lactações")
}
