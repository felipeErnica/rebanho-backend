package animalDashboard

import "github.com/felipeErnica/rebanho-backend/app"

func InitDashboard(app *app.App) {
    repository := NewRepository(app.DBconn)
    handler := DashboardHandler{repository}
    app.HandleFunc("POST /animals/dashboard/total-general", handler.TotalBySex)
    app.HandleFunc("POST /animals/dashboard/group-age", handler.GroupByAgeAndFarm)
}
