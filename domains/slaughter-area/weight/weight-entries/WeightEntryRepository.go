package weightEntries

import (
	repositoriesUtil "github.com/felipeErnica/rebanho-backend/util/repositories-util"
	"github.com/jmoiron/sqlx"
)

type WeightEntryRepository struct {
    SelectQuery string
    TableName string
    Db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *WeightEntryRepository {
    selectQuery := `
        SELECT weight_entries.*,
            animals.name as animal_name, animals.number as animal_number, 
            animals.order as animal_order, animals.sex as animal_sex,
            group.date as group_date
        LEFT JOIN animals ON animals.id = weight_entries.animal_id 
        LEFT JOIN weight_groups ON weight_groups.id = weight_entries.group_id 
    `
    return &WeightEntryRepository{selectQuery, "weight_entries", db}
}

func (r *WeightEntryRepository) FindByGroupId(groupId string) (*[]WeightEntry, error) {
	query := r.SelectQuery + " WHERE weight_entries.group_id = $1 and weight_entries.deleted_at is null"
	return repositoriesUtil.GetList[WeightEntry](r.Db, query, groupId)
}

func (r *WeightEntryRepository) FindByAnimalId(animalId string) (*[]WeightEntry, error) {
	query := r.SelectQuery + " WHERE weight_entries.animal_id = $1 and weight_entries.deleted_at is null"
	return repositoriesUtil.GetList[WeightEntry](r.Db, query, animalId)
}

func (r *WeightEntryRepository) FindById(id string) (*WeightEntry, error) {
	query := r.SelectQuery + " WHERE weight_entries.id = $1 and weight_entries.deleted_at is null"
	return repositoriesUtil.GetOne[WeightEntry](r.Db, query, id)
}

func (r *WeightEntryRepository) Add(newModel *WeightEntry) (*WeightEntry, error) {
    return repositoriesUtil.Add(r.Db, r.TableName, newModel)
}

func (r *WeightEntryRepository) Update(model *WeightEntry) error {
    return repositoriesUtil.Update(r.Db, r.TableName, model)
}

func (r *WeightEntryRepository) Delete(id string) error {
	return repositoriesUtil.Delete(r.Db, r.TableName, id)
}
