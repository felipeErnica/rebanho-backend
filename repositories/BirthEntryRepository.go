package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type BirthEntryRepository struct {
	Impl            PageRepositoryImpl[entity.BirthEntry]
	SelectQueryBody util.SelectConstructor
}

func (r *BirthEntryRepository) Init() {
	selectQuery := util.NewSelectQuery(util.SELECT,
		*util.NewNamedGroup("birth", "id", "observation"),
		*util.NewNamedGroup("calf", "id", "name", "identification_number", "sex", "birth_date"),
		*util.NewNamedGroup("mother", "id", "name", "identification_number", "animal_order"),
		*util.NewNamedGroup("father", "id", "name")).
		From("birth_entries")

	insertQuery := util.NewInsertQuery("birth_entries", "id", "animal_id", "calf_id", "observation")
	updateQuery := util.NewUpdateQuery("birth_entries", "id", "animal_id", "calf_id", "observation")

	base := RepositoryImpl[entity.BirthEntry]{
		TableName:       "birth_entries",
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		Repository:      r,
	}
	r.Impl = PageRepositoryImpl[entity.BirthEntry]{
		Base: &base,
	}
}

func (r *BirthEntryRepository) setNewEntity(model *entity.BirthEntry, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
}

func (r *BirthEntryRepository) buildEntity(row *sql.Row) (model *entity.BirthEntry, err error) {
	var entry entity.BirthEntry
	err = row.Scan(&entry.Id, &entry.Observation,
		&entry.CalfId, &entry.CalfSex, &entry.CalfBirthDate,
		&entry.CalfFatherName,
		&entry.CalfId, &entry.AnimalName, &entry.AnimalNumber)
	return &entry, err
}

func (r *BirthEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.BirthEntry, err error) {
	var entries []entity.BirthEntry
	for rows.Next() {
		var entry entity.BirthEntry
		err = rows.Scan(&entry.Id, &entry.Observation,
			&entry.CalfId, &entry.CalfSex, &entry.CalfBirthDate,
			&entry.CalfFatherName,
			&entry.CalfId, &entry.AnimalName, &entry.AnimalNumber)
		if err != nil {
			return
		}
		entries = append(entries, entry)
	}
	return &entries, err
}

func (r *BirthEntryRepository) saveOrUpdateScan(query string, model *entity.BirthEntry) error {
	return execQuery(query, model.Id, model.AnimalId, model.CalfId, model.Observation)
}

func (r *BirthEntryRepository) FindByMotherId(motherId string) (*[]entity.BirthEntry, error) {
	query := r.SelectQueryBody.Where("mother.id = $1")
	return r.Impl.FindListByQuery(query, motherId)
}

func (r *BirthEntryRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
