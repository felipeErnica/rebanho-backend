package pasture

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
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
        select 
			id, 
			name,
			farm_id,
			bull_id,
			f.name as farm_name,
			concat_ws(' - ', b.ring_number, b.name) as bull_name
        from pastures p
			join farms f on f.id = p.farm_id
			left join animals b on b.id = p.bull_id 
    `
	filterExpression, _, err := repositoriesUtil.GetFilterExpressions(filter, "p", 2)
	if err != nil {
		return nil, err
	}

	whereExpression := repositoriesUtil.GetWhereExpression("p.user_id = $1 and p.deleted_at is null", filterExpression)
	query += whereExpression + " order by p.name"

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)

	return repositoriesUtil.GetList[Pasture](r.DB, query, args...)
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

	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
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

	return repositoriesUtil.GetList[PastureAnimal](r.DB, query, pastureId, userId)
}
