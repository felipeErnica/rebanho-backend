package insemination

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity/insemination"
	"github.com/felipeErnica/rebanho-backend/repositories"
	"github.com/felipeErnica/rebanho-backend/util"
)

type InseminationGroupRepository struct {
	Impl repositories.RepositoryImpl[insemination.InseminationGroup]
}

func (r *InseminationGroupRepository) Init() {
    selectQuery:=new(util.QueryConstructor).Select("group", "id", "insemination_date")
        selectQuery.AndSelect("bull", "id", "name")
        selectQuery.AndSelect("bull_mother", "id", "name")
        selectQuery.AndSelect("bull_father", "id", "name")
        selectQuery.From("insemination_groups", "group")
        selectQuery.LeftJoin("animals", "bull").On("bull.id", "group.bull_id")
        selectQuery.LeftJoin("animals", "bull_mother").On("bull_mother.id", "bull.mother_id")
        selectQuery.LeftJoin("animals", "bull_father").On("bull_father.id", "bull.father_id")
    insertQuery:= new(util.QueryConstructor).Insert("insemination_groups", "id", "bull_id", 
        "insemination_date", "created_at", "user_id")
    updateQuery:= new(util.QueryConstructor).Update("insemination_groups", "bull_id", 
        "insemination_date", "created_at", "user_id")
    r.Impl = repositories.RepositoryImpl[insemination.InseminationGroup]{
        TableName: "insemination_groups",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *InseminationGroupRepository) SetNewEntity(model *insemination.InseminationGroup, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *InseminationGroupRepository) BuildEntity(row *sql.Row) (model *insemination.InseminationGroup, err error) {
    var group insemination.InseminationGroup
    err = row.Scan(&group.Id, &group.InseminationDate, &group.Bull.Id, &group.Bull.Name, 
        &group.Bull.Mother.Id, &group.Bull.Mother.Name, 
        &group.Bull.Father.Id, &group.Bull.Father.Name)
    return &group, err
}

func (r *InseminationGroupRepository) BuildListEntity(rows *sql.Rows) (arr *[]insemination.InseminationGroup, err error) {
    var groups []insemination.InseminationGroup
    for rows.Next() {
        var group insemination.InseminationGroup
        err = rows.Scan(&group.Id, &group.InseminationDate, &group.Bull.Id, &group.Bull.Name, 
            &group.Bull.Mother.Id, &group.Bull.Mother.Name, 
            &group.Bull.Father.Id, &group.Bull.Father.Name)
        if err != nil {
            return
        }
        groups = append(groups, group)
    }
    return &groups, err
}

func (r *InseminationGroupRepository) SaveOrUpdateScan(query string, model *insemination.InseminationGroup) error {
    return repositories.ExecQuery(query, model.Id, model.Bull.Id, model.InseminationDate, model.CreatedAt, model.UserId)
}

func (r *InseminationGroupRepository) FindAll() (*[]insemination.InseminationGroup, error) {
    return r.Impl.FindAll()
}

func (r *InseminationGroupRepository) FindById(id string) (*insemination.InseminationGroup, error) {
    return r.Impl.FindById(id)
}

func (r *InseminationGroupRepository) Add(newModel insemination.InseminationGroup) (*insemination.InseminationGroup, error) {
    return r.Impl.Add(newModel)
}

func (r *InseminationGroupRepository) Save(model *insemination.InseminationGroup) error {
    return r.Impl.Save(model)
}

func (r *InseminationGroupRepository) Delete(id string) error {
    return r.Impl.HardDelete(id)
}
