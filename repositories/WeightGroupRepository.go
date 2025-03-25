package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type WeightGroupRepository struct {
    Impl RepositoryImpl[entity.WeightGroup]
}

func (r *WeightGroupRepository) Init() {
    selectQuery:=new(util.QueryConstructor).Select("", "id", "weighted_date").From("weight_groups", "")
    updateQuery:=new(util.QueryConstructor).Update("weight_groups", "weighted_date", "created_at", "user_id")
    insertQuery:=new(util.QueryConstructor).Insert("weight_groups", "id", "weighted_date", "created_at", "user_id")

    r.Impl = RepositoryImpl[entity.WeightGroup]{
        Repository: r,
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
    }
}

func (r *WeightGroupRepository) setNewEntity(model *entity.WeightGroup, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *WeightGroupRepository) buildEntity(row *sql.Row) (model *entity.WeightGroup, err error) {
    var group entity.WeightGroup
    err = row.Scan(group.Id, group.WeightDate)
    return &group, err
}

func (r *WeightGroupRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.WeightGroup, err error) {
    var groups []entity.WeightGroup
    for rows.Next() {
        var group entity.WeightGroup
        err = rows.Scan(group.Id, group.WeightDate)
        if err != nil {
            return
        }
        groups = append(groups, group)
    }
    return &groups, err
}

func (r *WeightGroupRepository) saveOrUpdateScan(query string, model *entity.WeightGroup) error {
    return ExecQuery(query, model.Id, model.WeightDate, model.CreatedAt, model.UserId)
}

func (r *WeightGroupRepository) FindAll() (*[]entity.WeightGroup, error) {
    return r.Impl.FindAll()
}

func (r *WeightGroupRepository) FindById(id string) (*entity.WeightGroup, error) {
    return r.Impl.FindById(id)
}

func (r *WeightGroupRepository) Add(newModel entity.WeightGroup) (*entity.WeightGroup, error) {
    return r.Impl.Add(newModel)
}

func (r *WeightGroupRepository) Save(model *entity.WeightGroup) error {
    return r.Impl.Save(model)
}

func (r *WeightGroupRepository) Delete(id string) error {
    return r.Impl.HardDelete(id)
}
