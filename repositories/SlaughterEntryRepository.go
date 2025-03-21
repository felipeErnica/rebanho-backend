package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type SlaughterEntryRepository struct {
    Impl RepositoryImpl[entity.SlaughterEntry]
}

func (r *SlaughterEntryRepository) Init() {
    selectQuery:=new(util.QueryConstructor).Select("entry", "id", "weight", "dead_weight")
        selectQuery.AndSelect("group", "id", "weight_decrease", "slaugther_date")
        selectQuery.From("slaughter_entries", "entry")
        selectQuery.LeftJoin("slaughter_groups", "group").On("group.id", "entry.group_id")
    insertQuery:=new(util.QueryConstructor).Insert("slaughter_entries", "id", "group_id", "weight", 
        "dead_weight", "created_at", "user_id")
    updateQuery:=new(util.QueryConstructor).Insert("slaughter_entries", "id", "group_id", "weight", 
        "dead_weight", "created_at")
    r.Impl = RepositoryImpl[entity.SlaughterEntry]{
        TableName: "slaughter_entries",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *SlaughterEntryRepository) setNewEntity(model *entity.SlaughterEntry, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *SlaughterEntryRepository) buildEntity(row *sql.Row) (model *entity.SlaughterEntry, err error) {
    var entry entity.SlaughterEntry
    err = row.Scan(&entry.Id, &entry.Weight, &entry.DeadWeight, &entry.Group.Id, &entry.Group.WeightDecrease, &entry.Group.SlaughterDate)
    return &entry, err
}

func (r *SlaughterEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.SlaughterEntry, err error) {
    var entries []entity.SlaughterEntry
    for rows.Next() {
        var entry entity.SlaughterEntry
        err = rows.Scan(&entry.Id, &entry.Weight, &entry.DeadWeight, &entry.Group.Id, &entry.Group.WeightDecrease, &entry.Group.SlaughterDate)
        if err != nil {
            return
        }
        entries = append(entries, entry)
    }
    return &entries, err
}

func (r *SlaughterEntryRepository) saveOrUpdateScan(query string, model *entity.SlaughterEntry) error {
    return execQuery(query, model.Id, model.Group.Id, model.Weight, model.DeadWeight, model.CreatedAt)
}

func (r *SlaughterEntryRepository) FindByGroupId(groupId string) (*[]entity.SlaughterEntry, error) {
    query:="WHERE entry.group_id = $1"
    return r.Impl.FindListByQuery(query, groupId)
}

func (r *SlaughterEntryRepository) FindById(id string) (*entity.SlaughterEntry, error) {
    return r.Impl.FindById(id)
}

func (r *SlaughterEntryRepository) Add(newModel entity.SlaughterEntry) (*entity.SlaughterEntry, error) {
    return r.Impl.Add(newModel)
}

func (r *SlaughterEntryRepository) Save(model *entity.SlaughterEntry) error {
    return r.Impl.Save(model)
}

func (r *SlaughterEntryRepository) Delete(id string) error {
    return r.Impl.HardDelete(id)
}
