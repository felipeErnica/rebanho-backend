package domains

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/animals"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction"
)

func InitDomains(app *app.App) {
    animals.InitAnimal(app)
    reproduction.InitReproduction(app)
}
