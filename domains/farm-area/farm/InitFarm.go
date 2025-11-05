package farm

import "github.com/felipeErnica/rebanho-backend/app"

func InitFarm(app *app.App) {
    repository := NewRepository(app.DBconn)
    handler := FarmHandler{repository}
        app.HandleFunc("POST /farm-area/farms/{id}/animals/total", handler.FindFarmAnimalsTotal)
    app.HandleFunc("POST /farm-area/farms/{id}/animals", handler.FindFarmAnimals)
    app.HandleFunc("GET /farm-area/farms/search", handler.SearchFarm)
}
