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
    go animals.InitAnimal(app)
    go auth.InitAuth(app)
    go cors.InitCorsOptions(app)
    go farmArea.InitFarmArea(app)
	go lactation.InitLactationArea(app)
    go reproduction.InitReproduction(app)
	go weight.InitWeight(app)
    go slaughter.InitSlaughter(app)
}
