package birth

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type BirthRepository struct {
	SelectQuery string
	TableName   string
	Db          *sqlx.DB
}

func NewRepository(db *sqlx.DB) *BirthRepository {
	selectQuery := `
        SELECT birth.*, 
            animals.name as animal_name, animals.number as animal_number, animals.order as animal_order,
            calf.birth_date as calf_birth_date, calf.sex as calf_sex, father.name as calf_father
        FROM birth_entries as birth 
            LEFT JOIN animals ON animals.id = birth.animal_id
            LEFT JOIN animals as calf ON calf.id = birth.calf_id
            LEFT JOIN animals as father ON father.id = calf.father_id
    `
    return &BirthRepository{selectQuery, "birth", db}
}

func (r *BirthRepository) FindByMotherId(motherId string) (*[]BirthEntry, error) {
	query := r.SelectQuery + " WHERE birth.animal_id = $1 AND birth.deleted_at is null"
	return repositoriesUtil.GetList[BirthEntry](r.Db, query, motherId)
}

func (r *BirthRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
