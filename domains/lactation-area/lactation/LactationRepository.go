package lactation

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type LactationRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *LactationRepository {
	selectQuery := `
	SELECT lactations.*, 
		cow.name as cow_name, cow.number as cow_number, cow.pasture as cow_pasture, cow.order as cow_order,
		calf.birth_data as calf_birth_date, calf.sex as calf_sex, calf.father as calf_father
	FROM lactations
		LEFT JOIN cow as animals ON cow.id = lactations.cow_id
		LEFT JOIN calf as animals ON calf.id = lactations.calf_id 
	`
	return &LactationRepository{selectQuery, "lactations", db}
}

func (r *LactationRepository) FindByAnimal(animalId string) (arr *[]Lactation, err error) {
	query := r.SelectQuery + " WHERE lactations.deleted_at is null and lactations.animal_id = $1"
	return repositoriesUtil.GetList[Lactation](r.Db, query, animalId)
}

func (r *LactationRepository) Add(newLactation *LactationSave) (*LactationSave, error) {
	return repositoriesUtil.Add(r.Db, r.TableName, newLactation)
}

func (r *LactationRepository) Update(lactation *LactationSave) error {
	return repositoriesUtil.Update(r.Db, r.TableName, lactation)
}

func (r *LactationRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
