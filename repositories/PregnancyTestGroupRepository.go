package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PregnancyTestGroupRepository struct {
    Impl RepositoryImpl[entity.PregancyTestGroup]
}

func (r *PregnancyTestGroupRepository) Init() {
    selectQuery:=new(util.SelectConstructor).Select("", "id", "test_date", "number_entries", "number_pregnants")
        selectQuery.From("pregnancy_test_groups", "")
    insertQuery:=new(util.SelectConstructor).Insert("pregnancy_test_groups", "id", "test_date", "number_entries", "number_pregnants")
    updateQuery:=new(util.SelectConstructor).Update("pregnancy_test_groups", "id", "test_date", "number_entries", "number_pregnants")
    r.Impl = RepositoryImpl[entity.PregancyTestGroup]{
        TableName: "pregnancy_test_groups",
        SelectQueryBody: *selectQuery,
        InsertQuery: *insertQuery,
        UpdateQuery: *updateQuery,
        Repository: r,
    }
}

func (r *PregnancyTestGroupRepository) setNewEntity(model *entity.PregancyTestGroup, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
    model.UserId = GetUserId()
}

func (r *PregnancyTestGroupRepository) buildEntity(row *sql.Row) (model *entity.PregancyTestGroup, err error) {
    var test entity.PregancyTestGroup
    err = row.Scan(&test.Id, &test.TestDate, &test.NumberEntries, &test.NumberPregnants)
    return &test, err
}

func (r *PregnancyTestGroupRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.PregancyTestGroup, err error) {
    var tests  []entity.PregancyTestGroup
    for rows.Next() {
        var test entity.PregancyTestGroup
        err = rows.Scan(&test.Id, &test.TestDate, &test.NumberEntries, &test.NumberPregnants)
        if err != nil {
            return
        }
        tests = append(tests, test)
    }
    return &tests, err
}

func (r *PregnancyTestGroupRepository) saveOrUpdateScan(query string, model *entity.PregancyTestGroup) error {
    return execQuery(query, model.Id, model.TestDate, model.NumberEntries, model.NumberPregnants, model.CreatedAt, model.UserId)
}

func (r *PregnancyTestGroupRepository) FindAll() (*[]entity.PregancyTestGroup, error) {
    return r.Impl.FindAll()
} 

func (r *PregnancyTestGroupRepository) FindById(id string) (*entity.PregancyTestGroup, error) {
    return r.Impl.FindById(id)
} 

func (r *PregnancyTestGroupRepository) Add(newModel *entity.PregancyTestGroup) (*entity.PregancyTestGroup, error) {
    return r.Impl.Add(newModel)
} 

func (r *PregnancyTestGroupRepository) Save(model *entity.PregancyTestGroup) error {
    return r.Impl.Save(model)
} 

func (r *PregnancyTestGroupRepository) Delete(id string) error {
    return r.Impl.Delete(id)
} 

