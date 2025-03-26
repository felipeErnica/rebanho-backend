package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type InseminationGroupRepository struct {
	Impl RepositoryImpl[entity.InseminationGroup]
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
    r.Impl = RepositoryImpl[entity.InseminationGroup]{
        TableName: "insemination_groups",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *InseminationGroupRepository) setNewEntity(model *entity.InseminationGroup, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *InseminationGroupRepository) buildEntity(row *sql.Row) (model *entity.InseminationGroup, err error) {
    var group entity.InseminationGroup
    err = row.Scan(&group.Id, &group.InseminationDate, &group.Bull.Id, &group.Bull.Name, 
        &group.Bull.Mother.Id, &group.Bull.Mother.Name, 
        &group.Bull.Father.Id, &group.Bull.Father.Name)
    return &group, err
}

func (r *InseminationGroupRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.InseminationGroup, err error) {
    var groups []entity.InseminationGroup
    for rows.Next() {
        var group entity.InseminationGroup
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

func (r *InseminationGroupRepository) saveOrUpdateScan(query string, model *entity.InseminationGroup) error {
    return execQuery(query, model.Id, model.Bull.Id, model.InseminationDate, model.CreatedAt, model.UserId)
}

func (r *InseminationGroupRepository) FindAll() (*[]entity.InseminationGroup, error) {
    return r.Impl.FindAll()
}

func (r *InseminationGroupRepository) FindById(id string) (*entity.InseminationGroup, error) {
    return r.Impl.FindById(id)
}

func (r *InseminationGroupRepository) Add(newModel entity.InseminationGroup) (*entity.InseminationGroup, error) {
    return r.Impl.Add(newModel)
}

func (r *InseminationGroupRepository) Save(model *entity.InseminationGroup) error {
    return r.Impl.Save(model)
}

func (r *InseminationGroupRepository) Delete(id string) error {
    return r.Impl.HardDelete(id)
}
