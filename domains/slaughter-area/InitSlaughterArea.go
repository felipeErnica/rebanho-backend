package slaughterArea

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/domains/slaughter-area/slaughter"
	"github.com/felipeErnica/rebanho-backend/domains/slaughter-area/slaughterhouses"
	"github.com/felipeErnica/rebanho-backend/domains/slaughter-area/weight"
)

func InitSlaughterArea(app *app.App) {
    slaughter.InitSlaughter(app)
    slaughterhouses.InitSlaughterhouse(app)
    weight.InitWeight(app)
}
