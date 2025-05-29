package embryoTransfer

import "github.com/jmoiron/sqlx"

type TransferRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *TransferRepository {
    selectQuery := ``
    return &TransferRepository{selectQuery, "embryo_transfers", db}
}
