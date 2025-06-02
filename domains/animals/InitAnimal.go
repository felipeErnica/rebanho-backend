package animals

import (
	"github.com/felipeErnica/rebanho-backend/app"
	animalDashboard "github.com/felipeErnica/rebanho-backend/domains/animals/animal-dashboard"
	animalTable "github.com/felipeErnica/rebanho-backend/domains/animals/animal-table"
)

func InitAnimal(app *app.App) {
    animalDashboard.InitDashboard(app)
    animalTable.InitTable(app)
}
