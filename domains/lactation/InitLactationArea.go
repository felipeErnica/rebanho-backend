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
}
