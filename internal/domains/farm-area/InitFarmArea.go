package farmArea

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/domains/farm-area/dashboard"
	"github.com/felipeErnica/rebanho-backend/internal/domains/farm-area/farm"
	"github.com/felipeErnica/rebanho-backend/internal/domains/farm-area/pasture"
	pastureEntries "github.com/felipeErnica/rebanho-backend/internal/domains/farm-area/pasture-entries"
)

func InitFarmArea(app *app.App) {
	farm.InitFarm(app)
	dashboard.InitFarmDashboard(app)
	pasture.InitPasture(app)
	pastureEntries.InitPastureEntries(app)
}
