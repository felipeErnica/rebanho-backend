package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type BirthEntryRepository struct {
	Impl RepositoryImpl[entity.BirthEntry]
}

func (r *BirthEntryRepository) Init() {
    selectQuery:=new(util.QueryConstructor).Select("", "id", "observation", "calf_id")
        selectQuery.From("birth_entries", "")
    insertQuery:=new(util.QueryConstructor).Insert("birth_entries", "id", "calf_id", "observation")
    updateQuery:=new(util.QueryConstructor).Update("birth_entries", "id", "calf_id", "observation")
    r.Impl = RepositoryImpl[entity.BirthEntry]{
        TableName: "birth_entries",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *BirthEntryRepository) setNewEntity(model *entity.BirthEntry, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *BirthEntryRepository) buildEntity(row *sql.Row) (model *entity.BirthEntry, err error) {
    var entry entity.BirthEntry
    err = row.Scan(&entry.Id, &entry.Observation, &entry.CalfId)
    return &entry, err
}

func (r *BirthEntryRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.BirthEntry, err error) {
    var entries []entity.BirthEntry
    for rows.Next() {
        var entry entity.BirthEntry
        err = rows.Scan(&entry.Id, &entry.Observation, &entry.CalfId)
        if err != nil {
            return
        }
        entries = append(entries, entry)
    }
    return &entries, err
}

func (r *BirthEntryRepository) saveOrUpdateScan(query string, model *entity.BirthEntry) error {
    return execQuery(query, model.Id, model.CalfId, model.Observation)
}

func (r *BirthEntryRepository) Delete(id string) error {
    return r.Impl.HardDelete(id)
}

