package insemination

import (
	"github.com/felipeErnica/rebanho-backend/app"
	inseminationEntries "github.com/felipeErnica/rebanho-backend/domains/reproduction/insemination/insemination-entries"
	inseminationGroup "github.com/felipeErnica/rebanho-backend/domains/reproduction/insemination/insemination-group"
)

func InitInsemination(app *app.App) {
    inseminationEntries.InitEntries(app)
    inseminationGroup.InitGroup(app)
}
