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
    animals.InitAnimal(app)
    auth.InitAuth(app)
    cors.InitCorsOptions(app)
    farmArea.InitFarmArea(app)
	lactationArea.InitLactationArea(app)
    reproduction.InitReproduction(app)
    slaughterArea.InitSlaughterArea(app)
}
