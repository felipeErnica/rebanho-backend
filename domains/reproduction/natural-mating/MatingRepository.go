package naturalMating

import "github.com/jmoiron/sqlx"

type MatingRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *MatingRepository {
    selectQuery := ``
    return &MatingRepository{selectQuery, "natural_matings", db}
}
