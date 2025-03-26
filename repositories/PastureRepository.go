package repositories

import (
	"database/sql"
	"time"

	"github.com/felipeErnica/rebanho-backend/entity"
	"github.com/felipeErnica/rebanho-backend/util"
)

type PastureRepository struct{
    Impl    RepositoryImpl[entity.Pasture]
}

func (r *PastureRepository) Init() {
    selectQuery:= new(util.QueryConstructor).Select("pasture", "id", "name")
        selectQuery.AndSelect("bull", "id", "name")
        selectQuery.From("pastures_active", "pasture")
        selectQuery.LeftJoin("animals", "bull").On("bull.id", "pasture.bull_id")
    updateQuery:= new(util.QueryConstructor).Update("pastures", "name", "bull_id", "created_at")
    insertQuery:= new(util.QueryConstructor).Insert("pastures", "id", "name", "bull_id", "created_at")

    r.Impl = RepositoryImpl[entity.Pasture]{
        TableName: "pastures",
        SelectQueryBody: selectQuery.Build(),
        InsertQuery: insertQuery.Build(),
        UpdateQuery: updateQuery.Build(),
        Repository: r,
    }
}

func (r *PastureRepository) setNewEntity(model *entity.Pasture, id string, createdAt time.Time) {
    model.Id = id
    model.CreatedAt = createdAt
}

func (r *PastureRepository) buildEntity(row *sql.Row) (model *entity.Pasture, err error) {
    var pasture entity.Pasture
    err = row.Scan(&pasture.Id, &pasture.Name, &pasture.Bull.Id, &pasture.Bull.Name)
    return &pasture, err
}

func (r *PastureRepository) buildListEntity(rows *sql.Rows) (arr *[]entity.Pasture, err error) {
    var pastures []entity.Pasture
    for rows.Next() {
        var pasture entity.Pasture
        err = rows.Scan(&pasture.Id, &pasture.Name, &pasture.Bull.Id, &pasture.Bull.Name)
        if err != nil {
            return
        }
        pastures = append(pastures, pasture)
    }
    return &pastures, err
}

func (r *PastureRepository) saveOrUpdateScan(query string, model *entity.Pasture) error {
    err:= execQuery(query, model.Id, model.Name, model.Bull.Id, model.CreatedAt)
    return err
}

func (r *PastureRepository) FindAll() (*[]entity.Pasture, error) {
    return r.Impl.FindAll()
}

func (r *PastureRepository) FindById(id string) (*entity.Pasture, error) {
    return r.Impl.FindById(id)
}

func (r *PastureRepository) Add(newPasture entity.Pasture) (*entity.Pasture, error){
    return r.Impl.Add(newPasture)
}

func (r *PastureRepository) Save(pasture *entity.Pasture) (error){
    return r.Impl.Save(pasture)
}

func (r *PastureRepository) Delete(id string) (error){
    return r.Impl.SoftDelete(id)
}
