package embryoTransfer

import "github.com/felipeErnica/rebanho-backend/app"

func InitTransfer(app *app.App) {
    repository := NewRepository(app.DBconn)
    handler := TransferHandler{repository}
    app.HandleFunc("GET /reproduction/embryo-transfer/", handler.FindAll)
}
