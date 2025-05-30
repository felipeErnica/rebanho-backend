package lactationArea

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/lactation-area/lactation"
)

func InitLactationArea(app *app.App) {
	lactation.InitLactation(app)
}
