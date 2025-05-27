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
    selectQuery := util.NewSelectQuery(util.SELECT,
		*util.NewNamedGroup("group", "id", "insemination_date"),
		*util.NewNamedGroup("bull", "id", "name"),
		*util.NewNamedGroup("bull_mother", "id", "name"),
		*util.NewNamedGroup("bull_father", "id", "name")).
		From("insemination_groups as group").
		Joins(
            "left join animals as bull on bull.id = group.bull_id",
			"left join animals as bull_mother on bull_mother.id = bull.mother_id",
			"left join animals as bull_father on bull_father.id = bull.father_id")

	insertQuery := util.NewInsertQuery("insemination_groups", "id", "bull_id",
		"insemination_date", "created_at", "user_id")
	updateQuery := util.NewUpdateQuery("insemination_groups", "bull_id",
		"insemination_date", "created_at", "user_id")

	r.Impl = RepositoryImpl[entity.InseminationGroup]{
		TableName:       "insemination_groups",
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		Repository:      r,
	}
}

func (r *InseminationGroupRepository) setNewEntity(model *entity.InseminationGroup, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
}

func (r *InseminationGroupRepository) buildEntity(row *sql.Row) (model *entity.InseminationGroup, err error) {
	var group entity.InseminationGroup
	err = row.Scan(&group.Id, &group.InseminationDate, &group.Bull.Id, &group.Bull.Name,
		&group.Bull.MotherId, &group.Bull.MotherName,
		&group.Bull.FatherId, &group.Bull.FatherName)
	return &group, err
}

func (r *InseminationGroupRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.InseminationGroup, err error) {
	var groups []entity.InseminationGroup
	for rows.Next() {
		var group entity.InseminationGroup
		err = rows.Scan(&group.Id, &group.InseminationDate, &group.Bull.Id, &group.Bull.Name,
			&group.Bull.MotherId, &group.Bull.MotherName,
			&group.Bull.FatherId, &group.Bull.FatherName)
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

func (r *InseminationGroupRepository) Add(newModel *entity.InseminationGroup) (*entity.InseminationGroup, error) {
	return r.Impl.Add(newModel)
}

func (r *InseminationGroupRepository) Save(model *entity.InseminationGroup) error {
	return r.Impl.Save(model)
}

func (r *InseminationGroupRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
