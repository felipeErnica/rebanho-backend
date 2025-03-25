package insemination

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity/insemination"
	"github.com/felipeErnica/rebanho-backend/repositories"
	"github.com/felipeErnica/rebanho-backend/util"
)

type InseminationEntryRepository struct {
	Impl repositories.RepositoryImpl[insemination.InseminationEntry]
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

    r.Impl = repositories.RepositoryImpl[insemination.InseminationEntry]{
        TableName: "insemination_entries",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *InseminationEntryRepository) SetNewEntity(model *insemination.InseminationEntry, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *InseminationEntryRepository) BuildEntity(row *sql.Row) (model *insemination.InseminationEntry, err error) {
    var entry insemination.InseminationEntry
    err = row.Scan(&entry.Id, &entry.GroupId, &entry.Observation, &entry.Status, 
        &entry.Animal.Id, &entry.Animal.Name, &entry.Animal.IdentificationNumber, &entry.Animal.AnimalOrder, 
        &entry.Loss.Id, &entry.Loss.LossType, &entry.Loss.LossDate,
        &entry.Calf.Id, &entry.Calf.Sex, &entry.Calf.BirthDate)
    return &entry, err
}

func (r *InseminationEntryRepository) BuildListEntity(rows *sql.Rows) (arr *[]insemination.InseminationEntry, err error) {
    var entries []insemination.InseminationEntry
    for rows.Next() {
        var entry insemination.InseminationEntry
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

func (r *InseminationEntryRepository) SaveOrUpdateScan(query string, model *insemination.InseminationEntry) error {
    return repositories.ExecQuery(query, model.Id, model.Animal.Id, model.GroupId, model.Observation, 
        model.Status, model.Loss.Id, model.CreatedAt)
}

func (r *InseminationEntryRepository) FindByGroupId(groupId string) (*[]insemination.InseminationEntry, error) {
    query:="WHERE entry.group_id = $1"
    return r.Impl.FindListByQuery(query, groupId)
}

func (r *InseminationEntryRepository) FindById(id string) (*insemination.InseminationEntry, error) {
    return r.Impl.FindById(id)
}

func (r *InseminationEntryRepository) Add(newModel insemination.InseminationEntry) (*insemination.InseminationEntry, error) {
    return r.Impl.Add(newModel)
}

func (r *InseminationEntryRepository) Save(model *insemination.InseminationEntry) error {
    return r.Impl.Save(model)
}

func (r *InseminationEntryRepository) Delete(id string) error {
    return r.Impl.SoftDelete(id)
}
