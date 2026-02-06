package farm

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/internal/entity"
	"github.com/felipeErnica/rebanho-backend/internal/util"
	"github.com/jmoiron/sqlx"
)

type FarmRepository struct {
	Query     string
	TableName string
	DB        *sqlx.DB
}

func (r *FarmRepository) FindFarmAnimals(
	farmId string,
	userId string,
	filter *FarmAnimalFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[FarmAnimal], error) {

	sort = util.AddCommonFields(sort)
	sortMap := map[string]util.SortField{
		"ring_order": {Field: "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"name":       {Field: "coalesce(animals.name, '')", Order: "asc"},
		"birth_date": {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"id":         {Field: "animals.id", Order: "asc"},
		"created_at": {Field: "animals.created_at", Order: "asc"},
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
            COALESCE(REGEXP_REPLACE(animals.ring_number, '[^0-9]', '', 'g')::int, 0) AS ring_order,
            animals.sex, 
            animals.father_id, 
            animals.mother_id, 
            animals.birth_date, 
            animals.pasture_id, 
            animals.animal_type,
            animals.created_at,
            CONCAT_WS(' - ', father.ring_number, father.name) AS father_name, 
            CONCAT_WS(' - ', mother.ring_number, mother.name) AS mother_name, 
            pastures.name AS pasture_name,
            pastures.farm_id AS farm_id
        FROM animals
            LEFT JOIN pastures ON pastures.id = animals.pasture_id
            LEFT JOIN animals AS father ON father.id = animals.father_id
            LEFT JOIN animals AS mother ON mother.id = animals.mother_id
    `
	whereExpression := " WHERE pastures.farm_id = $1 AND animals.user_id = $2 AND animals.deleted_at IS NULL "

	filterExpressions, nextParam, err := util.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpressions != "" {
		whereExpression = whereExpression + " AND " + filterExpressions
	}

	cursorArgs, err := util.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := util.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if cursorExpression != "" {
		whereExpression = whereExpression + " AND " + cursorExpression
	}

	orderExpression := " ORDER BY " + expression

	query = query + whereExpression + orderExpression
	args := []any{farmId, userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	return util.GetPage[FarmAnimal](r.DB, query, sort, 200, args...)
}

func (r *FarmRepository) FindFarmAnimalsTotal(
	farmId string,
	userId string,
	filter *FarmAnimalFilter,
) (*FarmAnimalTotal, error) {

	query := `
        SELECT COUNT(animals.id) AS total 
        FROM animals 
            LEFT JOIN pastures ON pastures.id = animals.pasture_id
        WHERE pastures.farm_id = $1 AND animals.user_id = $2 AND animals.deleted_at IS NULL 
    `
	filterExpressions, _, err := util.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpressions != "" {
		query += " AND " + filterExpressions
	}

	args := []any{farmId, userId}
	filterArgs := util.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return util.GetOne[FarmAnimalTotal](r.DB, query, args...)
}

func NewRepository(db *sqlx.DB) *FarmRepository {
	query := "SELECT farms.* FROM farms"
	return &FarmRepository{query, "farms", db}
}

func (r *FarmRepository) SearchFarm(userId string, input string) (*[]entity.SearchEntity, error) {
	query := `
        SELECT id, name AS label 
        FROM farms 
        WHERE user_id = $1 AND name ILIKE $2 AND deleted_at IS NULL
        ORDER BY label
        `
	return util.GetList[entity.SearchEntity](r.DB, query, userId, input)
}

func (r *FarmRepository) SearchFarmById(userId string, idList []string) (*[]entity.SearchEntity, error) {
	query := `
        SELECT id, name AS label 
        FROM farms 
        WHERE user_id = $1 AND deleted_at IS NULL
        ORDER BY name
        `
	if len(idList) != 0 {
		queryId := "SELECT id, name AS label FROM farms"
		idExpression, _ := util.GetSliceExpressions(idList, "id", 2)
		queryId += " WHERE " + idExpression
		query = fmt.Sprintf(`
            WITH farm_base AS (%s),
            farm_id AS (%s)
            SELECT * FROM farm_base
            UNION
            SELECT * FROM farm_id
        `, query, queryId)
		args := []any{userId}
		args = util.GetSliceArgs(idList, args)
		return util.GetList[entity.SearchEntity](r.DB, query, userId)
	}
	return util.GetList[entity.SearchEntity](r.DB, query, userId)
}
