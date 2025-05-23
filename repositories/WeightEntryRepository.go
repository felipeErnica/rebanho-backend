package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type WeightEntryRepository struct {
	Impl RepositoryImpl[entity.WeightEntry]
}

func (r *WeightEntryRepository) Init() {

	selectQuery := util.NewSelectQuery(util.SELECT,
		*util.NewNamedGroup("weight", "id", "group_id", "weight"),
		*util.NewNamedGroup("animals", "id", "name", "identification_number", "animal_order")).
		From("weight_entries as weight").
		Joins("animals on animals.id = weight_entries.animal_id")

	updateQuery := util.NewUpdateQuery("weight_entries", "animal_id", "group_id", "weight", "created_at")
	insertQuery := util.NewInsertQuery("weight_entries", "id", "animal_id", "group_id", "weight", "created_at")

	r.Impl = RepositoryImpl[entity.WeightEntry]{
		Repository:      r,
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		TableName:       "weight_entries",
	}

}

func (r *WeightEntryRepository) setNewEntity(model *entity.WeightEntry, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
	model.UserId = GetUserId()
}

func (r *WeightEntryRepository) buildEntity(row *sql.Row) (model *entity.WeightEntry, err error) {
	var weight entity.WeightEntry
	err = row.Scan(weight.Id, weight.GroupId, weight.Weight,
		weight.AnimalId, weight.AnimalName, weight.AnimalNumber, weight.AnimalOrder)
	return &weight, err
}

func (r *WeightEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.WeightEntry, err error) {
	var weights []entity.WeightEntry
	for rows.Next() {
		var weight entity.WeightEntry
		err = rows.Scan(weight.Id, weight.GroupId, weight.Weight,
			weight.AnimalId, weight.AnimalName, weight.AnimalNumber, weight.AnimalOrder)
		if err != nil {
			return
		}
		weights = append(weights, weight)
	}
	return &weights, err
}

func (r *WeightEntryRepository) saveOrUpdateScan(query string, model *entity.WeightEntry) error {
	return execQuery(query, model.Id, model.AnimalId, model.GroupId, model.Weight, model.CreatedAt, model.UserId)
}

func (r *WeightEntryRepository) FindByGroupId(groupId string) (*[]entity.WeightEntry, error) {
	query := r.Impl.SelectQueryBody.Where("weight.group_id = $1 and weight.deleted_at is null")
	return r.Impl.FindListByQuery(query, groupId)
}

func (r *WeightEntryRepository) FindByAnimalId(animalId string) (*[]entity.WeightEntry, error) {
	query := r.Impl.SelectQueryBody.Where("weight.animal_id = $1 and weight.deleted_at is null")
	return r.Impl.FindListByQuery(query, animalId)
}

func (r *WeightEntryRepository) FindById(id string) (*entity.WeightEntry, error) {
	return r.Impl.FindById(id)
}

func (r *WeightEntryRepository) Add(newModel *entity.WeightEntry) (*entity.WeightEntry, error) {
	return r.Impl.Add(newModel)
}

func (r *WeightEntryRepository) Save(model *entity.WeightEntry) error {
	return r.Impl.Save(model)
}

func (r *WeightEntryRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
