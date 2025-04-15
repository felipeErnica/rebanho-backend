package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type SlaughterGroupRepository struct {
    Impl  RepositoryImpl[entity.SlaughterGroup]
}

func (r *SlaughterGroupRepository) Init() {
    selectQuery:=new(util.SelectConstructor).Select("group", "id", "weight_decrease", "slaughter_date")
        selectQuery.AndSelect("slaughterhouses", "id", "name")
        selectQuery.From("slaughter_groups", "group")
        selectQuery.LeftJoin("slaughterhouses","").On("slaughterhouses.id", "group.slaughterhouse_id")
    updateQuery:=new(util.SelectConstructor).Update("slaughter_groups", "weight_decrease", "slaughter_date", 
        "slaughterhouse_id", "created_at", "user_id")
    insertQuery:=new(util.SelectConstructor).Insert("slaughter_groups", "id", "weight_decrease", "slaughter_date", 
        "slaughterhouse_id", "created_at", "user_id")
    
    r.Impl = RepositoryImpl[entity.SlaughterGroup]{
        TableName: "slaughter_groups",
        SelectQueryBody: *selectQuery,
        InsertQuery: *insertQuery,
        UpdateQuery: *updateQuery,
    }
}

func (r *SlaughterGroupRepository) setNewEntity(model *entity.SlaughterGroup, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
    model.UserId = GetUserId()
}

func (r *SlaughterGroupRepository) buildEntity(row *sql.Row) (model *entity.SlaughterGroup, err error) {
    var group entity.SlaughterGroup
    err = row.Scan(&group.Id, &group.WeightDecrease, &group.SlaughterDate, &group.Slaughterhouse.Id, &group.Slaughterhouse.Name)
    return &group, err
}

func (r *SlaughterGroupRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.SlaughterGroup, err error) {
    var groups []entity.SlaughterGroup
    for rows.Next() {
        var group entity.SlaughterGroup
        err = rows.Scan(&group.Id, &group.WeightDecrease, &group.SlaughterDate, &group.Slaughterhouse.Id, &group.Slaughterhouse.Name)
        if err != nil {
            return
        }
        groups = append(groups, group)
    }
    return &groups, err
}

func (r *SlaughterGroupRepository) saveOrUpdateScan(query string, model *entity.SlaughterGroup) error {
    return execQuery(query, model.Id, model.WeightDecrease, model.SlaughterDate, 
        model.Slaughterhouse.Id, model.CreatedAt, model.UserId)
}

func (r *SlaughterGroupRepository) FindAll() (*[]entity.SlaughterGroup, error) {
    return r.Impl.FindAll()
}

func (r *SlaughterGroupRepository) FindId(id string) (*entity.SlaughterGroup, error) {
    return r.Impl.FindById(id)
}

func (r *SlaughterGroupRepository) Add(newModel *entity.SlaughterGroup) (*entity.SlaughterGroup, error) {
    return r.Impl.Add(newModel)
}

func (r *SlaughterGroupRepository) Save(model *entity.SlaughterGroup) error {
    return r.Impl.Save(model)
}

func (r *SlaughterGroupRepository) Delete(id string) error {
    return r.Impl.Delete(id)
}
