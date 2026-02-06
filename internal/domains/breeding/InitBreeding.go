package breeding

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitBreeding(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewBreedingService(repository)
	handler := BreedingHandler{service}

	app.HandleFunc("GET /breeding/stats/pregnancy-rate", handler.GetPregnancyRateStats)
	app.HandleFunc("GET /breeding/stats/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /breeding/stats/insemination-hist", handler.GetInseminationHist)
	app.HandleFunc("GET /breeding/stats/future-births", handler.GetFutureBirths)
	app.HandleFunc("GET /breeding/stats/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /breeding/stats/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /breeding/stats/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /breeding/stats/last-entries", handler.GetLastEntries)

	app.HandleFunc("GET /breeding/page", handler.FindEntriesPage)
	app.HandleFunc("GET /breeding/page/foot", handler.GetEntriesFoot)

	app.HandleFunc("DELETE /breeding/{id}", handler.Delete)
	app.HandleFunc("PUT /breeding", handler.Update)
	app.HandleFunc("POST /breeding", handler.AddBreeding)

	app.HandleFunc("GET /breeding/groups/page", handler.FindGroups)
	app.HandleFunc("GET /breeding/groups/{breedingDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /breeding/groups/{breedingDate}/foot", handler.GetEntriesByGroupFoot)

	app.HandleFunc("PUT /breeding/groups/{breedingDate}", handler.UpdateGroup)
	app.HandleFunc("DELETE /breeding/groups/{breedingDate}", handler.DeleteGroup)

	log.LogDomainsInit("Cobertura")
}
