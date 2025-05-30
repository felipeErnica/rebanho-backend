package reproduction

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction/birth"
	embryoTransfer "github.com/felipeErnica/rebanho-backend/domains/reproduction/embryo-transfer"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction/insemination"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction/loss"
	naturalMating "github.com/felipeErnica/rebanho-backend/domains/reproduction/natural-mating"
	pregnancyTest "github.com/felipeErnica/rebanho-backend/domains/reproduction/pregnancy-test"
)

func InitReproduction(app *app.App) {
    insemination.InitInsemination(app)
    embryoTransfer.InitTransfer(app)
    birth.InitBirth(app)
    naturalMating.InitMating(app)
    loss.InitLoss(app)
	pregnancyTest.InitPregnacyTest(app)
}
