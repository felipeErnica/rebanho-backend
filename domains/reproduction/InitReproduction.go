package reproduction

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction/birth"
	embryoTransfer "github.com/felipeErnica/rebanho-backend/domains/reproduction/embryo-transfer"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction/insemination"
	naturalMating "github.com/felipeErnica/rebanho-backend/domains/reproduction/natural-mating"
	pregnancyTests "github.com/felipeErnica/rebanho-backend/domains/reproduction/pregnancy-tests"
)

func InitReproduction(app *app.App) {
    insemination.InitInsemination(app)
    embryoTransfer.InitTransfer(app)
    birth.InitBirth(app)
    naturalMating.InitMating(app)
    pregnancyTests.InitTests(app)
}
