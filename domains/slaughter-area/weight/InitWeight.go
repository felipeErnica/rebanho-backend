package weight

import (
	"github.com/felipeErnica/rebanho-backend/app"
	weightEntries "github.com/felipeErnica/rebanho-backend/domains/slaughter-area/weight/weight-entries"
	weightGroup "github.com/felipeErnica/rebanho-backend/domains/slaughter-area/weight/weight-group"
)

func InitWeight(app *app.App) {
    weightEntries.InitWeightEntries(app)
    weightGroup.InitWeightGroup(app)
}
