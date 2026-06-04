package birth

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitBirth(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewBirthService(repository)
	handler := BirthHandler{service}

	app.HandleFunc("GET /births/page", handler.FindPage)
	app.HandleFunc("GET /births/page/foot", handler.GetPageFoot)

	app.HandleFunc("POST /births", handler.AddBirth)
	app.HandleFunc("PUT /births", handler.UpdateBirth)
	app.HandleFunc("DELETE /births", handler.DeleteBirth)

	app.HandleFunc("GET /births/{id}", handler.GetById)
	app.HandleFunc("GET /births/potential-father", handler.GetPotentialFather)

	app.HandleFunc("GET /births/stats/last-births", handler.GetLastBirths)
	app.HandleFunc("GET /births/stats/births-number", handler.GetLastBirthsNumber)
	app.HandleFunc("GET /births/stats/year-births", handler.GetYearBirthsNumber)
	app.HandleFunc("GET /births/stats/year-deaths", handler.GetYearDeathsNumber)
	app.HandleFunc("GET /births/stats/year-sex", handler.GetYearBySex)
	app.HandleFunc("GET /births/stats/best-intervals", handler.GetBestIntervals)
	app.HandleFunc("GET /births/stats/worst-intervals", handler.GetWorstIntervals)
	app.HandleFunc("GET /births/stats/interval-stats", handler.GetBirthIntervalHistory)
	app.HandleFunc("GET /births/stats/birth-history", handler.GetBirthHistory)
	app.HandleFunc("GET /births/stats/death-index", handler.GetDeathIndex)
	app.HandleFunc("GET /births/stats/total-sex", handler.TotalBySex)
	log.LogDomainsInit("Nascimentos")
}
