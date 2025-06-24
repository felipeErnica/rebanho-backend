package animalTable

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type AnimalRepository struct {
	SelectQuery string
	NullFields  []string
	TableName   string
	DB          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *AnimalRepository {
	selectQuery := `
        SELECT animals.*, 
            father.name as father_name, father.number as father_number,
            mother.number as mother_number, mother.name as mother_name,
            pastures.name as pasture_name
        FROM animals
            LEFT JOIN animals as father ON father.id = animals.father_id
            LEFT JOIN animals as mother ON mother.id = animals.mother_id
            LEFT JOIN pastures ON pastures.id = animals.pasture_id
    `
	nullFields := []string{
		"chip_id",
		"name",
		"number",
		"color",
		"pasture_id",
		"father_id",
		"mother_id",
		"weaning_date",
		"birth_date",
		"death_date",
		"observation",
	}
	return &AnimalRepository{selectQuery, nullFields, "animals", db}
}

func (r *AnimalRepository) FindPage(props repositoriesUtil.PageProps) (page *entity.Page[Animal], err error) {
	buildProps := repositoriesUtil.PageBuilderProps{
		QueryBody:  r.SelectQuery,
		NullFields: r.NullFields,
		TableName:  r.TableName,
		DbConn:     r.DB,
		PageProps:  props,
	}
	return repositoriesUtil.BuildPage[Animal](buildProps)
}

func (r *AnimalRepository) FindById(id string) (*Animal, error) {
	query := r.SelectQuery + "WHERE animals.id = $1"
	return repositoriesUtil.GetOne[Animal](r.DB, query, id)
}

func (r *AnimalRepository) FindByFatherId(fatherId string) (*[]Animal, error) {
	query := r.SelectQuery + "WHERE animals.father_id = $1"
	return repositoriesUtil.GetList[Animal](r.DB, query, fatherId)
}

func (r *AnimalRepository) FindByMotherId(motherId string) (*[]Animal, error) {
	query := r.SelectQuery + "WHERE animals.mother_id = $1"
	return repositoriesUtil.GetList[Animal](r.DB, query, motherId)
}

func (r *AnimalRepository) FindByName(name string, userId string) (*[]Animal, error) {
	query := r.SelectQuery + "WHERE animals.name = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[Animal](r.DB, query, name, userId)
}

func (r *AnimalRepository) FindByNumber(number string, userId string) (*[]Animal, error) {
	query := r.SelectQuery + "WHERE animals.number = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[Animal](r.DB, query, number, userId)
}

func (r *AnimalRepository) SearchFather(userId string, input string) (*[]entity.SearchEntity, error) {
    query := `
        select id, concat_ws(' - ', number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'M' 
            and type <> 'OFFSPRING'
            and (name is not null)
            and concat_ws(' - ', number, name) ilike $2
            and deleted_at is null
        order by label
        limit 20
    `
    return repositoriesUtil.GetList[entity.SearchEntity](r.DB, query, userId, input)
}

func (r *AnimalRepository) SearchMother(userId string, input string) (*[]entity.SearchEntity, error) {
    query := `
        select id, concat_ws(' - ', number, name) as label 
            from animals 
        where user_id = $1 
            and sex = 'F' 
            and type <> 'OFFSPRING'
            and (name is not null)
            and concat_ws(' - ', number, name) ilike $2
            and deleted_at is null
        order by label
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
