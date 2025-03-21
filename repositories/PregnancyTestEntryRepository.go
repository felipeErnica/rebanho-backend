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
    selectQuery:=new(util.QueryConstructor).Select("entry", "id", "group_id", "is_pregnant", "birth_forecast")
        selectQuery.AndSelect("animals", "id", "name", "identification_number", "animal_order")
        selectQuery.From("pregancy_test_entries", "entry")
        selectQuery.LeftJoin("animals", "").On("animals.id", "entry.animal_id")
    insertQuery:=new(util.QueryConstructor).Insert("pregancy_test_entries", "id", "animal_id", "group_id", 
        "is_pregnant", "birth_forecast", "created_at", "user_id")
    updateQuery:=new(util.QueryConstructor).Update("pregancy_test_entries", "id", "animal_id", "group_id", 
        "is_pregnant", "birth_forecast", "created_at", "user_id")
    r.Impl = RepositoryImpl[entity.PregnancyTestEntry]{
        TableName: "pregancy_test_entries",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *PregancyTestEntryRepository) setNewEntity(model *entity.PregnancyTestEntry, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
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
            &entry.Animal.Id, &entry.Animal.Name, &entry.Animal.IdentificationNumber, &entry.Animal.AnimalOrder)
        if err != nil {
            return
        }
        entries = append(entries, entry)
    }
    return &entries, err
}

func (r *PregancyTestEntryRepository) saveOrUpdateScan(query string, model *entity.PregnancyTestEntry) error {
    return execQuery(query, model.Id, model.Animal.Id, model.GroupId, model.IsPregnant, 
        model.BirthForecast, model.CreatedAt)
}

func (r *PregancyTestEntryRepository) FindByGroupId(groupId string) (*[]entity.PregnancyTestEntry, error) {
    query:="WHERE entry.group_id = $1"
    return r.Impl.FindListByQuery(query, groupId)
}

func (r *PregancyTestEntryRepository) FindById(id string) (*entity.PregnancyTestEntry, error) {
    return r.Impl.FindById(id)
}

func (r *PregancyTestEntryRepository) Add(newModel entity.PregnancyTestEntry) (*entity.PregnancyTestEntry, error) {
    return r.Impl.Add(newModel)
}

func (r *PregancyTestEntryRepository) Save(model *entity.PregnancyTestEntry) error {
    return r.Impl.Save(model)
}

func (r *PregancyTestEntryRepository) Delete(id string) error {
    return r.Impl.HardDelete(id)
}
