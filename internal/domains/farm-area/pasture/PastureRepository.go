package pasture

import (
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type PastureRepository struct {
	DB *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PastureRepository {
	return &PastureRepository{db}
}

func (r *PastureRepository) Search(filter *PastureFilter, userId string) (*[]PastureDB, error) {
	query := `
		WITH bull_cte AS (
			SELECT DISTINCT ON (pe.animal_id)
				pe.pasture_id,
				pe.animal_id AS bull_id,
				b.name AS bull_name,
				b.tag AS bull_tag
			FROM pasture_entries pe
			JOIN animals b ON b.id = pe.animal_id
			WHERE b.sex = 'M'
				AND b.name IS NOT NULL
				AND pe.exit_date IS NULL
				AND pe.user_id = $1
			ORDER BY pe.animal_id, pe.entry_date DESC
		),

        cte as (
			SELECT 
				p.id, 
				p.name,

				p.farm_id,
				f.name AS farm_name,

				b.bull_id,
				b.bull_name,
				b.bull_tag
			FROM pastures p 
			JOIN farms f ON f.id = p.farm_id
			LEFT JOIN bull_cte b ON b.pasture_id = p.id
			WHERE p.user_id = $1 AND p.deleted_at IS NULL
		)

		SELECT * FROM cte
    `
	filterExpression, _, err := util.GetFilterExpressions(filter, "cte", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression(filterExpression)
	query += whereExpression + " ORDER BY cte.name"

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetList[PastureDB](r.DB, query, args...)
}

func (r *PastureRepository) FindAnimalsById(
	pastureId string,
	userId string,
	sort string,
	order string,
) (*[]PastureAnimal, error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"tag":        {Field: "COALESCE(NULL_IF(REGEX_REPLACE(animals.tag, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"name":       {Field: "COALESCE(animals.name, '')", Order: "asc"},
		"birth_date": {Field: "COALESCE(animals.birth_date, '-infinity')", Order: "asc"},
		"death_date": {Field: "COALESCE(animals.death_date, '-infinity')", Order: "asc"},
		"id":         {Field: "animals.id", Order: "asc"},
		"created_at": {Field: "animals.created_at", Order: "asc"},
	}

	expression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query := `
		WITH entries_cte AS (
			SELECT DISTINCT ON (pe.animal_id)
				pe.animal_id
			FROM pasture_entries pe
			WHERE pe.pasture_id = $1
				AND pe.exit_date IS NULL
				AND pe.user_id = $2
				AND pe.deleted_at IS NULL
			ORDER BY pe.animal_id, pe.entry_date DESC
		)

        SELECT
            a.id, 
            a.name, 
            a.tag, 
            a.sex, 
            a.birth_date, 
            a.animal_type,
        FROM animals a
		JOIN entries_cte e ON e.animal_id = a.id
        WHERE a.death_date IS NULL
			AND a.user_id = $2
			AND a.deleted_at IS NULL
    `
	orderExpression := " ORDER BY " + expression
	query = query + orderExpression

	return util.GetList[PastureAnimal](r.DB, query, pastureId, userId)
}
