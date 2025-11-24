package reproduction

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction/birth"
	embryoTransfer "github.com/felipeErnica/rebanho-backend/domains/reproduction/embryo-transfer"
	"github.com/felipeErnica/rebanho-backend/domains/reproduction/insemination"
	naturalBreeding "github.com/felipeErnica/rebanho-backend/domains/reproduction/natural-breeding"
	pregnancyTests "github.com/felipeErnica/rebanho-backend/domains/reproduction/pregnancy-tests"
)

func InitReproduction(app *app.App) {
    insemination.InitInsemination(app)
    embryoTransfer.InitTransfer(app)
    birth.InitBirth(app)
    naturalBreeding.InitBreeding(app)
    pregnancyTests.InitTests(app)
}
