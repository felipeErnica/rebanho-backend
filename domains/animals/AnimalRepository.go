package animals

import (
	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/repositories"
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type AnimalRepository struct {
	SelectQuery string
	NullFields  []string
	TableName   string
	DB          *sqlx.DB
}

func NewAnimalRepository(db *sqlx.DB) *AnimalRepository {
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

func (r *AnimalRepository) FindPage(
	sort string,
	direction string,
	cursor string,
	filter *entity.AnimalFilter,
) (page *entity.Page[entity.Animal], err error) {
	props := repositoriesUtil.PageProps{
		QueryBody:  r.SelectQuery,
		Sort:       sort,
		Order:      direction,
		Cursor:     cursor,
		Filter:     filter,
		NullFields: r.NullFields,
		Limit:      repositories.PAGE_LIMIT,
		TableName:  r.TableName,
		DbConn:     r.DB,
		UserId:     repositories.GetUserId(),
	}
	return repositoriesUtil.BuildPage[entity.Animal](props)
}

func (r *AnimalRepository) FindById(id string) (*entity.Animal, error) {
	query := r.SelectQuery + "WHERE animals.id = $1"
	return repositoriesUtil.GetOne[entity.Animal](r.DB, query, id)
}

func (r *AnimalRepository) FindByFatherId(fatherId string) (*[]entity.Animal, error) {
	query := r.SelectQuery + "WHERE animals.father_id = $1"
	return repositoriesUtil.GetList[entity.Animal](r.DB, query, fatherId)
}

func (r *AnimalRepository) FindByMotherId(motherId string) (*[]entity.Animal, error) {
	query := r.SelectQuery + "WHERE animals.mother_id = $1"
	return repositoriesUtil.GetList[entity.Animal](r.DB, query, motherId)
}

func (r *AnimalRepository) FindByPastureId(
	sort string,
	direction string,
	cursor string,
	filter entity.AnimalFilter,
	pastureId string,
) (page *entity.Page[entity.Animal], err error) {
	//Modifica o filtro para obter uma página dos animais neste pasto
	filter.IsFiltered = true
	filter.Pastures = &[]string{pastureId}

	props := repositoriesUtil.PageProps{
		QueryBody:  r.SelectQuery,
		Sort:       sort,
		Order:      direction,
		Cursor:     cursor,
		Filter:     filter,
		NullFields: r.NullFields,
		Limit:      repositories.PAGE_LIMIT,
		TableName:  r.TableName,
		DbConn:     r.DB,
		UserId:     repositories.GetUserId(),
	}
	return repositoriesUtil.BuildPage[entity.Animal](props)
}

func (r *AnimalRepository) FindByName(name string) (*[]entity.Animal, error) {
	query := r.SelectQuery + "WHERE animals.name = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[entity.Animal](r.DB, query, name, repositories.GetUserId())
}

func (r *AnimalRepository) FindByIdentificationNumber(number string) (*[]entity.Animal, error) {
	query := r.SelectQuery + "WHERE animals.number = $1 AND animals.user_id = $2"
	return repositoriesUtil.GetList[entity.Animal](r.DB, query, number, repositories.GetUserId())
}

func (r *AnimalRepository) Add(create *entity.AnimalSave) (*entity.AnimalSave, error) {
	err := repositoriesUtil.Add(r.DB, r.TableName, create)
	return create, err
}

func (r *AnimalRepository) Save(animal *entity.AnimalSave) error {
	return repositoriesUtil.Update(r.DB, r.TableName, animal)
}

func (r *AnimalRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.DB, r.TableName, id)
}
