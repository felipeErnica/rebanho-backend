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

func (r *PastureRepository) SearchPasture(filter *PastureFilter, userId string) (*[]Pasture, error) {
	query := `
        SELECT 
			p.id, 
			p.name,
			p.farm_id,
			p.bull_id,
			f.name AS farm_name
        FROM pastures p 
		JOIN farms f ON f.id = p.farm_id
    `
	filterExpression, _, err := util.GetFilterExpressions(filter, "p", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := util.GetWhereExpression("p.user_id = $1 AND p.deleted_at IS NULL", filterExpression)
	query += whereExpression + " ORDER BY p.name"

	args := []any{userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetList[Pasture](r.DB, query, args...)
}

func (r *PastureRepository) FindAnimalsByPasture(
	pastureId string,
	userId string,
	sort string,
	order string,
) (*[]PastureAnimal, error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"ring_number": {Field: "coalesce(nullif(regexp_replace(animals.ring_number, '[^0-9]', '', 'g'), '')::int, 0)", Order: "asc"},
		"name":        {Field: "coalesce(animals.name, '')", Order: "asc"},
		"birth_date":  {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"death_date":  {Field: "coalesce(animals.death_date, '-infinity')", Order: "asc"},
		"id":          {Field: "animals.id", Order: "asc"},
		"created_at":  {Field: "animals.created_at", Order: "asc"},
	}

	expression, err := util.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query := `
        SELECT
            animals.id, 
            animals.name, 
            animals.ring_number, 
            animals.sex, 
            animals.father_id, 
            animals.mother_id, 
            animals.birth_date, 
            animals.death_date, 
            animals.animal_type,
            CONCAT_WS(' - ', father.ring_number, father.name) AS father_name, 
            CONCAT_WS(' - ', mother.ring_number, mother.name) AS mother_name
        FROM animals
            LEFT JOIN animals AS father ON father.id = animals.father_id
            LEFT JOIN animals AS mother ON mother.id = animals.mother_id
        WHERE animals.pasture_id = $1 AND animals.user_id = $2 AND animals.deleted_at IS NULL
    `
	orderExpression := " ORDER BY " + expression
	query = query + orderExpression

	return util.GetList[PastureAnimal](r.DB, query, pastureId, userId)
}
