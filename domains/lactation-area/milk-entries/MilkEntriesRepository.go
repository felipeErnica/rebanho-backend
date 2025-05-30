package milkEntries

import "github.com/jmoiron/sqlx"

type MilkEntriesRepository struct {
	SelectQuery string
	TableName string
	Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *MilkEntriesRepository {
	selectQuery := ``
	return &MilkEntriesRepository{selectQuery, "milk_entries", db}
}
