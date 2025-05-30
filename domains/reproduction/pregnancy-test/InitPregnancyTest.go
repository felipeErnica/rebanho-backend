package pregnancyTest

import (
	"github.com/felipeErnica/rebanho-backend/app"
	testEntries "github.com/felipeErnica/rebanho-backend/domains/reproduction/pregnancy-test/test-entries"
	testGroup "github.com/felipeErnica/rebanho-backend/domains/reproduction/pregnancy-test/test-group"
)

func InitPregnacyTest(app *app.App) {
	testEntries.InitEntries(app)
	testGroup.InitGroup(app)
}
