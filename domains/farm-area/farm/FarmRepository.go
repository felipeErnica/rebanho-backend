package farm

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type FarmRepository struct {
	Query     string
	TableName string
	DB        *sqlx.DB
}

func (r *FarmRepository) FindAnimalsByFarm(
    farmId string, 
    userId string, 
    filter FarmAnimalFilter,
    sort string, 
    order string,
) (*[]FarmAnimal, error) {

    sortMap := map[string]string{
        "ring_number": "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)",
        "name": "coalesce(animals.name, '')",
        "birth_date": "coalesce(animals.birth_date, '-infinity')",
        "death_date": "coalesce(animals.death_date, '-infinity')",
    }

    expression, err := repositoriesUtil.GetSortExpressionFromMap(sortMap, sort, order)
    if err != nil {
        return nil, err
    }

    query := `
        select
            animals.id, 
            animals.name, 
            animals.ring_number, 
            animals.sex, 
            animals.father_id, 
            animals.mother_id, 
            animals.birth_date, 
            animals.death_date, 
            animals.pasture_id, 
            animals.animal_type,
            concat_ws(' - ', father.ring_number, father.name) as father_name, 
            concat_ws(' - ', mother.ring_number, mother.name) as mother_name, 
            pastures.name as pasture_name
        from animals
            left join pastures on pastures.id = animals.pasture_id
            left join animals as father on father.id = animals.father_id
            left join animals as mother on mother.id = animals.mother_id
    `
    whereExpression := " where pastures.farm_id = $1 and animals.user_id = $2 and animals.deleted_at is null "
    filterExpressions, _, err := repositoriesUtil.BuildFilterExpressions(filter, "animals", 3)
    if err != nil {
        return nil, err
    }

    whereExpression = whereExpression + filterExpressions
    orderExpression := " order by " + expression

    query = query + whereExpression + orderExpression
    args := []any{ farmId, userId }
    filterArgs := repositoriesUtil.GetFilterArgs(filter)
    args = append(args, filterArgs...)

    return repositoriesUtil.GetList[FarmAnimal](r.DB, query, args...)
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

func (r *FarmRepository) Add(farm *Farm) (*Farm, error) {
	return repositoriesUtil.Add(r.DB, r.TableName, farm)
}
