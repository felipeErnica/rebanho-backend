package domains

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/animals"
	"github.com/felipeErnica/rebanho-backend/domains/auth"
	"github.com/felipeErnica/rebanho-backend/domains/cors"
	farmArea "github.com/felipeErnica/rebanho-backend/domains/farm-area"
	lactationArea "github.com/felipeErnica/rebanho-backend/domains/lactation-area"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction"
	slaughterArea "github.com/felipeErnica/rebanho-backend/domains/slaughter-area"
)

func InitDomains(app *app.App) {
    go animals.InitAnimal(app)
    go auth.InitAuth(app)
    go cors.InitCorsOptions(app)
    go farmArea.InitFarmArea(app)
	go lactationArea.InitLactationArea(app)
    go reproduction.InitReproduction(app)
    go slaughterArea.InitSlaughterArea(app)
}
