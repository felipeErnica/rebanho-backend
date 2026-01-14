package farm

import "github.com/felipeErnica/rebanho-backend/app"

func InitFarm(app *app.App) {
    repository := NewRepository(app.DBconn)
    handler := FarmHandler{repository}

	app.HandleFunc("GET /farm-area/farms/{id}/animals/total", handler.FindFarmAnimalsTotal)
    app.HandleFunc("GET /farm-area/farms/{id}/animals", handler.FindFarmAnimals)
}
