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
    selectQuery:=new(util.QueryConstructor).Select("entry", "id", "group_id", "observation", "status")
        selectQuery.AndSelect("animals", "id", "name", "identification_number", "animal_order")
        selectQuery.AndSelect("loss", "id", "loss_type", "loss_date")
        selectQuery.AndSelect("calf", "id", "sex", "birth_date")
        selectQuery.From("insemination_entries_active", "entry")
        selectQuery.LeftJoin("animals","").On("animals.id", "entry.animal_id")
        selectQuery.LeftJoin("pregnancy_losses","loss").On("loss.id", "entry.loss_id")
    insertQuery:=new(util.QueryConstructor).Insert("insemination_entries", "id", "animal_id", 
        "group_id", "observation", "status", "loss_id", "calf_id", "created_at")
    updateQuery:=new(util.QueryConstructor).Update("insemination_entries", "animal_id", 
        "group_id", "observation", "status", "loss_id", "calf_id", "created_at")
    r.Impl = RepositoryImpl[entity.InseminationEntry]{
        TableName: "insemination_entries",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *InseminationEntryRepository) setNewEntity(model *entity.InseminationEntry, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *InseminationEntryRepository) buildEntity(row *sql.Row) (model *entity.InseminationEntry, err error) {
    var entry entity.InseminationEntry
    err = row.Scan(&entry.Id, &entry.GroupId, &entry.Observation, &entry.Status, 
        &entry.Animal.Id, &entry.Animal.Name, &entry.Animal.IdentificationNumber, &entry.Animal.AnimalOrder, 
        &entry.Loss.Id, &entry.Loss.LossType, &entry.Loss.LossDate,
        &entry.Calf.Id, &entry.Calf.Sex, &entry.Calf.BirthDate)
    return &entry, err
}

func (r *InseminationEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.InseminationEntry, err error) {
    var entries []entity.InseminationEntry
    for rows.Next() {
        var entry entity.InseminationEntry
        err = rows.Scan(&entry.Id, &entry.GroupId, &entry.Observation, &entry.Status, 
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

func (r *InseminationEntryRepository) saveOrUpdateScan(query string, model *entity.InseminationEntry) error {
    return execQuery(query, model.Id, model.Animal.Id, model.GroupId, model.Observation, 
        model.Status, model.Loss.Id, model.CreatedAt)
}

func (r *InseminationEntryRepository) FindByGroupId(groupId string) (*[]entity.InseminationEntry, error) {
    query:="WHERE entry.group_id = $1"
    return r.Impl.FindListByQuery(query, groupId)
}

func (r *InseminationEntryRepository) FindById(id string) (*entity.InseminationEntry, error) {
    return r.Impl.FindById(id)
}

func (r *InseminationEntryRepository) Add(newModel entity.InseminationEntry) (*entity.InseminationEntry, error) {
    return r.Impl.Add(newModel)
}

func (r *InseminationEntryRepository) Save(model *entity.InseminationEntry) error {
    return r.Impl.Save(model)
}

func (r *InseminationEntryRepository) Delete(id string) error {
    return r.Impl.SoftDelete(id)
}
