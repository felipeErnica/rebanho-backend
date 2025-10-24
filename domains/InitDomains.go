package domains

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/animals"
	"github.com/felipeErnica/rebanho-backend/domains/auth"
	"github.com/felipeErnica/rebanho-backend/domains/cors"
	farmArea "github.com/felipeErnica/rebanho-backend/domains/farm-area"
	"github.com/felipeErnica/rebanho-backend/domains/lactation"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction"
	"github.com/felipeErnica/rebanho-backend/domains/slaughter"
	"github.com/felipeErnica/rebanho-backend/domains/weight"
)

func InitDomains(app *app.App) {
    cors.InitCorsOptions(app)
    animals.InitAnimal(app)
    auth.InitAuth(app)
    farmArea.InitFarmArea(app)
	lactation.InitLactationArea(app)
    reproduction.InitReproduction(app)
	weight.InitWeight(app)
    slaughter.InitSlaughter(app)
}
