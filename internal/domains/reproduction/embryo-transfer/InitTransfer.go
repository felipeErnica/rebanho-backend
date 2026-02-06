package embryoTransfer

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitTransfer(app *app.App) {
	repository := NewRepository(app.DBconn)
	handler := TransferHandler{repository}

	app.HandleFunc("GET /transfer/dashboard/pregnancy-rate", handler.GetPregnancyRateStats)
	app.HandleFunc("GET /transfer/dashboard/birth-rate", handler.GetBirthRateStats)
	app.HandleFunc("GET /transfer/dashboard/insemination-hist", handler.GetTransferHist)
	app.HandleFunc("GET /transfer/dashboard/future-births", handler.GetFutureBirths)
	app.HandleFunc("GET /transfer/dashboard/best-bull", handler.GetBestBull)
	app.HandleFunc("GET /transfer/dashboard/best-receivers", handler.GetBestReceivers)
	app.HandleFunc("GET /transfer/dashboard/best-donors", handler.GetBestDonors)
	app.HandleFunc("GET /transfer/dashboard/animals-number", handler.GetAnimalsNumber)
	app.HandleFunc("GET /transfer/dashboard/last-groups", handler.GetLastGroups)
	app.HandleFunc("GET /transfer/dashboard/last-entries", handler.GetLastEntries)

	app.HandleFunc("GET /transfer/page", handler.FindEntriesPage)
	app.HandleFunc("GET /transfer/page/foot", handler.GetEntriesFoot)

	app.HandleFunc("POST /transfer", handler.AddTransfer)
	app.HandleFunc("PUT /transfer", handler.Update)
	app.HandleFunc("DELETE /transfer/{id}", handler.Delete)

	app.HandleFunc("GET /transfer/groups/page", handler.FindGroups)
	app.HandleFunc("GET /transfer/groups/{transferDate}/entries", handler.FindEntriesByGroup)
	app.HandleFunc("GET /transfer/groups/{transferDate}/foot", handler.GetEntriesByGroupFoot)

	app.HandleFunc("PUT /transfer/groups", handler.UpdateGroup)
	app.HandleFunc("DELETE /transfer/groups/{transferDate}", handler.DeleteGroup)

	log.LogDomainsInit("Transferência Embrionária")
}
