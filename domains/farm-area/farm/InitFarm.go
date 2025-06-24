package farm

import "github.com/felipeErnica/rebanho-backend/app"

func InitFarm(app *app.App) {
    repository := NewRepository(app.DBconn)
    handler := FarmHandler{repository}
    app.HandleFunc("GET /farm-area/farms/search", handler.SearchFarm)
    app.HandleFunc("POST /farm-area/farms/add", handler.Add)
}
