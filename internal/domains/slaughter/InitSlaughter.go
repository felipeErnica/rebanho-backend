package slaughter

import (
	"github.com/felipeErnica/rebanho-backend/internal/app"
	"github.com/felipeErnica/rebanho-backend/internal/log"
)

func InitSlaughter(app *app.App) {
	repository := NewRepository(app.DBconn)
	service := NewService(repository)
	handler := SlaughterHandler{service}

	app.HandleFunc("GET /slaughter/stats/last-weight", handler.GetLastAverageWeight)
	app.HandleFunc("GET /slaughter/stats/last-dead-weight", handler.GetLastDeadWeight)
	app.HandleFunc("GET /slaughter/stats/last-performance", handler.GetLastPerformance)
	app.HandleFunc("GET /slaughter/stats/weight-hist", handler.GetWeightHist)
	app.HandleFunc("GET /slaughter/stats/rate-hist", handler.GetRateHist)
	app.HandleFunc("GET /slaughter/stats/best-fathers", handler.GetBestFathers)
	app.HandleFunc("GET /slaughter/stats/best-mothers", handler.GetBestMothers)
	app.HandleFunc("GET /slaughter/stats/best-slaughterhouses", handler.GetBestButchers)
	app.HandleFunc("GET /slaughter/stats/last-entries", handler.GetLastEntries)

	app.HandleFunc("GET /slaughter/page", handler.FindPage)
	app.HandleFunc("GET /slaughter/page/foot", handler.GetPageFoot)
	app.HandleFunc("GET /slaughter/groups/{groupId}/entries", handler.FindEntries)
	app.HandleFunc("GET /slaughter/groups/{groupId}/entries/foot", handler.GetEntriesFoot)

	app.HandleFunc("POST /slaughter", handler.Add)
	app.HandleFunc("PUT /slaughter", handler.Update)
	app.HandleFunc("PUT /slaughter/batch", handler.UpdateBatch)
	app.HandleFunc("DELETE /slaughter/{id}", handler.Delete)
	app.HandleFunc("DELETE /slaughter/batch", handler.DeleteBatch)

	app.HandleFunc("GET /slaughter/butchers/{butcherId}/page", handler.FindButcherPage)
	app.HandleFunc("GET /slaughter/butchers/{butcherId}/page/foot", handler.GetButcherPageFoot)

	log.LogDomainsInit("Entradas de Abate")
}

