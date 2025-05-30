package lactationArea

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/lactation-area/lactation"
	milkEntries "github.com/felipeErnica/rebanho-backend/domains/lactation-area/milk-entries"
)

func InitLactationArea(app *app.App) {
	lactation.InitLactation(app)
    milkEntries.InitMilkEntries(app)
}
