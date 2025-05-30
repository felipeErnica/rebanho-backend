package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/app"
	slaughterEntry "github.com/felipeErnica/rebanho-backend/domains/slaughter-area/slaughter/slaughter-entry"
	slaughterGroup "github.com/felipeErnica/rebanho-backend/domains/slaughter-area/slaughter/slaughter-group"
)

func InitSlaughter(app *app.App) {
    slaughterEntry.InitSlaughterEntry(app)
    slaughterGroup.InitSlaughterGroup(app)
}
