package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PregancyTestEntryRepository struct {
	Impl RepositoryImpl[entity.PregnancyTestEntry]
}

func (r *PregancyTestEntryRepository) Init() {
	selectQuery := util.NewSelectQuery(util.SELECT, 
        *util.NewNamedGroup("entry", "id", "group_id", "is_pregnant", "birth_forecast"),
        *util.NewNamedGroup("animal", "id", "name", "identification_number", "animal_order"),
        *util.NewNamedGroup("loss", "id", "loss_type", "loss_date"),
        *util.NewNamedGroup("calf", "id", "sex", "birth_date")).
        From("pregancy_test_entries as entry").
        Joins(
            "left join animals as animal on animal.id = entry.animal_id",
            "left join animals as calf on calf.id = entry.calf_id",
            "left join pregnancy_losses as loss on loss.id = entry.loss_id")

	insertQuery := util.NewInsertQuery("pregancy_test_entries", "id", "animal_id", "group_id",
		"is_pregnant", "birth_forecast", "loss_id", "calf_id", "created_at", "user_id")
	updateQuery := util.NewUpdateQuery("pregancy_test_entries", "id", "animal_id", "group_id",
		"is_pregnant", "birth_forecast", "loss_id", "calf_id", "created_at", "user_id")

	r.Impl = RepositoryImpl[entity.PregnancyTestEntry]{
		TableName:       "pregancy_test_entries",
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		Repository:      r,
	}
}

func (r *PregancyTestEntryRepository) setNewEntity(model *entity.PregnancyTestEntry, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
	model.UserId = GetUserId()
}

func (r *PregancyTestEntryRepository) buildEntity(row *sql.Row) (model *entity.PregnancyTestEntry, err error) {
	var entry entity.PregnancyTestEntry
	err = row.Scan(&entry.Id, &entry.GroupId, &entry.IsPregnant, &entry.BirthForecast,
		&entry.Animal.Id, &entry.Animal.Name, &entry.Animal.IdentificationNumber, &entry.Animal.AnimalOrder)
	return &entry, err
}

func (r *PregancyTestEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.PregnancyTestEntry, err error) {
	var entries []entity.PregnancyTestEntry
	for rows.Next() {
		var entry entity.PregnancyTestEntry
		err = rows.Scan(&entry.Id, &entry.GroupId, &entry.IsPregnant, &entry.BirthForecast,
			&entry.Animal.Id, &entry.Animal.Name, &entry.Animal.IdentificationNumber, &entry.Animal.AnimalOrder,
			&entry.Loss.Id, &entry.Loss.LossType, &entry.Loss.LossDate,
			&entry.Calf.Id, &entry.Calf.Sex, &entry.Calf.BirthDate)
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}
	return &entries, err
}

func (r *PregancyTestEntryRepository) saveOrUpdateScan(query string, model *entity.PregnancyTestEntry) error {
	return execQuery(query, model.Id, model.Animal.Id, model.GroupId, model.IsPregnant,
		model.BirthForecast, model.Calf.Id, model.Loss.Id, model.CreatedAt, model.UserId)
}

func (r *PregancyTestEntryRepository) FindByGroupId(groupId string) (*[]entity.PregnancyTestEntry, error) {
	query := r.Impl.SelectQueryBody.Where("entry.group_id = $1 and entry.deleted_at is null")
	return r.Impl.FindListByQuery(query, groupId)
}

func (r *PregancyTestEntryRepository) FindByAnimalId(animalId string) (*[]entity.PregnancyTestEntry, error) {
	query := r.Impl.SelectQueryBody.Where("entry.animal_id = $1 and entry.deleted_at is null")
	return r.Impl.FindListByQuery(query, animalId)
}

func (r *PregancyTestEntryRepository) FindById(id string) (*entity.PregnancyTestEntry, error) {
	return r.Impl.FindById(id)
}

func (r *PregancyTestEntryRepository) Add(newModel *entity.PregnancyTestEntry) (*entity.PregnancyTestEntry, error) {
	return r.Impl.Add(newModel)
}

func (r *PregancyTestEntryRepository) Save(model *entity.PregnancyTestEntry) error {
	return r.Impl.Save(model)
}

func (r *PregancyTestEntryRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
