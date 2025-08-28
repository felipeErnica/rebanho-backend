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
	app.HandleFunc("GET /reproduction/births/dashboard/best-intervals", handler.GetBestIntervals)
	app.HandleFunc("GET /reproduction/births/dashboard/birth-stats", handler.GetBirthStats)
	app.HandleFunc("GET /reproduction/births/dashboard/total-sex", handler.TotalBySex)
	util.LogDomainsInit("Nascimentos")
}
