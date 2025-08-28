package farm

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
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
	filter FarmAnimalFilter,
	sort string,
	order string,
	cursor string,
) (*entity.Page[FarmAnimal], error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"ring_order": {Field: "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"name":       {Field: "coalesce(animals.name, '')", Order: "asc"},
		"birth_date": {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"id":         {Field: "animals.id", Order: "asc"},
		"created_at": {Field: "animals.created_at", Order: "asc"},
	}

	expression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}

	query := `
        select
            animals.id, 
            animals.name, 
            animals.ring_number, 
            coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0) as ring_order,
            animals.sex, 
            animals.father_id, 
            animals.mother_id, 
            animals.birth_date, 
            animals.pasture_id, 
            animals.animal_type,
            animals.created_at,
            concat_ws(' - ', father.ring_number, father.name) as father_name, 
            concat_ws(' - ', mother.ring_number, mother.name) as mother_name, 
            pastures.name as pasture_name,
            pastures.farm_id as farm_id
        from animals
            left join pastures on pastures.id = animals.pasture_id
            left join animals as father on father.id = animals.father_id
            left join animals as mother on mother.id = animals.mother_id
    `
	whereExpression := " where pastures.farm_id = $1 and animals.user_id = $2 and animals.deleted_at is null "

	filterExpressions, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpressions != "" {
		whereExpression = whereExpression + " and " + filterExpressions
	}

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if cursorExpression != "" {
		whereExpression = whereExpression + " and " + cursorExpression
	}

	orderExpression := " order by " + expression

	query = query + whereExpression + orderExpression
	args := []any{farmId, userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	return repositoriesUtil.GetPage[FarmAnimal](r.DB, query, sort, 200, args...)
}

func (r *FarmRepository) FindFarmAnimalsTotal(
	farmId string,
	userId string,
	filter FarmAnimalFilter,
) (*FarmAnimalTotal, error) {

	query := `
        select count(animals.id) as total 
        from animals 
            left join pastures on pastures.id = animals.pasture_id
        where pastures.farm_id = $1 and animals.user_id = $2 and animals.deleted_at is null 
    `
	filterExpressions, _, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 3)
	if err != nil {
		return nil, err
	}

	if filterExpressions != "" {
		query += " and " + filterExpressions
	}

	args := []any{farmId, userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetOne[FarmAnimalTotal](r.DB, query, args...)
}

func NewRepository(db *sqlx.DB) *FarmRepository {
	query := "SELECT farms.* FROM farms"
	return &FarmRepository{query, "farms", db}
}

func (r *FarmRepository) SearchFarm(userId string, input string) (*[]entity.SearchEntity, error) {
	query := `
        select id, name as label 
        from farms 
        where user_id = $1 and name ilike $2 and deleted_at is null
        order by label
        `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId, input)
}

func (r *FarmRepository) SearchFarmById(userId string, idList []string) (*[]entity.SearchEntity, error) {
	query := `
        select id, name as label 
        from farms 
        where user_id = $1 and deleted_at is null
        order by name
        `
	if len(idList) != 0 {
		queryId := "select id, name as label from farms"
		idExpression, _ := repositoriesUtil.GetSliceExpressions(idList, "id", 2)
		queryId += " where " + idExpression
		query = fmt.Sprintf(`
            with farm_base as (%s),
            farm_id as (%s)
            select * from farm_base
            union
            select * from farm_id
        `, query, queryId)
		args := []any{userId}
		args = repositoriesUtil.GetSliceArgs(idList, args)
		return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
	}
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *FarmRepository) Add(farm *Farm) (*Farm, error) {
	return repositoriesUtil.Add(r.DB, r.TableName, farm)
}
