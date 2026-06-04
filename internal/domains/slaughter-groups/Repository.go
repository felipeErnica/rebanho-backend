package slaughtergroups

import (
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) FindAll(userId string) (*[]GroupDB, error) {
	query := `
		WITH entries_stats AS (
			SELECT
				group_id,
				COUNT(*) AS animals_number,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0))) AS average_rate
				AVG(weight) AS average_weight,
				AVG(dead_weight) AS avearge_dead_weight
			FROM slaughter_entries
			GROUP BY group_id
		)

		SELECT
			g.id,
			g.entry_date,
			g.discount_rate,

			g.butcher_id,
			b.name AS butcher_name,
			b.discount_rate AS butcher_discount

			s.animals_number,
			s.average_weight,
			s.avearge_dead_weight,
			s.average_rate
		FROM slaughter_groups g
		JOIN butchers b ON b.id = g.butcher_id
		LEFT JOIN entries_stats s ON s.group_id = g.id 
		WHERE g.user_id = $1 AND g.deleted_at IS NULL
	`
	return util.GetList[GroupDB](r.DB, query, userId)
}

func (r *Repository) FindById(id string, userId string) (*GroupDB, error) {
	query := `
		WITH entries_stats AS (
			SELECT
				group_id,
				COUNT(*) AS animals_number,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0))) AS average_rate
				AVG(weight) AS average_weight,
				AVG(dead_weight) AS avearge_dead_weight
			FROM slaughter_entries
			GROUP BY group_id
		)

		SELECT
			g.id,
			g.entry_date,
			g.discount_rate,

			g.butcher_id,
			b.name AS butcher_name,
			b.discount_rate AS butcher_discount

			s.animals_number,
			s.average_weight,
			s.avearge_dead_weight,
			s.average_rate
		FROM slaughter_groups g
		JOIN butchers b ON b.id = g.butcher_id
		LEFT JOIN entries_stats s ON s.group_id = g.id 
		WHERE g.id = $1
			AND g.user_id = $2
			AND g.deleted_at IS NULL
	`
	return util.GetOne[GroupDB](r.DB, query, id, userId)
}

func (r *Repository) FindLast(userId string) (*[]GroupDB, error) {
	query := `
		WITH entries_stats AS (
			SELECT
				group_id,
				COUNT(*) AS animals_number,
				AVG((dead_weight / NULLIF(weight * (1 - discount_rate), 0))) AS average_rate
				AVG(weight) AS average_weight,
				AVG(dead_weight) AS avearge_dead_weight
			FROM slaughter_entries
			WHERE user_id = $1 AND deleted_at IS NULL
			GROUP BY group_id
		)

		SELECT
			g.id,
			g.entry_date,
			g.discount_rate,

			g.butcher_id,
			b.name AS butcher_name,
			b.discount_rate AS butcher_discount

			s.animals_number,
			s.average_weight,
			s.avearge_dead_weight,
			s.average_rate
		FROM slaughter_groups g
		JOIN butchers b ON b.id = g.butcher_id
		LEFT JOIN entries_stats s ON s.group_id = g.id 
		WHERE g.user_id = $1 AND g.deleted_at IS NULL
		ORDER BY g.entry_date DESC
		LIMIT 10
	`
	return util.GetList[GroupDB](r.DB, query, userId)
}
