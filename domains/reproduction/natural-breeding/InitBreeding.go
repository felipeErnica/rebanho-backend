package naturalBreeding

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitBreeding(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := BreedingHandler{repository}

	app.HandleFunc("GET /reproduction/breeding/dashboard/pregnancy-rate", handler.GetPregnancyRateStats)
	app.HandleFunc("GET /reproduction/breeding/dashboard/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /reproduction/breeding/dashboard/insemination-hist", handler.GetInseminationHist)
	app.HandleFunc("GET /reproduction/breeding/dashboard/future-births", handler.GetFutureBirths)
	app.HandleFunc("GET /reproduction/breeding/dashboard/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /reproduction/breeding/dashboard/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /reproduction/breeding/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /reproduction/breeding/dashboard/last-entries", handler.GetLastEntries)

	app.HandleFunc("POST /reproduction/breeding/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /reproduction/breeding/entries/page/foot", handler.GetEntriesFoot)
	app.HandleFunc("DELETE /reproduction/breeding/entries/{id}/delete", handler.Delete)
	app.HandleFunc("DELETE /reproduction/breeding/entries/{id}/delete-no-validation", handler.DeleteNoValidation)
	app.HandleFunc("DELETE /reproduction/breeding/entries/{id}/delete-change-father", handler.DeleteAndChangeFather)
	app.HandleFunc("PUT /reproduction/breeding/entries/update", handler.Update)
	app.HandleFunc("PUT /reproduction/breeding/entries/update-no-validation", handler.UpdateNoValidation)
	app.HandleFunc("PUT /reproduction/breeding/entries/add", handler.AddBreeding)
	app.HandleFunc("PUT /reproduction/breeding/entries/replace", handler.ReplaceBreeding)

	app.HandleFunc("GET /reproduction/breeding/groups/page", handler.FindGroups)
	app.HandleFunc("GET /reproduction/breeding/groups/{breedingDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /reproduction/breeding/groups/{breedingDate}/entries/foot", handler.GetEntriesByGroupFoot)
	app.HandleFunc("PUT /reproduction/breeding/groups/{breedingDate}/update", handler.UpdateBatch)
	app.HandleFunc("DELETE /reproduction/breeding/groups/{breedingDate}/delete", handler.DeleteBatch)

	app.HandleFunc("GET /reproduction/breeding/bulls/search-non-breeding", handler.SearchNonBreedingBulls)
	app.HandleFunc("GET /reproduction/breeding/bulls/search", handler.SearchBreedingBulls)
	app.HandleFunc("GET /reproduction/breeding/bulls/add/{id}", handler.AddBreedingBull)
	util.LogDomainsInit("Cobertura")
}
