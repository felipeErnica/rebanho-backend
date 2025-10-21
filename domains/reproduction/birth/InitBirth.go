package birth

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitBirth(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := BirthHandler{repository}
	app.HandleFunc("POST /reproduction/births/table/page/footer", handler.FindPageFooter)
	app.HandleFunc("POST /reproduction/births/table/page", handler.FindPage)

	app.HandleFunc("GET /reproduction/births/dashboard/last-births", handler.GetLastBirths)
	app.HandleFunc("GET /reproduction/births/dashboard/births-number", handler.GetLastBirthsNumber)
	app.HandleFunc("GET /reproduction/births/dashboard/year-births", handler.GetYearBirthsNumber)
	app.HandleFunc("GET /reproduction/births/dashboard/year-deaths", handler.GetYearDeathsNumber)
	app.HandleFunc("GET /reproduction/births/dashboard/year-sex", handler.GetYearBySex)
	app.HandleFunc("GET /reproduction/births/dashboard/best-intervals", handler.GetBestIntervals)
	app.HandleFunc("GET /reproduction/births/dashboard/worst-intervals", handler.GetWorstIntervals)
	app.HandleFunc("GET /reproduction/births/dashboard/interval-stats", handler.GetBirthIntervalHistory)
	app.HandleFunc("GET /reproduction/births/dashboard/birth-history", handler.GetBirthHistory)
	app.HandleFunc("GET /reproduction/births/dashboard/death-index", handler.GetDeathIndex)
	app.HandleFunc("GET /reproduction/births/dashboard/total-sex", handler.TotalBySex)
	util.LogDomainsInit("Nascimentos")
}
