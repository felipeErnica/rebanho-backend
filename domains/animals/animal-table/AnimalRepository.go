package animalTable

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type AnimalRepository struct {
	SelectQuery string
	TableName   string
	DB          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *AnimalRepository {
	selectQuery := `
        select animals.*, 
            coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0) as animal_order,
            concat_ws(' - ', father.ring_number, father.name) as father_name, 
            concat_ws(' - ', mother.ring_number, mother.name) as mother_name,
            pastures.name as pasture_name,
            farms.id as farm_id, farms.name as farm_name
        from animals
            left join animals as father ON father.id = animals.father_id
            left join animals as mother ON mother.id = animals.mother_id
            left join pastures ON pastures.id = animals.pasture_id
            left join farms ON farms.id = pastures.farm_id
    `
	return &AnimalRepository{selectQuery, "animals", db}
}

func (r *AnimalRepository) FindPage(
	userId string,
	cursor string,
	sort string,
	order string,
	filter AnimalFilter,
) (page *entity.Page[Animal], err error) {

	sort = repositoriesUtil.AddCommonFields(sort)
	sortMap := map[string]repositoriesUtil.SortField{
		"name":                   {Field: "coalesce(animals.name, '')", Order: "asc"},
		"isr":                    {Field: "coalesce(animals.isr, 0)", Order: "asc"},
		"average_birth_interval": {Field: "coalesce(animals.average_birth_interval, 0)", Order: "asc"},
		"average_prod_interval":  {Field: "coalesce(animals.average_prod_interval, 0)", Order: "asc"},
		"average_prod":           {Field: "coalesce(animals.average_prod, 0)", Order: "asc"},
		"average_peak":           {Field: "coalesce(animals.average_peak, 0)", Order: "asc"},
		"death_date":             {Field: "coalesce(animals.death_date, '-infinity')", Order: "asc"},
		"weaning_date":           {Field: "coalesce(animals.weaning_date, '-infinity')", Order: "asc"},
		"birth_date":             {Field: "coalesce(animals.birth_date, '-infinity')", Order: "asc"},
		"animal_order":           {Field: "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)", Order: "asc"},
	}

	whereExpression := "where animals.deleted_at is null and animals.user_id = $1"
	sortExpression, err := repositoriesUtil.GetSortExpression(sortMap, sort, order)
	if err != nil {
		return nil, err
	}
	sortExpression = " order by " + sortExpression

	cursorArgs, err := repositoriesUtil.GetCursorArgs(cursor)
	if err != nil {
		return nil, err
	}

	filterExpression, nextParam, err := repositoriesUtil.GetFilterExpressions(filter, "animals", 2)
	if err != nil {
		return nil, err
	}

	cursorExpression, _, err := repositoriesUtil.GetCursorExpression(sortMap, sort, order, cursor, nextParam)
	if err != nil {
		return nil, err
	}

	if filterExpression != "" {
		whereExpression = whereExpression + " and " + filterExpression
	}

	if cursorExpression != "" {
		whereExpression = whereExpression + " and " + cursorExpression
	}

	args := []any{userId}
	filterArgs := repositoriesUtil.GetFilterArgs(filter)
	args = append(args, filterArgs...)
	args = append(args, cursorArgs...)

	query := r.SelectQuery + whereExpression + sortExpression
	return repositoriesUtil.GetPage[Animal](r.DB, query, sort, 200, args...)
}

func (r *AnimalRepository) FindById(id string) (*Animal, error) {
	query := r.SelectQuery + "where animals.id = $1"
	return repositoriesUtil.GetOne[Animal](r.DB, query, id)
}

func (r *AnimalRepository) FindByFatherId(fatherId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.father_id = $1"
	return repositoriesUtil.GetList[Animal](r.DB, query, fatherId)
}

func (r *AnimalRepository) FindByMotherId(motherId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.mother_id = $1"
	return repositoriesUtil.GetList[Animal](r.DB, query, motherId)
}

func (r *AnimalRepository) FindByName(name string, userId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.name = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[Animal](r.DB, query, name, userId)
}

func (r *AnimalRepository) FindByNumber(ringNumber string, userId string) (*[]Animal, error) {
	query := r.SelectQuery + "where animals.ring_number = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[Animal](r.DB, query, ringNumber, userId)
}

func (r *AnimalRepository) SearchFather(userId string) (*[]entity.SearchEntity, error) {
	queryInput := `
        select id, concat_ws(' - ', ring_number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'M' 
            and animal_type <> 'OFFSPRING'
            and (name is not null)
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), coalesce(name, '')
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, queryInput, userId)
}

func (r *AnimalRepository) SearchAnimals(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name, to_char(birth_date, 'DD/MM/YYYY')) as label 
            from animals 
        where user_id = $1 
            and animal_type <> 'OUTSIDE_ANIMAL'
            and deleted_at is null
        order by 
			coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), 
			coalesce(name, ''), 
			coalesce(birth_date, '-infinity')
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) SearchMother(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'F' 
            and animal_type <> 'OFFSPRING'
            and (name is not null)
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), label
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) SearchDairyAnimals(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
		from animals 
        where user_id = $1 
            and sex = 'F' 
            and animal_type = 'DAIRY_ANIMAL'
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), name
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) SearchBull(userId string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'M' 
            and animal_type = 'REPRODUCTION_ANIMAL'
            and (name is not null)
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), label
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId)
}

func (r *AnimalRepository) Add(create *AnimalSave) (*AnimalSave, error) {
	return repositoriesUtil.Add(r.DB, r.TableName, create)
}

func (r *AnimalRepository) Update(animal *AnimalSave) error {
	return repositoriesUtil.Update(r.DB, r.TableName, animal)
}

func (r *AnimalRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.DB, r.TableName, id)
}
