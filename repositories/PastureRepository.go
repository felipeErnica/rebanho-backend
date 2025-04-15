package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PastureRepository struct {
	Impl RepositoryImpl[entity.Pasture]
}

func (r *PastureRepository) Init() {
	selectQuery := new(util.SelectConstructor).Select("pasture", "id", "name")
	selectQuery.AndSelect("bull", "id", "name")
	selectQuery.AndSelect("farm", "id", "name")
	selectQuery.From("pastures", "pasture")
	selectQuery.LeftJoin("animals", "bull").On("bull.id", "pasture.bull_id")
	selectQuery.LeftJoin("farms", "farm").On("farm.id", "pasture.farm_id")
	updateQuery := new(util.SelectConstructor).Update("pastures", "name", "bull_id", "farm_id", "created_at", "user_id")
	insertQuery := new(util.SelectConstructor).Insert("pastures", "id", "name", "bull_id", "farm_id", "created_at", "user_id")

	r.Impl = RepositoryImpl[entity.Pasture]{
		TableName:       "pastures",
		SelectQueryBody: *selectQuery,
		InsertQuery:     *insertQuery,
		UpdateQuery:     *updateQuery,
		Repository:      r,
	}
}

func (r *PastureRepository) setNewEntity(model *entity.Pasture, id string, createdAt time.Time) {
	model.Id = id
	model.CreatedAt = createdAt
	model.UserId = GetUserId()
}

func (r *PastureRepository) buildEntity(row *sql.Row) (model *entity.Pasture, err error) {
	var pasture entity.Pasture
	err = row.Scan(&pasture.Id, &pasture.Name,
		&pasture.Bull.Id, &pasture.Bull.Name,
		&pasture.Farm.Id, &pasture.Farm.Name)
	return &pasture, err
}

func (r *PastureRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.Pasture, err error) {
	var pastures []entity.Pasture
	for rows.Next() {
        var pasture entity.Pasture
        err = rows.Scan(&pasture.Id, &pasture.Name,
            &pasture.Bull.Id, &pasture.Bull.Name,
            &pasture.Farm.Id, &pasture.Farm.Name)
        if err != nil {
            return
        }
        pastures = append(pastures, pasture)
	}
	return &pastures, err
}

func (r *PastureRepository) saveOrUpdateScan(query string, model *entity.Pasture) error {
	err := execQuery(query, model.Id, model.Name, model.Bull.Id, model.Farm.Id, model.CreatedAt, model.UserId)
	return err
}

func (r *PastureRepository) FindAll() (*[]entity.Pasture, error) {
	return r.Impl.FindAll()
}

func (r *PastureRepository) FindById(id string) (*entity.Pasture, error) {
	return r.Impl.FindById(id)
}

func (r *PastureRepository) Add(newPasture *entity.Pasture) (*entity.Pasture, error) {
	return r.Impl.Add(newPasture)
}

func (r *PastureRepository) Save(pasture *entity.Pasture) error {
	return r.Impl.Save(pasture)
}

func (r *PastureRepository) Delete(id string) error {
	return r.Impl.Delete(id)
}
