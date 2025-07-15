package dashboard

import "github.com/felipeErnica/rebanho-backend/app"

func InitFarmDashboard(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := FarmDashboardHandler{repository}
	app.HandleFunc("GET /farm-area/dashboard/farm-info", handler.FarmInfo)
	app.HandleFunc("GET /farm-area/dashboard/pasture-info", handler.PastureInfo)
}
