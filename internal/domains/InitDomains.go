package domains

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/domains/animals"
	"github.com/felipeErnica/rebanho-backend/internal/domains/auth"
	"github.com/felipeErnica/rebanho-backend/internal/domains/birth"
	"github.com/felipeErnica/rebanho-backend/internal/domains/breeding"
	"github.com/felipeErnica/rebanho-backend/internal/domains/butcher"
	"github.com/felipeErnica/rebanho-backend/internal/domains/cors"
	farmArea "github.com/felipeErnica/rebanho-backend/internal/domains/farm-area"
	"github.com/felipeErnica/rebanho-backend/internal/domains/lactation"
	"github.com/felipeErnica/rebanho-backend/internal/domains/milk"
	"github.com/felipeErnica/rebanho-backend/internal/domains/tests"
	"github.com/felipeErnica/rebanho-backend/internal/domains/reproduction"
	"github.com/felipeErnica/rebanho-backend/internal/domains/slaughter"
	"github.com/felipeErnica/rebanho-backend/internal/domains/weight"
)

func InitDomains(app *app.App) {
	cors.InitCorsOptions(app)
	animals.InitAnimal(app)
	auth.InitAuth(app)
	farmArea.InitFarmArea(app)
	lactation.InitLactationArea(app)
	milk.InitMilk(app)
	birth.InitBirth(app)
	tests.InitTests(app)
	reproduction.InitReproduction(app)
	weight.InitWeight(app)
	slaughter.InitSlaughter(app)
	butcher.InitButcher(app)
	breeding.InitBreeding(app)
}
