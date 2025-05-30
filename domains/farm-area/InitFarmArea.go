package farmArea

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/farm-area/farm"
	"github.com/felipeErnica/rebanho-backend/domains/farm-area/pasture"
	pastureEntries "github.com/felipeErnica/rebanho-backend/domains/farm-area/pasture-entries"
)

func InitFarmArea(app *app.App) {
    farm.InitFarm(app)
    pasture.InitPasture(app)
    pastureEntries.InitPastureEntries(app)
}
