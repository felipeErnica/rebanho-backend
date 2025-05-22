package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type InseminationEntryRepository struct {
	Impl RepositoryImpl[entity.InseminationEntry]
}

func (r *InseminationEntryRepository) Init() {
	selectQuery := util.NewSelectQuery(util.SELECT, 
        *util.NewNamedGroup("entry", "id", "group_id", "observation", "status"),
	    *util.NewNamedGroup("animals", "id", "name", "identification_number", "animal_order"),
	    *util.NewNamedGroup("loss", "id", "loss_type", "loss_date"),
	    *util.NewNamedGroup("calf", "id", "sex", "birth_date")).
        From("insemination_entries as entry").
        Joins(
            "left join animals on animals.id = entry.animal_id",
            "left join pregnancy_losses as loss on loss.id = entry.loss_id")

    insertQuery := util.NewInsertQuery("insemination_entries", "id", "animal_id",
		"group_id", "observation", "status", "loss_id", "calf_id", "created_at")
    updateQuery := util.NewUpdateQuery("insemination_entries", "animal_id",
		"group_id", "observation", "status", "loss_id", "calf_id", "created_at")

	r.Impl = RepositoryImpl[entity.InseminationEntry]{
		TableName:       "insemination_entries",
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		Repository:      r,
	}
}

func (r *InseminationEntryRepository) setNewEntity(model *entity.InseminationEntry, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
}

func (r *InseminationEntryRepository) buildEntity(row *sql.Row) (model *entity.InseminationEntry, err error) {
	var entry entity.InseminationEntry
	err = row.Scan(&entry.Id, &entry.GroupId, &entry.Observation, &entry.Status,
		&entry.AnimalId, &entry.AnimalName, &entry.AnimalNumber, &entry.AnimalOrder,
		&entry.LossId, &entry.CalfId)
	return &entry, err
}

func (r *InseminationEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.InseminationEntry, err error) {
	var entries []entity.InseminationEntry
	for rows.Next() {
		var entry entity.InseminationEntry
		err = rows.Scan(&entry.Id, &entry.GroupId, &entry.Observation, &entry.Status,
			&entry.AnimalId, &entry.AnimalName, &entry.AnimalNumber, &entry.AnimalOrder,
			&entry.LossId, &entry.CalfId)
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}
	return &entries, err
}

func (r *InseminationEntryRepository) saveOrUpdateScan(query string, model *entity.InseminationEntry) error {
	return execQuery(query, model.Id, model.AnimalId, model.GroupId, model.Observation,
		model.Status, model.LossId, model.CalfId, model.CreatedAt)
}

func (r *InseminationEntryRepository) FindByGroupId(groupId string) (*[]entity.InseminationEntry, error) {
	query := r.Impl.SelectQueryBody.Where("groups.id = $1")
	return r.Impl.FindListByQuery(query, groupId)
}

func (r *InseminationEntryRepository) FindById(id string) (*entity.InseminationEntry, error) {
	return r.Impl.FindById(id)
}

func (r *InseminationEntryRepository) Add(newModel *entity.InseminationEntry) (*entity.InseminationEntry, error) {
	return r.Impl.Add(newModel)
}

func (r *InseminationEntryRepository) Save(model *entity.InseminationEntry) error {
	return r.Impl.Save(model)
}

func (r *InseminationEntryRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
