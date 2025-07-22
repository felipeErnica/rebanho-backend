package animalTable

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type AnimalRepository struct {
	SelectQuery     string
	TableName       string
	SortExpressions []repositoriesUtil.SortExpression
	DB              *sqlx.DB
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
	sortExpressions := []repositoriesUtil.SortExpression{
		*repositoriesUtil.NewSort("name", "coalesce(animals.name, '')"),
		*repositoriesUtil.NewSort("isr", "coalesce(animals.isr, 0)"),
		*repositoriesUtil.NewSort("average_birth_interval", "coalesce(animals.average_birth_interval, 0)"),
		*repositoriesUtil.NewSort("average_prod_interval", "coalesce(animals.average_prod_interval, 0)"),
		*repositoriesUtil.NewSort("average_prod", "coalesce(animals.average_prod, 0)"),
		*repositoriesUtil.NewSort("average_peak", "coalesce(animals.average_peak, 0)"),
		*repositoriesUtil.NewSort("death_date", "coalesce(animals.death_date, '-infinity')"),
		*repositoriesUtil.NewSort("weaning_date", "coalesce(animals.weaning_date, '-infinity')"),
		*repositoriesUtil.NewSort("birth_date", "coalesce(animals.birth_date, '-infinity')"),
		*repositoriesUtil.NewSort("animal_order", "coalesce(regexp_replace(animals.ring_number, '[^0-9]', '', 'g')::int, 0)"),
	}
	return &AnimalRepository{selectQuery, "animals", sortExpressions, db}
}

func (r *AnimalRepository) FindPage(props repositoriesUtil.PageProps) (page *entity.Page[Animal], err error) {
	countQuery := `
        select count (animals.id) as total
        from animals
            left join animals as father ON father.id = animals.father_id
            left join animals as mother ON mother.id = animals.mother_id
            left join pastures ON pastures.id = animals.pasture_id
            left join farms ON farms.id = pastures.farm_id
    `
	buildProps := repositoriesUtil.PageBuilderProps{
		CountQuery:      countQuery,
		QueryBody:       r.SelectQuery,
		TableName:       r.TableName,
		DbConn:          r.DB,
		PageProps:       props,
		SortExpressions: r.SortExpressions,
	}
	return repositoriesUtil.BuildPage[Animal](buildProps)
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

func (r *AnimalRepository) SearchFather(userId string, input string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'M' 
            and animal_type <> 'OFFSPRING'
            and (name is not null)
            and concat_ws(' - ', ring_number, name) ilike $2
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), label
        limit 20
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId, input)
}

func (r *AnimalRepository) SearchAnimals(userId string, input string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name, to_char(birth_date, 'DD/MM/YYYY')) as label 
            from animals 
        where user_id = $1 
            and animal_type <> 'OUTSIDE_ANIMAL'
            and concat_ws(' - ', ring_number, name) ilike $2
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), birth_date
        limit 20
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId, input)
}

func (r *AnimalRepository) SearchMother(userId string, input string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'F' 
            and animal_type <> 'OFFSPRING'
            and (name is not null)
            and concat_ws(' - ', ring_number, name) ilike $2
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), label
        limit 20
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId, input)
}

func (r *AnimalRepository) SearchBull(userId string, input string) (*[]entity.SearchEntity, error) {
	query := `
        select id, concat_ws(' - ', ring_number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'M' 
            and animal_type = 'REPRODUCTION_ANIMAL'
            and (name is not null)
            and concat_ws(' - ', ring_number, name) ilike $2
            and deleted_at is null
        order by coalesce(regexp_replace(ring_number, '[^0-9]' ,'', 'g')::int, 0), label
        limit 20
    `
	return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId, input)
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
