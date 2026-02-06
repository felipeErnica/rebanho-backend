package reproduction

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	embryoTransfer "github.com/felipeErnica/rebanho-backend/internal/domains/reproduction/embryo-transfer"
	"github.com/felipeErnica/rebanho-backend/internal/domains/reproduction/insemination"
)

func InitReproduction(app *app.App) {
	insemination.InitInsemination(app)
	embryoTransfer.InitTransfer(app)
}
