package pasture

import (
	"fmt"

	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type PastureRepository struct {
	Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PastureRepository {
	return &PastureRepository{db}
}

func (r *PastureRepository) SearchPasture(userId string, farmsId []string) (*[]entity.SearchEntity, error) {
	args := []any{userId}

	arrayStatement := ""
	if len(farmsId) != 0 {
		args = repositoriesUtil.GetSliceArgs(farmsId, args)
		farmExpression, _ := repositoriesUtil.GetSliceExpressions(farmsId, "farm_id", 3)
		arrayStatement = " and " + farmExpression
	}

	query := fmt.Sprintf(`
        select id, name as label 
        from pastures 
        where 
            user_id = $1 
            %s
            and deleted_at is null
        order by label
    `, arrayStatement)

	return repositoriesUtil.GetList[entity.SearchEntity](r.Db, query, args...)
}

func (r *PastureRepository) SearchAllPastures(userId string) (*[]entity.SearchEntity, error) {

	query := `
        select 
			p.id, 
			format('%s (%s)', p.name, f.name) as label 
        from pastures p
			join farms f on f.id = p.farm_id
        where p.user_id = $1 and p.deleted_at is null
        order by label
    `

	return repositoriesUtil.GetList[entity.SearchEntity](r.Db, query, userId)
}

func (r *PastureRepository) FindAnimalsByPasture(
	pastureId string,
	userId string,
	sort string,
	order string,
) (*[]PastureAnimal, error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"ring_number": {Field: "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
		"name":        {Field: "coalesce(animals.name, '')", Order: "asc"},
		"birth_date":  {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"death_date":  {Field: "coalesce(animals.death_date, '-infinity')", Order: "asc"},
		"id":          {Field: "animals.id", Order: "asc"},
		"created_at":  {Field: "animals.created_at", Order: "asc"},
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
            animals.sex, 
            animals.father_id, 
            animals.mother_id, 
            animals.birth_date, 
            animals.death_date, 
            animals.animal_type,
            concat_ws(' - ', father.ring_number, father.name) as father_name, 
            concat_ws(' - ', mother.ring_number, mother.name) as mother_name
        from animals
            left join animals as father on father.id = animals.father_id
            left join animals as mother on mother.id = animals.mother_id
        where animals.pasture_id = $1 and animals.user_id = $2 and animals.deleted_at is null
    `
	orderExpression := " order by " + expression
	query = query + orderExpression

	return repositoriesUtil.GetList[PastureAnimal](r.Db, query, pastureId, userId)
}
