package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/app"
	"github.com/felipeErnica/rebanho-backend/util"
)

func InitTransfer(app *app.App) {
    repository := NewRepository(app.DBconn)
    handler := TransferHandler{repository}

	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/pregnancy-rate", handler.GetPregnancyRateStats)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/insemination-hist", handler.GetTransferHist)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/future-births", handler.GetFutureBirths)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/best-receivers", handler.GetBestReceivers)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/best-donors", handler.GetBestDonors)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /reproduction/embryo-transfer/dashboard/last-entries", handler.GetLastEntries)

	app.HandleFunc("POST /reproduction/embryo-transfer/entries/page", handler.FindEntriesPage)
	app.HandleFunc("POST /reproduction/embryo-transfer/entries/page/foot", handler.GetEntriesFoot)
	app.HandleFunc("PUT /reproduction/embryo-transfer/entries/add", handler.AddTransfer)
	app.HandleFunc("PUT /reproduction/embryo-transfer/entries/replace", handler.Replace)
	app.HandleFunc("PUT /reproduction/embryo-transfer/entries/update", handler.Update)
	app.HandleFunc("DELETE /reproduction/embryo-transfer/entries/{id}/delete", handler.Delete)

	app.HandleFunc("GET /reproduction/embryo-transfer/groups/page", handler.FindGroups)
	app.HandleFunc("GET /reproduction/embryo-transfer/groups/{transferDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /reproduction/embryo-transfer/groups/{transferDate}/entries/foot", handler.GetEntriesByGroupFoot)
	app.HandleFunc("PUT /reproduction/embryo-transfer/groups/{transferDate}/update", handler.UpdateGroup)
	app.HandleFunc("DELETE /reproduction/embryo-transfer/groups/{transferDate}/delete", handler.DeleteGroup)

	app.HandleFunc("GET /reproduction/embryo-transfer/bulls/search", handler.SearchTransferBulls)
	app.HandleFunc("GET /reproduction/embryo-transfer/bulls/search-non-transfer", handler.SearchNonTransferBulls)
	app.HandleFunc("PUT /reproduction/embryo-transfer/bulls/{id}/add", handler.UpdateAsTransferBulls)

	app.HandleFunc("GET /reproduction/embryo-transfer/donors/search", handler.SearchEmbryoDonors)
	app.HandleFunc("GET /reproduction/embryo-transfer/donors/search-non-embryo", handler.SearchNonEmbryoDonors)
	app.HandleFunc("PUT /reproduction/embryo-transfer/donors/{id}/add", handler.UpdateAsEmbryoDonors)
	util.LogDomainsInit("Transferência Embrionária")
}
